package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/executor/pkg"
)

func main() {
	r := gin.Default()

	PORT := pkg.GetEnv("PORT")
	err := r.Run(":" + PORT)
	if err != nil {
		return
	}
}
