package v1

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Resposne struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, Resposne{Message: message})
}
