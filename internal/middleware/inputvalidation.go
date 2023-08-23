package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ankit/project/credit-card-offer-limit/internal/constants"
	"github.com/ankit/project/credit-card-offer-limit/internal/models"
	"github.com/ankit/project/credit-card-offer-limit/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// This function gets the unique transactionID
func getTransactionID(c *gin.Context) string {
	transactionID := c.GetHeader(constants.TransactionID)
	_, err := uuid.Parse(transactionID)
	if err != nil {
		transactionID = uuid.New().String()
		c.Set(constants.TransactionID, transactionID)
	}
	return transactionID
}

func ValidateInputRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get the transactionID from headers if not present create a new.
		transactionID := getTransactionID(ctx)
		fmt.Printf("TimeStamp : %v", time.Now().UTC)
		path := ctx.Request.URL.String()
		switch {
		case strings.Contains(path, constants.CreateAccount):
			validateCreateAccountInput(ctx, transactionID)
		case strings.Contains(path, constants.GetAccount):
			validateGetAccountInput(ctx, transactionID)
		case strings.Contains(path, constants.CreateLimitOffer):
			validateCreateLimitOfferInput(ctx, transactionID)
		case strings.Contains(path, constants.ListActiveLimitOffers):
			validateListActiveLimitOffersInput(ctx, transactionID)
		case strings.Contains(path, constants.UpdateLimitOfferStatus):
			validateUpdateLimitOfferStatusInput(ctx, transactionID)
		}
		fmt.Println("txid : ", transactionID)

		ctx.Next()
	}
}

func validateCreateAccountInput(ctx *gin.Context, txid string) {
	var accountInfo models.Account
	err := ctx.ShouldBindBodyWith(&accountInfo, binding.JSON)
	if err != nil {
		utils.Logger.Error("error while unmarshaling the request field for create account data validation")
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidBodyCreateAccount)
		return
	}

	if accountInfo.AccountLimit == nil {
		utils.Logger.Error(fmt.Sprintf("account_limit field is missing while creating an account, txid : %v", txid))
		errMessage := "account_limit field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	if accountInfo.LastAccountLimit == nil {
		utils.Logger.Error(fmt.Sprintf("last_account_limit field is missing while creating an account, txid : %v", txid))
		errMessage := "last_account_limit field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	if accountInfo.PerTransactionLimit == nil {
		utils.Logger.Error(fmt.Sprintf("per_transaction_limit field is missing while creating an account, txid : %v", txid))
		errMessage := "per_transaction_limit field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	if accountInfo.LastPerTransactionLimit == nil {
		utils.Logger.Error(fmt.Sprintf("last_per_transaction_limit field is missing while creating an account, txid : %v", txid))
		errMessage := "last_per_transaction_limit field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}

	if *accountInfo.AccountLimit < *accountInfo.LastAccountLimit {
		utils.Logger.Error(fmt.Sprintf("amount_limit field should be greater than or equal to last_amount_limit while creating an account, txid : %v", txid))
		errMessage := "amount_limit is less than last_amount_limit"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}

	if *accountInfo.PerTransactionLimit < *accountInfo.LastPerTransactionLimit {
		utils.Logger.Error(fmt.Sprintf("per_transaction_limit field should be greater than or equal to last_per_transaction_limit while creating an account, txid : %v", txid))
		errMessage := "per_transaction_limit is less than last_per_transaction_limit"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
}

func validateGetAccountInput(ctx *gin.Context, txid string) {
	accountID := ctx.Param(constants.AccountID)
	fmt.Println("accountID : ", accountID)
	utils.Logger.Info(fmt.Sprintf("request received for get %v account, txid : %v", accountID, txid))
	_, erraccountUUID := uuid.Parse(accountID)
	if erraccountUUID != nil {
		utils.Logger.Error(fmt.Sprintf("Error parsing the %v accountID, txid : %v", accountID, txid))
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidAccountID)
		return
	}
}

func validateCreateLimitOfferInput(ctx *gin.Context, txid string) {
	var limitOffer models.LimitOffer
	err := ctx.ShouldBindBodyWith(&limitOffer, binding.JSON)
	if err != nil {
		utils.Logger.Error("error while unmarshaling the request field for create limit offer data validation")
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidBodyCreateLimitOffer)
		return
	}

	if limitOffer.LimitType == nil {
		utils.Logger.Error(fmt.Sprintf("limit_type field is missing while creating limit offer, txid : %v", txid))
		errMessage := "limit_type field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	if limitOffer.AccountID == nil {
		utils.Logger.Error(fmt.Sprintf("account_id field is missing while creating limit offer, txid : %v", txid))
		errMessage := "account_id field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	if limitOffer.NewLimit == nil {
		utils.Logger.Error(fmt.Sprintf("new_limit field is missing while creating limit offer, txid : %v", txid))
		errMessage := "new_limit field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	if limitOffer.OfferActivationTime == nil {
		utils.Logger.Error(fmt.Sprintf("offer_activation_time field is missing while creating limit offer, txid : %v", txid))
		errMessage := "offer_activation_time field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	if limitOffer.OfferExpiryTime == nil {
		utils.Logger.Error(fmt.Sprintf("offer_expiry_time field is missing while creating limit offer, txid : %v", txid))
		errMessage := "offer_expiry_time field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	// offer_expiry_time field should be greater than offer_activation_time
	isExpireTimeBeforeActivationTime := limitOffer.OfferExpiryTime.Before(*limitOffer.OfferActivationTime)
	if isExpireTimeBeforeActivationTime {
		utils.Logger.Error(fmt.Sprintf("offer_expiry_time field should be greater than offer_activation_time, txid : %v", txid))
		errMessage := "offer_expiry_time field should be greater than offer_activation_time"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
}

func validateListActiveLimitOffersInput(ctx *gin.Context, txid string) {
	var activeLimitOffer models.ActiveLimitOffer
	err := ctx.ShouldBindBodyWith(&activeLimitOffer, binding.JSON)
	if err != nil {
		utils.Logger.Error("error while unmarshaling the request field to list active limit offer data validation")
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidBody)
		return
	}

	if activeLimitOffer.AccountID == "" {
		utils.Logger.Error(fmt.Sprintf("account_id field is missing to list active limit offers, txid : %v", txid))
		errMessage := "account_id field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}

	_, erraccountUUID := uuid.Parse(activeLimitOffer.AccountID)
	if erraccountUUID != nil {
		utils.Logger.Error(fmt.Sprintf("Error parsing the %v accountID, txid : %v", activeLimitOffer.AccountID, txid))
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidAccountID)
		return
	}
}

func validateUpdateLimitOfferStatusInput(ctx *gin.Context, txid string) {
	var updateLimitOfferStatus models.UpdateLimitOfferStatus
	err := ctx.ShouldBindBodyWith(&updateLimitOfferStatus, binding.JSON)
	if err != nil {
		utils.Logger.Error("error while unmarshaling the request field to update limit offer status data validation")
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidBodyUpdateLimitOfferStatus)
		return
	}
	utils.Logger.Info(fmt.Sprintf("received request for account creation is unmarshalled successfully, txid : %v", txid))

	if updateLimitOfferStatus.LimitOfferID == "" {
		utils.Logger.Error(fmt.Sprintf("limit offer id is missing, txid : %v", txid))
		errMessage := "limit_offer_id field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}

	_, errlimitOfferUUID := uuid.Parse(updateLimitOfferStatus.LimitOfferID)
	if errlimitOfferUUID != nil {
		utils.Logger.Error(fmt.Sprintf("Error parsing the %v limitOfferID, txid : %v", updateLimitOfferStatus.LimitOfferID, txid))
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidOfferLimitID)
		return
	}

	switch updateLimitOfferStatus.Status {
	case string(models.Accepted), string(models.Rejected):
	default:
		utils.Logger.Error(fmt.Sprintf("invalid status is provided, txid : %v", txid))

		errMessage := "received status is not supported"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
}
