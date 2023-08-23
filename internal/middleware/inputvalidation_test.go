package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ankit/project/credit-card-offer-limit/internal/constants"
	"github.com/ankit/project/credit-card-offer-limit/internal/models"
	"github.com/ankit/project/credit-card-offer-limit/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateAccountRequestInput(t *testing.T) {
	// init logging client
	utils.InitLogClient()

	accountLimit := 1000
	perTransactionLimit := 10000
	lastAccountLimit := 1000
	lastPerTransactionLimit := 1000

	// Case 1 : account_limit missing
	requestFields := models.Account{
		//AccountLimit: &accountLimit,
		PerTransactionLimit:           &perTransactionLimit,
		LastAccountLimit:              &lastAccountLimit,
		LastPerTransactionLimit:       &lastPerTransactionLimit,
		AccountLimitUpdateTime:        time.Now().UTC(),
		PerTransactionLimitUpdateTime: time.Now().UTC(),
	}

	jsonValue, _ := json.Marshal(requestFields)

	w := httptest.NewRecorder()
	_, e := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/v1/create_account", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Case 2 : per_transaction_limit missing
	requestFields = models.Account{
		AccountLimit: &accountLimit,
		//PerTransactionLimit: &perTransactionLimit,
		LastAccountLimit:              &lastAccountLimit,
		LastPerTransactionLimit:       &lastPerTransactionLimit,
		AccountLimitUpdateTime:        time.Now().UTC(),
		PerTransactionLimitUpdateTime: time.Now().UTC(),
	}

	jsonValue, _ = json.Marshal(requestFields)

	w = httptest.NewRecorder()
	_, e = gin.CreateTestContext(w)
	req, _ = http.NewRequest(http.MethodPost, "/v1/create_account", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Case 3 : last_account_limit missing
	requestFields = models.Account{
		AccountLimit:        &accountLimit,
		PerTransactionLimit: &perTransactionLimit,
		//LastAccountLimit: &lastAccountLimit,
		LastPerTransactionLimit:       &lastPerTransactionLimit,
		AccountLimitUpdateTime:        time.Now().UTC(),
		PerTransactionLimitUpdateTime: time.Now().UTC(),
	}

	jsonValue, _ = json.Marshal(requestFields)

	w = httptest.NewRecorder()
	_, e = gin.CreateTestContext(w)
	req, _ = http.NewRequest(http.MethodPost, "/v1/create_account", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Case 4 : last_per_transaction_limit missing
	requestFields = models.Account{
		AccountLimit:        &accountLimit,
		PerTransactionLimit: &perTransactionLimit,
		LastAccountLimit:    &lastAccountLimit,
		//LastPerTransactionLimit: &lastPerTransactionLimit,
		AccountLimitUpdateTime:        time.Now().UTC(),
		PerTransactionLimitUpdateTime: time.Now().UTC(),
	}

	jsonValue, _ = json.Marshal(requestFields)

	w = httptest.NewRecorder()
	_, e = gin.CreateTestContext(w)
	req, _ = http.NewRequest(http.MethodPost, "/v1/create_account", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateGetAccountRequestInput(t *testing.T) {
	// init logging client
	utils.InitLogClient()

	// case 1 : sending invalid uuid as account_id

	// Create a test context with the Gin engine
	r := gin.Default()
	r.Use(ValidateInputRequest())
	r.GET("/v1/get_account/:account_id", func(c *gin.Context) {})

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/get_account/74edf2ad-7ee8-45f7-a7d2-7a3287c33", nil)
	ctx.Request.Header.Add(constants.ContentType, "application/json")
	ctx.Params = []gin.Param{
		{Key: "account_id", Value: "74edf2ad-7ee8-45f7-a7d2-7a3287c33"},
	}

	// Serve the request through the Gin engine
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// case 2 : sending empty account_id
	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/get_account/", nil)
	ctx.Request.Header.Add(constants.ContentType, "application/json")
	ctx.Params = []gin.Param{
		{Key: "account_id", Value: ""},
	}

	r.ServeHTTP(w, ctx.Request)
	// Assert the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateCreateLimitOfferRequestInput(t *testing.T) {
	// init logging client
	utils.InitLogClient()

	accountID := "f83513e1-f0cb-4a49-85e4-8e9ddb1f3417"
	accountLimit := models.AccountLimit
	newLimit := 5000
	offerActivationTime := time.Now().UTC()
	OfferExpiryTime := time.Now().Add(500).UTC()

	limtOffer := models.LimitOffer{
		ID: "74edf2ad-7ee8-45f7-a7d2-7a3287c33ffe",
		//AccountID: &accountID,
		LimitType:           &accountLimit,
		NewLimit:            &newLimit,
		OfferActivationTime: &offerActivationTime,
		OfferExpiryTime:     &OfferExpiryTime,
	}

	// case 1 : account_id missing
	jsonValue, _ := json.Marshal(limtOffer)

	w := httptest.NewRecorder()
	_, e := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/v1/create_limit_offer", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// case 2 : limit_type missing
	limtOffer = models.LimitOffer{
		ID:        "74edf2ad-7ee8-45f7-a7d2-7a3287c33ffe",
		AccountID: &accountID,
		//LimitType: &accountLimit,
		NewLimit:            &newLimit,
		OfferActivationTime: &offerActivationTime,
		OfferExpiryTime:     &OfferExpiryTime,
	}

	jsonValue, _ = json.Marshal(limtOffer)

	req, _ = http.NewRequest(http.MethodPost, "/v1/create_limit_offer", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// case 3 : new_limit missing
	limtOffer = models.LimitOffer{
		ID:        "74edf2ad-7ee8-45f7-a7d2-7a3287c33ffe",
		AccountID: &accountID,
		LimitType: &accountLimit,
		//NewLimit: &newLimit,
		OfferActivationTime: &offerActivationTime,
		OfferExpiryTime:     &OfferExpiryTime,
	}

	jsonValue, _ = json.Marshal(limtOffer)

	req, _ = http.NewRequest(http.MethodPost, "/v1/create_limit_offer", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// case 4 : offer_activation_time missing
	limtOffer = models.LimitOffer{
		ID:        "74edf2ad-7ee8-45f7-a7d2-7a3287c33ffe",
		AccountID: &accountID,
		LimitType: &accountLimit,
		NewLimit:  &newLimit,
		//OfferActivationTime: &offerActivationTime,
		OfferExpiryTime: &OfferExpiryTime,
	}

	jsonValue, _ = json.Marshal(limtOffer)

	req, _ = http.NewRequest(http.MethodPost, "/v1/create_limit_offer", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// case 4 : offer_expiry_time missing
	limtOffer = models.LimitOffer{
		ID:                  "74edf2ad-7ee8-45f7-a7d2-7a3287c33ffe",
		AccountID:           &accountID,
		LimitType:           &accountLimit,
		NewLimit:            &newLimit,
		OfferActivationTime: &offerActivationTime,
		//OfferExpiryTime: &OfferExpiryTime,
	}

	jsonValue, _ = json.Marshal(limtOffer)

	req, _ = http.NewRequest(http.MethodPost, "/v1/create_limit_offer", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// case 5 : offer_expiry_time field should be greater than offer_activation_time
	today := time.Now().UTC()
	OfferExpiryTime = today.AddDate(0, 0, -1)
	limtOffer = models.LimitOffer{
		ID:                  "74edf2ad-7ee8-45f7-a7d2-7a3287c33ffe",
		AccountID:           &accountID,
		LimitType:           &accountLimit,
		NewLimit:            &newLimit,
		OfferActivationTime: &offerActivationTime,
		OfferExpiryTime:     &OfferExpiryTime,
	}

	jsonValue, _ = json.Marshal(limtOffer)

	req, _ = http.NewRequest(http.MethodPost, "/v1/create_limit_offer", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestValidateListActiveLimitOffersRequestInput(t *testing.T) {
	// init logging client
	utils.InitLogClient()

	// case 1 : accountID missing

	//accountID := "f83513e1-f0cb-4a49-85e4-8e9ddb1f3417"
	activeLimitOffer := models.ActiveLimitOffer{}
	jsonValue, _ := json.Marshal(activeLimitOffer)

	w := httptest.NewRecorder()
	_, e := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/v1/list_active_limit_offers", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// case 2 : invalid uuid as accountID

	accountID := "f83513e1-f0cb-4a49-85e4-"
	activeLimitOffer = models.ActiveLimitOffer{
		AccountID: accountID,
	}
	jsonValue, _ = json.Marshal(activeLimitOffer)

	req, _ = http.NewRequest(http.MethodPost, "/v1/list_active_limit_offers", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestValidateUpdateLimitOfferStatusRequestInput(t *testing.T) {
	// init logging client
	utils.InitLogClient()

	// case 1 : limit_offer_id missing

	limitOfferID := "f83513e1-f0cb-4a49-85e4-8e9ddb1f3417"
	updateLimitOfferStatus := models.UpdateLimitOfferStatus{
		Status: string(models.Accepted),
	}
	jsonValue, _ := json.Marshal(updateLimitOfferStatus)

	w := httptest.NewRecorder()
	_, e := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/v1/update_limit_offer_status", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// case 2 : some status is missing
	updateLimitOfferStatus = models.UpdateLimitOfferStatus{
		LimitOfferID: limitOfferID,
	}
	jsonValue, _ = json.Marshal(updateLimitOfferStatus)

	req, _ = http.NewRequest(http.MethodPost, "/v1/update_limit_offer_status", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// case 2 : some random status is sent
	updateLimitOfferStatus = models.UpdateLimitOfferStatus{
		LimitOfferID: limitOfferID,
		Status:       "some-non-sense-status",
	}
	jsonValue, _ = json.Marshal(updateLimitOfferStatus)

	req, _ = http.NewRequest(http.MethodPost, "/v1/update_limit_offer_status", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}
