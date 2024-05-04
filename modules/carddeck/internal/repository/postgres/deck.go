package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
)

// Deck defines deck repository
type Deck struct {
	db *sqlx.DB
}

// NewDeck returns new deck repository
func NewDeck(db *sqlx.DB) *Deck {
	return &Deck{db: db}
}

func (d *Deck) Insert(ctx context.Context, deck *entity.Deck) (*entity.Deck, error) {
	panic("not implemented")
}

func (d *Deck) GetByID(ctx context.Context, id string) (*entity.Deck, error) {
	panic("not implemented")
}

func (d *Deck) DrawCards(ctx context.Context, id string, count int64) (*entity.Cards, error) {
	panic("not implemented")
}
