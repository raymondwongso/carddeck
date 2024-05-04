// Package carddeck contains implementation for usecases related to French card deck.
package carddeck

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/raymondwongso/carddeck/config"
	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
	"github.com/raymondwongso/carddeck/modules/carddeck/internal/repository/postgres"
	"github.com/raymondwongso/carddeck/modules/carddeck/internal/rest"
	"github.com/raymondwongso/carddeck/modules/carddeck/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib" // driver for postgres
)

// BuildHandler build and returns handler
func BuildHandler(cfg *config.Config) (*rest.Handler, error) {
	connCfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Pass,
		cfg.Postgres.DatabaseName,
	)
	db, err := sqlx.Connect("pgx", connCfg)
	if err != nil {
		return nil, err
	}

	deckRepository := postgres.NewDeck(db)

	randGenerator := func() *rand.Rand {
		return rand.New(rand.NewSource(time.Now().Unix()))
	}

	cardShuffler := func(r *rand.Rand, cards []*entity.Card) []*entity.Card {
		for i := range cards {
			j := r.Intn(i + 1)
			cards[i], cards[j] = cards[j], cards[i]
		}
		return cards
	}

	svc := service.New(deckRepository, randGenerator, cardShuffler)
	return rest.NewHandler(svc), nil
}
