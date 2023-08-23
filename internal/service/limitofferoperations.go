package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ankit/project/credit-card-offer-limit/internal/constants"
	"github.com/ankit/project/credit-card-offer-limit/internal/limitoffererror"
	"github.com/ankit/project/credit-card-offer-limit/internal/models"
	"github.com/ankit/project/credit-card-offer-limit/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// This function is responsible for account creation
func CreateLimitOffer() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for account creation, txid : %v", txid))
		var limitOffer models.LimitOffer
		if err := ctx.ShouldBindBodyWith(&limitOffer, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request for account creation is unmarshalled successfully, txid : %v", txid))

			offerLimitID, err := creditCardLimitOfferClient.createLimitOffer(ctx, limitOffer)
			if err != nil {
				utils.RespondWithError(ctx, err.Code, err.Message)
				return
			}

			ctx.JSON(http.StatusOK, map[string]string{
				"offer_limit_id": offerLimitID,
			})

			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *CreditCardLimitOfferService) createLimitOffer(ctx *gin.Context, limitOffer models.LimitOffer) (string, *limitoffererror.CreditCardError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	utils.Logger.Info(fmt.Sprintf("calling db layer for fetching account %v info to get the existing limit, txid : %v", limitOffer.AccountID, txid))
	
	fetchedAccount, err := service.repo.GetAccount(ctx, *limitOffer.AccountID)
	if err != nil {
		return constants.EmptyString, err
	}
	isLimitOfferExsits, offerLimitID, err := service.repo.IsLimitOfferExists(ctx, limitOffer)
	if err != nil {
		return constants.EmptyString, err
	}
	fmt.Println("isLimitOfferExsits : ",isLimitOfferExsits)
	var currentLimit int
	if *limitOffer.LimitType == models.AccountLimit {
		currentLimit = *fetchedAccount.AccountLimit
	} else if *limitOffer.LimitType == models.PerTransactionLimit {
		currentLimit = *fetchedAccount.PerTransactionLimit
	}
	fmt.Println("*limitOffer.NewLimit <= currentLimit ", *limitOffer.NewLimit , ": ", currentLimit)
	if *limitOffer.NewLimit <= currentLimit {
		return constants.EmptyString, &limitoffererror.CreditCardError{
			Code:    http.StatusBadRequest,
			Message: "offer limit is less than or equal to existing limit",
			Trace:   txid,
		}
	}

	if !isLimitOfferExsits {
		// create the offer id
		limitOffer.ID = uuid.New().String()
	} else {
		limitOffer.ID = offerLimitID
	}

	limitOffer.Status = models.Pending

	// if current limit is greater than existing limit for an account, then update the existing time with status as PENDING
	utils.Logger.Info(fmt.Sprintf("calling db layer for creating limit offer for %v account, txid : %v", limitOffer.AccountID, txid))
	err = service.repo.CreateLimitOffer(ctx, limitOffer, isLimitOfferExsits)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer while creating creating limit offer for %v account, txid : %v", limitOffer.AccountID, txid))
		return constants.EmptyString, err
	}

	return limitOffer.ID, nil
}

// This function is responsible to list all active limit offers
func ListActiveLimitOffers() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for to list all active limit offer, txid : %v", txid))
		var activeLimitOffer models.ActiveLimitOffer
		if err := ctx.ShouldBindBodyWith(&activeLimitOffer, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("received request for account creation is unmarshalled successfully, txid : %v", txid))

			activeLimitOffers, err := creditCardLimitOfferClient.listActiveLimitOffers(ctx, activeLimitOffer)
			if err != nil {
				utils.RespondWithError(ctx, err.Code, err.Message)
				return
			}
		
			ctx.JSON(http.StatusOK, activeLimitOffers)

			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *CreditCardLimitOfferService) listActiveLimitOffers(ctx *gin.Context, activeLimitOffer models.ActiveLimitOffer) ([]models.LimitOffer, *limitoffererror.CreditCardError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	
	if activeLimitOffer.ActiveDate == nil {
		time := time.Now().UTC()
		activeLimitOffer.ActiveDate = &time
	}

	utils.Logger.Info(fmt.Sprintf("calling db layer to list all active limit offer for %v account, txid : %v", activeLimitOffer.AccountID, txid))
	fetchedAccount, err := service.repo.ListActiveLimitOffers(ctx, activeLimitOffer)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer while fetching all active limit offer for %v account, txid : %v", activeLimitOffer.AccountID, txid))
		return []models.LimitOffer{}, err
	}

	return fetchedAccount, nil
}

// This function is responsible to update limit offer status
func UpdateLimitOfferStatus() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for to list all active limit offer, txid : %v", txid))
		var updateLimitOfferStatus models.UpdateLimitOfferStatus
		if err := ctx.ShouldBindBodyWith(&updateLimitOfferStatus, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("received request for account creation is unmarshalled successfully, txid : %v", txid))

			err := creditCardLimitOfferClient.updateLimitOfferStatus(ctx, updateLimitOfferStatus)
			if err != nil {
				utils.RespondWithError(ctx, err.Code, err.Message)
				return
			}
		
			ctx.JSON(http.StatusOK, nil)
			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *CreditCardLimitOfferService) updateLimitOfferStatus(ctx *gin.Context, updateLimitOfferStatus models.UpdateLimitOfferStatus) *limitoffererror.CreditCardError {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	limitOfferInfo, err := service.repo.GetLimitOffer(ctx, updateLimitOfferStatus.LimitOfferID)
	if err != nil {
		return err
	}
	fmt.Println("limitOfferInfo.Status : ", limitOfferInfo.Status, limitOfferInfo)
	if limitOfferInfo.Status == models.Accepted {
		return &limitoffererror.CreditCardError{
			Code:    http.StatusUnprocessableEntity,
			Message: "limit offer is already in accepted state",
			Trace:   txid,
		}
	} else if limitOfferInfo.Status == models.Rejected {
		return &limitoffererror.CreditCardError{
			Code:    http.StatusUnprocessableEntity,
			Message: "limit offer is already in rejected state",
			Trace:   txid,
		}
	}

	utils.Logger.Info(fmt.Sprintf("calling db layer to update limit offer status account, txid : %v", txid))
	err = service.repo.UpdateLimitOfferStatus(ctx, updateLimitOfferStatus)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer while updating limit offer status, txid : %v", txid))
		return err
	}

	return nil
}
