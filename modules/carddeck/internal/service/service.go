// Package service implements usecases for card deck
package service

import (
	"context"

	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
)

type service struct {
}

// New creates new carddeck service layer (usecase)
func New() *service {
	return &service{}
}

// CreateDeck create deck
func (s *service) CreateDeck(ctx context.Context, shuffled bool, cardCodes []string) (*entity.Deck, error) {
	panic("not implemented")
}

// GetDeck get deck by ID
func (s *service) GetDeck(ctx context.Context, id string) (*entity.Deck, error) {
	panic("not implemented")
}
