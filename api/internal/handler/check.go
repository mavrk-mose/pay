package handler

import "github.com/gin-gonic/gin"

func (t *ApiHandler) Check(c *gin.Context) {
	c.Status(200)
}
