package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/api/pkg"
)

func main() {
	r := gin.Default()

	PORT := pkg.GetEnv("PORT")
	err := r.Run(":" + PORT)
	if err != nil {
		return
	}
}
