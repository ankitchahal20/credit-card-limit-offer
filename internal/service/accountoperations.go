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
func CreateAccount() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for account creation, txid : %v", txid))
		var accountInfo models.Account
		if err := ctx.ShouldBindBodyWith(&accountInfo, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request for account creation is unmarshalled successfully, txid : %v", txid))

		if accountInfo.AccountLimit == nil || accountInfo.LastAccountLimit == nil || accountInfo.PerTransactionLimit == nil ||
			accountInfo.LastPerTransactionLimit == nil {
			utils.Logger.Error(fmt.Sprintf("one of the following account limit, last account limit, per transaction limit or last per transaction limit field is missing while creating an account, txid : %v", txid))
			
			errMessage := "one of the following account limit, last account limit, per transaction limit or last per transaction limit field is missing while creating an account"
			utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
			return
		}

		createdAccount, err := creditCardLimitOfferClient.createAccount(ctx, accountInfo)
		if err != nil {
			utils.RespondWithError(ctx, err.Code, err.Message)
			return
		}

		ctx.JSON(http.StatusOK, map[string]string{
			"account_id": createdAccount.AccountID,
			"customer_id": createdAccount.CustomerID,
		})

		ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *CreditCardLimitOfferService) createAccount(ctx *gin.Context, accountInfo models.Account) (models.Account, *limitoffererror.CreditCardError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	
	// check if per transaction limit is greater than account limit
	if *accountInfo.PerTransactionLimit > *accountInfo.AccountLimit {
		utils.Logger.Info(fmt.Sprintf("per transaction limit can not be greater than account limit, txid : %v", txid))
		return models.Account{},  &limitoffererror.CreditCardError{
			Code:    http.StatusBadRequest,
			Message: "per transaction limit can not be greater than account limit",
			Trace:   txid,
		}
	}

	// generate the accountID and customerID from uuid package and set in the the accountInfo
	accountID := uuid.New().String()
	customerID := uuid.New().String()
	accountInfo.AccountID = accountID
	accountInfo.CustomerID = customerID
	
	// set the current time as accountCreationTime and AccountLimitUpdateTime in the accountInfo
	accountCreationTime := time.Now().UTC()
	accountInfo.AccountLimitUpdateTime = accountCreationTime
	accountInfo.PerTransactionLimitUpdateTime = accountCreationTime

	utils.Logger.Info(fmt.Sprintf("calling db layer for account creation, txid : %v", txid))
	err := service.repo.CreateAccount(ctx, accountInfo)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer during account creation, txid : %v", txid))
		return models.Account{}, err
	}

	return accountInfo, nil
}


// This function is responsible to get a specific account based on accountid
func GetAccount() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		accountID := ctx.Param(constants.AccountID)
		utils.Logger.Info(fmt.Sprintf("request received for get %v account, txid : %v", accountID, txid))
		_, erraccountUUID := uuid.Parse(accountID)
		if erraccountUUID != nil {
			utils.Logger.Error(fmt.Sprintf("Error parsing the %v accountID, txid : %v", accountID, txid))
			utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidAccountID)
			return
		}
		utils.Logger.Info(fmt.Sprintf("calling service layer for getting %v accountID, txid : %v", accountID, txid))
		fetchedAccount, err := creditCardLimitOfferClient.getAccount(ctx, accountID)
		if err != nil {
			utils.RespondWithError(ctx, err.Code, err.Message)
			return
		}

		ctx.JSON(http.StatusOK, fetchedAccount)
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

func (service *CreditCardLimitOfferService) getAccount(ctx *gin.Context, accountID string) (models.Account, *limitoffererror.CreditCardError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	
	utils.Logger.Info(fmt.Sprintf("calling db layer for getting %v account, txid : %v", accountID, txid))
	fetchedAccount, err := service.repo.GetAccount(ctx, accountID)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer during getting %v account, txid : %v", accountID, txid))
		return models.Account{}, err
	}

	return fetchedAccount, nil
}