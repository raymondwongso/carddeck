package rest

import "github.com/raymondwongso/carddeck/modules/carddeck/entity"

// CreateDeckResponse contains simplified deck information, only showing the ID, shuffled and remaining fields.
type CreateDeckResponse struct {
	ID        string `json:"id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int64  `json:"remaining"`
}

// DrawCardResponse defines custom response for GET /decks/{id}/cards
type DrawCardResponse struct {
	Cards *entity.Cards `json:"cards"`
}
