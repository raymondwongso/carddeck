package rest

import "github.com/raymondwongso/carddeck/modules/carddeck/entity"

// DrawCardResponse defines custom response for GET /decks/{id}/cards
type DrawCardResponse struct {
	Cards *entity.Cards `json:"cards"`
}
