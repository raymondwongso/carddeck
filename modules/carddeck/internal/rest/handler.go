// Package rest contains implementation for REST API handler for carddeck module
package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
	"github.com/rs/zerolog/log"
)

// Service defines interfaces for carddeck usecases
type Service interface {
	CreateDeck(ctx context.Context, shuffled bool, cardCodes []string) (*entity.Deck, error)
}

// Handler defines REST API Handler for card deck
type Handler struct {
	svc Service
}

// NewHandler creates new REST API handler.
func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// @summary	Create new deck
// @tags		carddeck
// @produce	json
// @param		shuffled	query	boolean	false	"Specify whether newly created deck is shuffled or not"
// @param		cards		query	string	false	"Specify cards used in this newly created deck"
// @router		/decks [get]
func (h *Handler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	shuffledParam := r.URL.Query().Get("shuffled")
	cardsParam := r.URL.Query().Get("cards")

	var (
		shuffled  = false
		cardCodes []string
	)

	if cardsParam != "" {
		cardCodes = strings.Split(cardsParam, ",")
	}

	if shuffledParam != "" {
		var parseErr error
		shuffled, parseErr = strconv.ParseBool(shuffledParam)
		if parseErr != nil {
			log.Error().Err(parseErr).Msg("[GET /decks] error parsing shuffled parameter")
			err := entity.NewError(entity.ErrParamInvalid, entity.ErrMsgParamInvalid)
			err.AddDetail(entity.NewErrorDetail("shuffled", "shuffled parameter is invalid"))
			handleError(w, err, http.StatusBadRequest)
			return
		}
	}

	deck, err := h.svc.CreateDeck(r.Context(), shuffled, cardCodes)
	if err != nil {
		log.Error().Err(err).Msg("[GET /decks] error creating deck")

		if perr, ok := err.(*entity.Error); ok {
			switch perr.Code {
			case entity.ErrCardCodeInvalid:
				handleError(w, perr, http.StatusUnprocessableEntity)
			default:
				handleError(w, perr, http.StatusInternalServerError)
			}
		} else {
			// error is not in custom error, assume unknown error
			handleError(
				w,
				entity.NewError(entity.ErrInternal, entity.ErrMsgInternal),
				http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&deck); err != nil {
		log.Error().Err(err).Msg("[GET /decks] error encoding response")
		http.Error(w, entity.ErrMsgInternal, http.StatusInternalServerError)
	}
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

func handleError(w http.ResponseWriter, err any, code int) {
	w.WriteHeader(code)

	// if encoding still failed, fallback to default internal error
	if err := json.NewEncoder(w).Encode(&err); err != nil {
		log.Error().Err(err).Msg("[GET /decks] error encoding response")
		http.Error(w, entity.ErrMsgInternal, http.StatusInternalServerError)
	}
}
