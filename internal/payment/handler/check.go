package handler

import "github.com/gin-gonic/gin"

func (t *PaymentHandler) Check(c *gin.Context) {
	c.Status(200)
}
