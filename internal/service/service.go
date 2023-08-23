package service

import (
	"sync"

	"github.com/ankit/project/credit-card-offer-limit/internal/db"
	"github.com/ankit/project/credit-card-offer-limit/internal/utils"
)

var (
	creditCardLimitOfferClient *CreditCardLimitOfferService
	once                       sync.Once
)

type CreditCardLimitOfferService struct {
	repo db.CreditCardLimitOfferService
}

// creditCardLimitOfferClient should only be created once throughtout the application lifetime
func NewCreditCardLimitOfferService(conn db.CreditCardLimitOfferService) *CreditCardLimitOfferService {
	if creditCardLimitOfferClient == nil {
		once.Do(
			func() {
				creditCardLimitOfferClient = &CreditCardLimitOfferService{
					repo: conn,
				}
			})
	} else {
		utils.Logger.Info("creditCardLimitOfferClient is alredy created")
	}
	return creditCardLimitOfferClient
}
