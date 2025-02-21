package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your_project/service"
)

type ReferralHandler struct {
	ReferralSvc service.ReferralService
}

func NewReferralHandler(refSvc service.ReferralService) *ReferralHandler {
	return &ReferralHandler{ReferralSvc: refSvc}
}

// GenerateReferralCode handles a GET request to generate a referral code for a user.
func (h *ReferralHandler) GenerateReferralCode(c *gin.Context) {
	userID := c.Param("userID")
	code, err := h.ReferralSvc.GenerateReferralCode(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"referral_code": code})
}

// ApplyReferralCode handles a POST request to apply a referral code.
func (h *ReferralHandler) ApplyReferralCode(c *gin.Context) {
	userID := c.Param("userID")
	var req struct {
		ReferralCode string `json:"referral_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.ReferralSvc.ApplyReferralCode(userID, req.ReferralCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Referral applied successfully"})
}

// GetReferralBonus handles a GET request to retrieve the referral bonus for a user.
func (h *ReferralHandler) GetReferralBonus(c *gin.Context) {
	userID := c.Param("userID")
	bonus, err := h.ReferralSvc.GetReferralBonus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bonus": bonus})
}
