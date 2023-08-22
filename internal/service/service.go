package service

import "github.com/ankit/project/credit-card-offer-limit/internal/db"

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