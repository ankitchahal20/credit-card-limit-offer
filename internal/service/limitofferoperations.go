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

		if limitOffer.LimitType == nil || limitOffer.AccountID == nil || limitOffer.NewLimit == nil ||
			limitOffer.OfferActivationTime == nil || limitOffer.OfferExpiryTime == nil {
			utils.Logger.Error(fmt.Sprintf("one of the following account LimitType, AccountID, NewLimit or OfferActivationTime or OfferExpiryTime field is missing while creating an account, txid : %v", txid))
			
			errMessage := "one of the following account limit, last account limit, per transaction limit or last per transaction limit field is missing while creating an account"
			utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
			return	
		}

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

	// create the offer id
	limitOffer.ID = uuid.New().String()
	limitOffer.Status = models.Pending

	// if current limit is greater than existing limit for an account, then update the existing time with status as PENDING
	utils.Logger.Info(fmt.Sprintf("calling db layer for creating limit offer for %v account, txid : %v", limitOffer.AccountID, txid))
	err = service.repo.CreateLimitOffer(ctx, limitOffer)
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
		
		if activeLimitOffer.AccountID == "" {
			utils.Logger.Error(fmt.Sprintf("account is missing, txid : %v", txid))
			
			errMessage := "account is missing to list all active limit offers"
			utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
			return
		}

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
		
		if updateLimitOfferStatus.LimitOfferID == "" {
			utils.Logger.Error(fmt.Sprintf("limit offer id is missing, txid : %v", txid))
			
			errMessage := "limit offer id is missing"
			utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
			return
		}

		_, errlimitOfferUUID := uuid.Parse(updateLimitOfferStatus.LimitOfferID)
		if errlimitOfferUUID != nil {
			utils.Logger.Error(fmt.Sprintf("Error parsing the %v limitOfferID, txid : %v", updateLimitOfferStatus.LimitOfferID, txid))
			utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidOfferLimitID)
			return
		}

		switch updateLimitOfferStatus.Status{
		case string(models.Accepted), string(models.Rejected):
		default:
			utils.Logger.Error(fmt.Sprintf("invalid status is provided, txid : %v", txid))
			
			errMessage := "received status is not supported"
			utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
			return
		}

		activeLimitOffers, err := creditCardLimitOfferClient.updateLimitOfferStatus(ctx, updateLimitOfferStatus)
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

func (service *CreditCardLimitOfferService) updateLimitOfferStatus(ctx *gin.Context, updateLimitOfferStatus models.UpdateLimitOfferStatus) ([]models.LimitOffer, *limitoffererror.CreditCardError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	utils.Logger.Info(fmt.Sprintf("calling db layer to update limit offer status account, txid : %v", txid))
	fetchedAccount, err := service.repo.UpdateLimitOfferStatus(ctx, updateLimitOfferStatus)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer while updating limit offer status, txid : %v", txid))
		return []models.LimitOffer{}, err
	}

	return fetchedAccount, nil
}