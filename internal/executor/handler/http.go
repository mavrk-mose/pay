package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/executor/service"
	"github.com/mavrk-mose/pay/pkg/nats"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type WebhookHandler struct {
	Logger    *zap.Logger
	Nats      *nats.Client
	Provider  *service.PayPalProvider
	WebhookID string
	Secret    string
}

func NewWebhookHandler(cfg *config.Config, natsClient *nats.Client) *WebhookHandler {
	return &WebhookHandler{
		Logger: zap.Must(zap.NewProduction()),
		Nats:   natsClient,
		Secret: cfg.Stripe.Secret,
	}
}

func (h *WebhookHandler) StripeWebhookHandler(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.Logger.Error("Failed to read request body", zap.Error(err))
		c.String(http.StatusRequestEntityTooLarge, "Request too large")
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEventWithOptions(payload, sigHeader, h.Secret, webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})
	if err != nil {
		h.Logger.Error("Invalid Stripe signature", zap.Error(err))
		c.String(http.StatusBadRequest, "Invalid signature")
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			h.Logger.Error("Failed to parse session", zap.Error(err))
			c.String(http.StatusBadRequest, "Bad session payload")
			return
		}

		err := h.Nats.Publish("payments.success", map[string]any{
			"event":          "checkout.session.completed",
			"session_id":     session.ID,
			"customer_id":    session.Customer.ID,
			"amount_total":   session.AmountTotal,
			"currency":       session.Currency,
			"payment_intent": session.PaymentIntent.ID,
		})
		if err != nil {
			h.Logger.Error("Failed to publish session", zap.Error(err))
			c.String(http.StatusInternalServerError, "Failed to publish event")
			return
		}

		h.Logger.Info("Published checkout.session.completed", zap.String("session_id", session.ID))

	case "payment_intent.succeeded":
		var intent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
			h.Logger.Error("Failed to parse payment intent", zap.Error(err))
			c.String(http.StatusBadRequest, "Bad intent payload")
			return
		}

		err := h.Nats.Publish("payments.success", map[string]any{
			"event":       "payment_intent.succeeded",
			"intent_id":   intent.ID,
			"amount":      intent.Amount,
			"currency":    intent.Currency,
			"customer_id": intent.Customer.ID,
			"description": intent.Description,
		})
		if err != nil {
			h.Logger.Error("Failed to publish payment intent", zap.Error(err))
			c.String(http.StatusInternalServerError, "Failed to publish event")
			return
		}

		h.Logger.Info("Published payment_intent.succeeded", zap.String("intent_id", intent.ID))

	default:
		h.Logger.Warn("Unhandled Stripe event", zap.String("type", string(event.Type)))
	}

	c.String(http.StatusOK, "Received")
}

func (h *WebhookHandler) PaypalWebhookHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.Logger.Error("Failed to read request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	headers := map[string]string{
		"PAYPAL-AUTH-ALGO":         c.GetHeader("Paypal-Auth-Algo"),
		"PAYPAL-CERT-URL":          c.GetHeader("Paypal-Cert-Url"),
		"PAYPAL-TRANSMISSION-ID":   c.GetHeader("Paypal-Transmission-Id"),
		"PAYPAL-TRANSMISSION-SIG":  c.GetHeader("Paypal-Transmission-Sig"),
		"PAYPAL-TRANSMISSION-TIME": c.GetHeader("Paypal-Transmission-Time"),
	}

	// Step 1: Verify the webhook
	valid, err := h.Provider.VerifyWebhook(headers, body, h.WebhookID)
	if err != nil || !valid {
		h.Logger.Error("Webhook verification failed", zap.Error(err), zap.Bool("valid", valid))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	// Step 2: Parse event payload
	var event map[string]any
	if err := json.Unmarshal(body, &event); err != nil {
		h.Logger.Error("Failed to parse webhook event", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	eventType, _ := event["event_type"].(string)
	resource := event["resource"]

	h.Logger.Info("Received PayPal webhook",
		zap.String("event_type", eventType),
		zap.Any("resource", resource),
	)

	subject := "payment.paypal." + eventType
	payload, _ := json.Marshal(event)

	if err := h.Nats.Publish(subject, payload); err != nil {
		h.Logger.Error("Failed to publish event to NATS", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish event"})
		return
	}

	c.Status(http.StatusOK)
}
