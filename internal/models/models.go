package models

import "time"

type LimitType string

const (
	AccountLimit         LimitType = "ACCOUNT_LIMIT"
	PerTransactionLimit LimitType = "PER_TRANSACTION_LIMIT"
)

type OfferStatus string

const (
	Pending  OfferStatus = "PENDING"
	Accepted OfferStatus = "ACCEPTED"
	Rejected OfferStatus = "REJECTED"
)

type LimitOffer struct {
	ID                	string      `json:"id"`
	AccountID         	string      `json:"account_id"`
	LimitType         	LimitType   `json:"limit_type"`
	NewLimit          	int 		`json:"new_limit"`
	OfferActivationTime time.Time   `json:"offer_activation"`
	OfferExpiryTime   	time.Time   `json:"offer_expiry_time"`
	Status            	OfferStatus `json:"status"`
}

type Account struct {
	AccountID                     string 	`json:"account_id"`
	CustomerID              	  string 	`json:"customer_id"`
	AccountLimit            	  *int 		`json:"account_limit"`
	PerTransactionLimit     	  *int 		`json:"per_transaction_limit"`
	LastAccountLimit        	  *int 		`json:"last_account_limit"`
	LastPerTransactionLimit 	  *int 		`json:"last_per_transaction_limit"`
	AccountLimitUpdateTime 		  time.Time `json:"account_limit_update_time,omitempty"`
	PerTransactionLimitUpdateTime time.Time `json:"per_transaction_limit_update_time,omitempty"`
}

