package service

import (
	"fmt"

	"github.com/ankit/project/credit-card-offer-limit/internal/db"
	"github.com/gin-gonic/gin"
)

var (
	creditCardLimitOfferClient *CreditCardLimitOfferService
)

type CreditCardLimitOfferService struct {
	repo db.CreditCardLimitOfferService
}

func NewCreditCardLimitOfferService(conn db.CreditCardLimitOfferService) *CreditCardLimitOfferService {
	creditCardLimitOfferClient = &CreditCardLimitOfferService{
		repo: conn,
	}
	return creditCardLimitOfferClient
}

// This is a function to create account
func CreateAccount() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Retrun from create Account")
	}
}

// This is a function to create account
func GetAccount() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Return from get Account")
	}
}