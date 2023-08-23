package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/ankit/project/credit-card-offer-limit/internal/constants"
	"github.com/ankit/project/credit-card-offer-limit/internal/limitoffererror"
	"github.com/ankit/project/credit-card-offer-limit/internal/models"
	"github.com/ankit/project/credit-card-offer-limit/internal/utils"
	"github.com/gin-gonic/gin"
)

func (p postgres) CreateAccount(ctx *gin.Context, accountInfo models.Account) *limitoffererror.CreditCardError {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	
	query := `
			INSERT INTO account(account_id, customer_id, account_limit, per_transaction_limit, last_account_limit, 
			last_per_transaction_limit, account_limit_update_time, per_transaction_limit_update_time) 
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	
	_, err := p.db.Exec(query, accountInfo.AccountID, accountInfo.CustomerID, accountInfo.AccountLimit, 
		accountInfo.PerTransactionLimit, accountInfo.LastAccountLimit, accountInfo.LastPerTransactionLimit, 
		accountInfo.AccountLimitUpdateTime, accountInfo.PerTransactionLimitUpdateTime)
	
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("error while running insert query, txid : %v", txid))
		if strings.Contains(err.Error(), "duplicate key value") {
			return &limitoffererror.CreditCardError{
				Trace:   ctx.Request.Header.Get(constants.TransactionID),
				Code:    http.StatusBadRequest,
				Message: "account already added",
			}
		}
		return &limitoffererror.CreditCardError{
			Trace:   txid,
			Code:    http.StatusInternalServerError,
			Message: "unable to add account info",
		}
	}
	utils.Logger.Info(fmt.Sprintf("successfully added the account entry in db, txid : %v", txid))
	return nil
}

func (p postgres) GetAccount(ctx *gin.Context, accountID string) (models.Account, *limitoffererror.CreditCardError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	scannedAccount := models.Account{}

	query := `SELECT * FROM account WHERE account_id=$1`
	row := p.db.QueryRow(query, accountID)

	err := row.Scan(
		&scannedAccount.AccountID,
		&scannedAccount.CustomerID,
		&scannedAccount.AccountLimit,
		&scannedAccount.PerTransactionLimit,
		&scannedAccount.LastAccountLimit,
		&scannedAccount.LastPerTransactionLimit,
		&scannedAccount.AccountLimitUpdateTime,
		&scannedAccount.PerTransactionLimitUpdateTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no rows were found
			return scannedAccount, &limitoffererror.CreditCardError{
				Code:    http.StatusNotFound,
				Message: "account not found",
				Trace:   txid,
			}
		}

		utils.Logger.Error(fmt.Sprintf("error while scanning account from db, txid : %v, error: %v", txid, err))
		return scannedAccount, &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "unable to get the account",
			Trace:   txid,
		}
	}

	utils.Logger.Info(fmt.Sprintf("successfully fetched account from db, txid : %v", txid))
	return scannedAccount, nil
}




