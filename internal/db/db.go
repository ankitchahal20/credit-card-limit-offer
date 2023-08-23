package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/ankit/project/credit-card-offer-limit/internal/config"
	"github.com/ankit/project/credit-card-offer-limit/internal/limitoffererror"
	"github.com/ankit/project/credit-card-offer-limit/internal/models"
	"github.com/gin-gonic/gin"
)

var (
	conn *sql.DB
	once sync.Once
)

type postgres struct{ db *sql.DB }

type CreditCardLimitOfferService interface {
	CreateAccount(*gin.Context, models.Account) *limitoffererror.CreditCardError
	GetAccount(*gin.Context, string) (models.Account, *limitoffererror.CreditCardError)
	CreateLimitOffer(*gin.Context, models.LimitOffer, bool) *limitoffererror.CreditCardError
	ListActiveLimitOffers(*gin.Context, models.ActiveLimitOffer) ([]models.LimitOffer, *limitoffererror.CreditCardError)
	UpdateLimitOfferStatus(*gin.Context, models.UpdateLimitOfferStatus) *limitoffererror.CreditCardError
	IsLimitOfferExists(*gin.Context, models.LimitOffer) (bool, string, *limitoffererror.CreditCardError)
	GetLimitOffer(*gin.Context, string) (models.LimitOffer, *limitoffererror.CreditCardError)
}

func New() (postgres, error) {
	// Initialize the connection only once
	once.Do(func() {
		cfg := config.GetConfig()
		connString := fmt.Sprintf(
			"host=%s dbname=%s password=%s user=%s port=%d",
			cfg.Database.Host, cfg.Database.DBname, cfg.Database.Password,
			cfg.Database.User, cfg.Database.Port,
		)

		var err error
		conn, err = sql.Open("pgx", connString)
		if err != nil {
			log.Fatalf("Unable to connect: %v\n", err)
		}

		log.Println("Connected to database")

		err = conn.Ping()
		if err != nil {
			log.Fatal("Cannot Ping the database")
		}
		log.Println("pinged database")
	})

	return postgres{db: conn}, nil
}
