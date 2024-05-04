// Package rest contains implementation for REST API handler for carddeck module
package rest

import (
	"context"
	"net/http"

	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
)

// Service defines interfaces for carddeck usecases
type Service interface {
	CreateDeck(ctx context.Context, shuffled bool, cardCodes []string) (*entity.Deck, error)
}

type Handler struct{}

// NewHandler creates new REST API handler.
func NewHandler() *Handler {
	return &Handler{}
}

// @summary	Create new deck
// @tags		carddeck
// @produce	json
// @param		shuffled	query	boolean	false	"Specify whether newly created deck is shuffled or not"
// @param		cards		query	string	false	"Specify cards used in this newly created deck"
// @router		/decks [get]
func (h *Handler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// @summary	"Open" a new deck, or get deck by specific ID
// @tags		carddeck
// @produce	json
// @param		id	path	string	true	"ID of the deck"
// @router		/decks/{id} [get]
func (h *Handler) GetDeck(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// @summary	Draw cards from specific deck
// @tags		carddeck
// @produce	json
// @param		id		path	string	true	"ID of the deck"
// @param		count	query	integer	true	"Number of cards to withdraw"
// @router		/decks/{id}/cards [get]
func (h *Handler) DrawCards(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
