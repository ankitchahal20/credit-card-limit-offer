package constants

const (
	ForwardSlash = "/"
	CreateAccount = "create_account"
	GetAccount = "get_account"
	CreateLimitOffer = "create_limit_offer"
	ListActiveLimitOffers = "list_active_limit_offers"
	UpdateLimitOfferStatus = "update_limit_offer_status"
	AccountID = "account_id"
	Colon = ":"
	EmptyString = ""

	Version = "v1"

	TransactionID = "transaction-id"
	InvalidBody   = "invalid value for body"
	InvalidAccountID  = "invalid value for accountID"
	InvalidOfferLimitID  = "invalid value for offer limit id"
	InvalidBodyCreateAccount = "invalid create account request body"
	InvalidBodyCreateLimitOffer = "invalid create limit offer request body"
	InvalidBodyUpdateLimitOfferStatus = "invalid update limit offer status request body"

	//http
	Accept          = "Accept"
	ContentType     = "Content-Type"
	Authorization   = "Authorization"
	ApplicationJSON = "application/json"
)
