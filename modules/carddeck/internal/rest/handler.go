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
	GetDeck(ctx context.Context, id string) (*entity.Deck, error)
	DrawCards(ctx context.Context, id string, n int64) (*entity.Cards, error)
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
			log.Error().Err(parseErr).Msg("[POST /decks] error parsing shuffled parameter")
			err := entity.NewError(entity.ErrParamInvalid, entity.ErrMsgParamInvalid)
			err.AddDetail(entity.NewErrorDetail("shuffled", "shuffled parameter is invalid"))
			handleError(w, err, http.StatusBadRequest)
			return
		}
	}

	deck, err := h.svc.CreateDeck(r.Context(), shuffled, cardCodes)
	if err != nil {
		log.Error().Err(err).Msg("[POST /decks] error creating deck")

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

	resp := CreateDeckResponse{
		ID:        deck.ID,
		Shuffled:  deck.Shuffled,
		Remaining: int64(deck.Remaining()),
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Error().Err(err).Msg("[POST /decks] error encoding response")
		http.Error(w, entity.ErrMsgInternal, http.StatusInternalServerError)
	}
}

// @summary	"Open" a new deck, or get deck by specific ID
// @tags		carddeck
// @produce	json
// @param		id	path	string	true	"ID of the deck"
// @router		/decks/{id} [get]
func (h *Handler) GetDeck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deck, err := h.svc.GetDeck(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("[GET /decks/{id}] error getting deck")

		if perr, ok := err.(*entity.Error); ok {
			switch perr.Code {
			case entity.ErrDeckNotFound:
				handleError(w, perr, http.StatusNotFound)
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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&deck); err != nil {
		log.Error().Err(err).Msg("[GET /decks/{id}] error encoding response")
		http.Error(w, entity.ErrMsgInternal, http.StatusInternalServerError)
	}
}

// @summary	Draw cards from specific deck
// @tags		carddeck
// @produce	json
// @param		id		path	string	true	"ID of the deck"
// @param		count	query	integer	true	"Number of cards to withdraw"
// @router		/decks/{id}/cards [get]
func (h *Handler) DrawCards(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	countParam := r.URL.Query().Get("count")
	count, err := strconv.ParseInt(countParam, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("[GET /decks/{id}/cards] error parsing count parameter")
		err := entity.NewError(entity.ErrParamInvalid, entity.ErrMsgParamInvalid)
		err.AddDetail(entity.NewErrorDetail("count", "count parameter is invalid"))
		handleError(w, err, http.StatusBadRequest)
		return
	}

	cards, err := h.svc.DrawCards(r.Context(), id, count)
	if err != nil {
		log.Error().Err(err).Msg("[GET /decks/{id}/cards] error drawing cards")

		if perr, ok := err.(*entity.Error); ok {
			switch perr.Code {
			case entity.ErrDeckNotFound:
				handleError(w, perr, http.StatusNotFound)
			case entity.ErrDeckCardInsufficient:
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

	resp := DrawCardResponse{
		Cards: cards,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Error().Err(err).Msg("[GET /decks/{id}/cards] error encoding response")
		http.Error(w, entity.ErrMsgInternal, http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, err any, code int) {
	w.WriteHeader(code)

	// if encoding still failed, fallback to default internal error
	if err := json.NewEncoder(w).Encode(&err); err != nil {
		log.Error().Err(err).Msg("[GET /decks] error encoding response")
		http.Error(w, entity.ErrMsgInternal, http.StatusInternalServerError)
	}
}
