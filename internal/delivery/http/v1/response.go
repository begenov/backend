package v1

import (
	"log"

	"github.com/begenov/backend/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Resposne struct {
	Message string `json:"message"`
}

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

func newResponse(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, Resposne{Message: message})
}
