package utils

import (
	"github.com/ankit/project/credit-card-offer-limit/internal/constants"
	"github.com/ankit/project/credit-card-offer-limit/internal/limitoffererror"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogClient() {
	Logger, _ = zap.NewDevelopment()
}

func RespondWithError(c *gin.Context, statusCode int, message string) {

	c.AbortWithStatusJSON(statusCode, limitoffererror.CreditCardError{
		Trace:   c.Request.Header.Get(constants.TransactionID),
		Code:    statusCode,
		Message: message,
	})
}
