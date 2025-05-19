package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/internal/executor/service"
	"github.com/mavrk-mose/pay/pkg/utils"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	"go.uber.org/zap"
)

type WebhookHandler struct {
	Logger     utils.Logger
	Provider   *service.PayPalProvider
	WebhookID  string
	Secret     string
}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{
		Secret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
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
	event, err := webhook.ConstructEvent(payload, sigHeader, h.Secret)
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

		//TODO: update wallet & ledger

		h.Logger.Info("Published checkout.session.completed", zap.String("session_id", session.ID))

	case "payment_intent.succeeded":
		var intent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
			h.Logger.Error("Failed to parse payment intent", zap.Error(err))
			c.String(http.StatusBadRequest, "Bad intent payload")
			return
		}

		//TODO: update wallet & ledger

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

	//TODO: process in wallet & ledger
	h.Logger.Info("Publishing event", zap.String("subject", subject), zap.ByteString("payload", payload))

	c.Status(http.StatusOK)
}