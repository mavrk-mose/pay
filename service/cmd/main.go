package main

import (
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()

  PORT := pkg.GetEnv("PORT")
  r.Run(":" + PORT)
}