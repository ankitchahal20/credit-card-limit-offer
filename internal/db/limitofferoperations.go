package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ankit/project/credit-card-offer-limit/internal/constants"
	"github.com/ankit/project/credit-card-offer-limit/internal/limitoffererror"
	"github.com/ankit/project/credit-card-offer-limit/internal/models"
	"github.com/ankit/project/credit-card-offer-limit/internal/utils"
	"github.com/gin-gonic/gin"
)

func (p postgres) IsLimitOfferExists(ctx *gin.Context, limitOffer models.LimitOffer) (bool, string, *limitoffererror.CreditCardError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	query := `
		SELECT id
		FROM limit_offer
		WHERE account_id = $1 AND limit_type = $2 AND status = $3
		LIMIT 1`

	var offerLimitId sql.NullString
	//var offerStatusId sql.NullString
	err := p.db.QueryRow(query, limitOffer.AccountID, limitOffer.LimitType, models.Pending).Scan(&offerLimitId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}

		utils.Logger.Error(fmt.Sprintf("error querying limit offer existence from db, txid: %v, error: %v", txid, err))
		return false, "", &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "error checking limit offer existence",
			Trace:   txid,
		}
	}

	if offerLimitId.Valid {
		return true, offerLimitId.String, nil
	}

	return false, "", nil
}

func (p postgres) CreateLimitOffer(ctx *gin.Context, limitOffer models.LimitOffer, isLimitOfferExsits bool) *limitoffererror.CreditCardError {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	if isLimitOfferExsits {
		fmt.Println("limitOffer.NewLimit, limitOffer.AccountID, limitOffer.LimitType :", *limitOffer.NewLimit, ":", *limitOffer.AccountID, ":", *limitOffer.LimitType)
		_, err := p.db.Exec("UPDATE limit_offer SET new_limit = $1 WHERE account_id = $2 AND limit_type = $3", *limitOffer.NewLimit, *limitOffer.AccountID, *limitOffer.LimitType)
		fmt.Println("err 3 ", err)
		if err != nil {
			log.Println("error updating limit offer status:", err)
			return &limitoffererror.CreditCardError{
				Code:    http.StatusInternalServerError,
				Message: "error while updating limit offer status",
				Trace:   txid,
			}
		}
	} else {
		query := `
			INSERT INTO limit_offer(id, account_id, limit_type, new_limit, offer_activation_time, offer_expiry_time, status) 
			VALUES($1, $2, $3, $4, $5, $6, $7)`

		_, err := p.db.Exec(query, limitOffer.ID, limitOffer.AccountID, limitOffer.LimitType, limitOffer.NewLimit,
			limitOffer.OfferActivationTime, limitOffer.OfferExpiryTime, limitOffer.Status)

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
				Message: "unable to add offer limit info",
			}
		}
	}
	utils.Logger.Info(fmt.Sprintf("successfully added the offer limit entry in db, txid : %v", txid))
	return nil
}

func (p postgres) ListActiveLimitOffers(ctx *gin.Context, limitOffer models.ActiveLimitOffer) ([]models.LimitOffer, *limitoffererror.CreditCardError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	// Check if the account with the provided account_id exists
	var accountExists bool
	accountCheckQuery := `SELECT EXISTS (SELECT 1 FROM account WHERE account_id = $1)`
	if err := p.db.QueryRow(accountCheckQuery, limitOffer.AccountID).Scan(&accountExists); err != nil {
		return nil, &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "error checking account existence",
			Trace:   txid,
		}
	}

	if !accountExists {
		return nil, &limitoffererror.CreditCardError{
			Code:    http.StatusNotFound,
			Message: "account not found",
			Trace:   txid,
		}
	}

	query := `
		SELECT id, account_id, limit_type, new_limit, offer_activation_time, offer_expiry_time, status
		FROM limit_offer
		WHERE account_id = $1 AND status = $2 AND offer_activation_time <= $3 AND offer_expiry_time >= $4`

	activeOffers := []models.LimitOffer{}
	rows, err := p.db.Query(query, limitOffer.AccountID, models.Pending, limitOffer.ActiveDate, limitOffer.ActiveDate)
	if err != nil {
		// Handle the error if the query fails
		return nil, &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "unable to retrieve active limit offers",
			Trace:   txid,
		}
	}
	defer rows.Close()

	for rows.Next() {
		var offer models.LimitOffer
		err := rows.Scan(
			&offer.ID, &offer.AccountID, &offer.LimitType, &offer.NewLimit,
			&offer.OfferActivationTime, &offer.OfferExpiryTime, &offer.Status,
		)
		if err != nil {
			// Handle the error if scanning fails
			return nil, &limitoffererror.CreditCardError{
				Code:    http.StatusInternalServerError,
				Message: "error scanning limit offer rows",
				Trace:   txid,
			}
		}
		activeOffers = append(activeOffers, offer)
	}

	return activeOffers, nil
}

func (p postgres) UpdateLimitOfferStatus(ctx *gin.Context, updateLimitOfferStatus models.UpdateLimitOfferStatus) *limitoffererror.CreditCardError {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	// Check if the account with the provided offer_limit_id exists
	var offerLimitExists bool
	offerLimitCheckQuery := `SELECT EXISTS (SELECT 1 FROM limit_offer WHERE id = $1)`
	if err := p.db.QueryRow(offerLimitCheckQuery, updateLimitOfferStatus.LimitOfferID).Scan(&offerLimitExists); err != nil {
		return &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "error checking account existence",
			Trace:   txid,
		}
	}

	if !offerLimitExists {
		return &limitoffererror.CreditCardError{
			Code:    http.StatusNotFound,
			Message: "offer limit not found",
			Trace:   txid,
		}
	}

	tx, err := p.db.Begin()
	fmt.Println("err 1 ", err)
	if err != nil {
		return &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "unable to begin transaction",
			Trace:   txid,
		}
	}
	defer tx.Rollback()

	var limitOffer models.LimitOffer
	query := `SELECT * FROM limit_offer WHERE id = $1`
	err = tx.QueryRow(query, updateLimitOfferStatus.LimitOfferID).Scan(&limitOffer.ID,
		&limitOffer.AccountID,
		&limitOffer.LimitType,
		&limitOffer.NewLimit,
		&limitOffer.OfferActivationTime,
		&limitOffer.OfferExpiryTime,
		&limitOffer.Status)
	fmt.Println("err 2 ", err)
	if err != nil {
		return &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "error while fetching limit offer details",
			Trace:   txid,
		}
	}

	isActiveOffer := limitOffer.OfferActivationTime.Before(time.Now().UTC()) && limitOffer.OfferExpiryTime.After(time.Now().UTC())
	if !isActiveOffer {
		// if not in range
		return &limitoffererror.CreditCardError{
			Code:    http.StatusNotFound,
			Message: "limit offer already expired",
			Trace:   txid,
		}
	}

	// update the status to ACCEPTED/REJECTED
	limitOffer.Status = models.OfferStatus(updateLimitOfferStatus.Status)
	_, err = tx.Exec("UPDATE limit_offer SET status = $1 WHERE id = $2", limitOffer.Status, limitOffer.ID)
	fmt.Println("err 3 ", err)
	if err != nil {
		log.Println("error updating limit offer status:", err)
		return &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "error while updating limit offer status",
			Trace:   txid,
		}
	}

	switch limitOffer.Status {
	case models.Rejected:
		// if status is REJECTED, your work is done, no updation required in db.
	case models.Accepted:
		// if status is ACCEPTED, update limit values (current and last), as well as limit update date in the account object.
		var accountInfo models.Account
		accountQuery := `SELECT * FROM account WHERE account_id = $1`
		err = tx.QueryRow(accountQuery, limitOffer.AccountID).Scan(&accountInfo.AccountID,
			&accountInfo.CustomerID,
			&accountInfo.AccountLimit,
			&accountInfo.PerTransactionLimit,
			&accountInfo.LastAccountLimit,
			&accountInfo.LastPerTransactionLimit,
			&accountInfo.AccountLimitUpdateTime,
			&accountInfo.PerTransactionLimitUpdateTime)
		fmt.Println("err 4 ", err)
		if err != nil {
			log.Println("error fetching account:", err)
			return &limitoffererror.CreditCardError{
				Code:    http.StatusInternalServerError,
				Message: "error while reteriving get account info",
				Trace:   txid,
			}
		}

		if *limitOffer.LimitType == models.AccountLimit {
			accountInfo.LastAccountLimit = accountInfo.AccountLimit
			accountInfo.AccountLimit = limitOffer.NewLimit
			accountInfo.AccountLimitUpdateTime = time.Now().UTC()

			// update the db
			_, err = tx.Exec("UPDATE account SET last_account_limit = $1, account_limit = $2, account_limit_update_time = $3 WHERE account_id = $4",
				accountInfo.LastAccountLimit, accountInfo.AccountLimit, accountInfo.AccountLimitUpdateTime, accountInfo.AccountID)
			if err != nil {
				log.Println("error updating account:", err)
				return &limitoffererror.CreditCardError{
					Code:    http.StatusInternalServerError,
					Message: "unable to update the account limit info in db",
					Trace:   txid,
				}
			}

		} else if *limitOffer.LimitType == models.PerTransactionLimit {
			accountInfo.LastPerTransactionLimit = accountInfo.PerTransactionLimit
			accountInfo.PerTransactionLimit = limitOffer.NewLimit
			accountInfo.PerTransactionLimitUpdateTime = time.Now().UTC()
			_, err = tx.Exec("UPDATE account SET last_per_transaction_limit = $1, per_transaction_limit = $2, per_transaction_limit_update_time = $3 WHERE account_id = $4",
				accountInfo.LastPerTransactionLimit, accountInfo.PerTransactionLimit, accountInfo.PerTransactionLimitUpdateTime, accountInfo.AccountID)
			if err != nil {
				log.Println("error updating account:", err)
				return &limitoffererror.CreditCardError{
					Code:    http.StatusInternalServerError,
					Message: "unable to update the account limit info in db",
					Trace:   txid,
				}
			}
		}

	default:
		return &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "received status not supported",
			Trace:   txid,
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("error committing transaction:", err)
		return &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "unable to commit changes in db",
			Trace:   txid,
		}
	}
	return nil
}

func (p postgres) GetLimitOffer(ctx *gin.Context, offerLimitID string) (models.LimitOffer, *limitoffererror.CreditCardError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	scannedLimitOffer := models.LimitOffer{}

	query := `SELECT * FROM limit_offer WHERE id=$1`
	row := p.db.QueryRow(query, offerLimitID)

	err := row.Scan(
		&scannedLimitOffer.ID,
		&scannedLimitOffer.AccountID,
		&scannedLimitOffer.LimitType,
		&scannedLimitOffer.NewLimit,
		&scannedLimitOffer.OfferActivationTime,
		&scannedLimitOffer.OfferExpiryTime,
		&scannedLimitOffer.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no rows were found
			return scannedLimitOffer, &limitoffererror.CreditCardError{
				Code:    http.StatusNotFound,
				Message: "limit offer not found",
				Trace:   txid,
			}
		}

		utils.Logger.Error(fmt.Sprintf("error while scanning account from db, txid : %v, error: %v", txid, err))
		return scannedLimitOffer, &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "unable to get the limit offer",
			Trace:   txid,
		}
	}

	utils.Logger.Info(fmt.Sprintf("successfully fetched limit offer from db, txid : %v", txid))
	return scannedLimitOffer, nil
}
