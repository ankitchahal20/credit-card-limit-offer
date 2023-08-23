package db

import (
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

func (p postgres) CreateLimitOffer(ctx *gin.Context, limitOffer models.LimitOffer) *limitoffererror.CreditCardError {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	
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
	utils.Logger.Info(fmt.Sprintf("successfully added the offer limit entry in db, txid : %v", txid))
	return nil
}

func (p postgres) ListActiveLimitOffers(ctx *gin.Context, limitOffer models.ActiveLimitOffer) ([]models.LimitOffer, *limitoffererror.CreditCardError) {
	query := `
		SELECT id, account_id, limit_type, new_limit, offer_activation_time, offer_expiry_time, status
		FROM limit_offer
		WHERE account_id = $1 AND status = $2 AND offer_activation_time <= $3 AND offer_expiry_time >= $4`

	rows, err := p.db.Query(query, limitOffer.AccountID, models.Pending, limitOffer.ActiveDate, limitOffer.ActiveDate)
	if err != nil {
		// Handle the error if the query fails
		return nil, &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "unable to retrieve active limit offers",
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
		}
	}
	defer rows.Close()
	activeOffers := []models.LimitOffer{}
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
				Trace:   ctx.Request.Header.Get(constants.TransactionID),
			}
		}
		activeOffers = append(activeOffers, offer)
	}

	return activeOffers, nil
}

func (p postgres) UpdateLimitOfferStatus(ctx *gin.Context, updateLimitOfferStatus models.UpdateLimitOfferStatus) ([]models.LimitOffer, *limitoffererror.CreditCardError){
	txid := ctx.Request.Header.Get(constants.TransactionID)
	tx, err := p.db.Begin()
	fmt.Println("err 1 ", err)
	if err != nil {
		return []models.LimitOffer{},  &limitoffererror.CreditCardError{
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
		&limitOffer.Status,)
	fmt.Println("err 2 ", err)
	if err != nil {
		return []models.LimitOffer{},  &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "error while fetching limit offer details",
			Trace:   txid,
		}
	}

	// update the status to ACCEPTED/REJECTED
	limitOffer.Status = models.OfferStatus(updateLimitOfferStatus.Status)
	_, err = tx.Exec("UPDATE limit_offer SET status = $1 WHERE id = $2", limitOffer.Status, limitOffer.ID)
	fmt.Println("err 3 ", err)
	if err != nil {
		log.Println("error updating limit offer status:", err)
		return []models.LimitOffer{},  &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "error while updating limit offer status",
			Trace:   txid,
		}
	}

	switch limitOffer.Status{
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
				&accountInfo.PerTransactionLimitUpdateTime,)
			fmt.Println("err 4 ", err)
			if err != nil {
				log.Println("error fetching account:", err)
				return []models.LimitOffer{},  &limitoffererror.CreditCardError{
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
					return []models.LimitOffer{},  &limitoffererror.CreditCardError{
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
					return []models.LimitOffer{},  &limitoffererror.CreditCardError{
						Code:    http.StatusInternalServerError,
						Message: "unable to update the account limit info in db",
						Trace:   txid,
					}
				}
			}

	default:
		return []models.LimitOffer{},  &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "received status not supported",
			Trace:   txid,
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("error committing transaction:", err)
		return []models.LimitOffer{},  &limitoffererror.CreditCardError{
			Code:    http.StatusInternalServerError,
			Message: "unable to commit changes in db",
			Trace:   txid,
		}
	}
	return []models.LimitOffer{}, nil
}