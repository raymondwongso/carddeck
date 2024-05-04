package rest_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
	"github.com/raymondwongso/carddeck/modules/carddeck/internal/rest"

	"github.com/golang/mock/gomock"
	mock_rest "github.com/raymondwongso/carddeck/test/mock/modules/carddeck/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	defaultTime = time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC)
	defaultDeck = func() *entity.Deck {
		deck := entity.NewDeck(false, &entity.Cards{
			{Val: "ACE", Suit: "SPADE", Code: "AS"},
			{Val: "2", Suit: "SPADE", Code: "2S"},
			{Val: "3", Suit: "SPADE", Code: "3S"},
		})
		deck.ID = "some-uuid-abc-def"
		deck.CreatedAt = defaultTime
		deck.UpdatedAt = defaultTime

		return deck
	}()
	defaultCards = entity.Cards{
		{Val: "ACE", Suit: "SPADE", Code: "AS"},
		{Val: "2", Suit: "SPADE", Code: "2S"},
	}
	defaultDrawCardResponse = rest.DrawCardResponse{
		Cards: &defaultCards,
	}
)

type HandlerTestSuite struct {
	suite.Suite
	svc *mock_rest.MockService
}

func (s *HandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.svc = mock_rest.NewMockService(ctrl)
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) TestCreateDeck() {
	s.Run("success - without shuffled and cards parameter", func() {
		r := httptest.NewRequest(http.MethodPost, "http://localhost", nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().CreateDeck(r.Context(), false, nil).Return(defaultDeck, nil)

		h := rest.NewHandler(s.svc)
		h.CreateDeck(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusCreated, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)
		expected, err := json.Marshal(&defaultDeck)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("success - with shuffled and without cards parameter", func() {
		r := httptest.NewRequest(http.MethodPost, "http://localhost?shuffled=true", nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().CreateDeck(r.Context(), true, nil).Return(defaultDeck, nil)

		h := rest.NewHandler(s.svc)
		h.CreateDeck(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusCreated, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)
		expected, err := json.Marshal(&defaultDeck)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("success - with shuffled and cards parameter", func() {
		r := httptest.NewRequest(http.MethodPost, "http://localhost?shuffled=true&cards=AS,AD,AC,AH", nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().CreateDeck(r.Context(), true, []string{"AS", "AD", "AC", "AH"}).Return(defaultDeck, nil)

		h := rest.NewHandler(s.svc)
		h.CreateDeck(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusCreated, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)
		expected, err := json.Marshal(&defaultDeck)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("failed - shuffled is invalid", func() {
		r := httptest.NewRequest(http.MethodPost, "http://localhost?shuffled=NOT_BOOL_PARSABLE", nil)
		w := httptest.NewRecorder()

		h := rest.NewHandler(s.svc)
		h.CreateDeck(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusBadRequest, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expectedError := entity.NewError(entity.ErrParamInvalid, entity.ErrMsgParamInvalid)
		expectedError.AddDetail(entity.NewErrorDetail("shuffled", "shuffled parameter is invalid"))
		expected, err := json.Marshal(&expectedError)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("failed - service layer returns unexpected error", func() {
		r := httptest.NewRequest(http.MethodPost, "http://localhost", nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().CreateDeck(r.Context(), false, nil).Return(nil, errors.New("unknown error"))

		h := rest.NewHandler(s.svc)
		h.CreateDeck(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusInternalServerError, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expectedError := entity.NewError(entity.ErrInternal, entity.ErrMsgInternal)
		expected, err := json.Marshal(&expectedError)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("failed - service layer returns invalid card codes", func() {
		r := httptest.NewRequest(http.MethodPost, "http://localhost", nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().CreateDeck(r.Context(), false, nil).Return(nil, entity.NewError(entity.ErrCardCodeInvalid, entity.ErrMsgCardCodeInvalid))

		h := rest.NewHandler(s.svc)
		h.CreateDeck(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusUnprocessableEntity, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expectedError := entity.NewError(entity.ErrCardCodeInvalid, entity.ErrMsgCardCodeInvalid)
		expected, err := json.Marshal(&expectedError)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})
}

func (s *HandlerTestSuite) TestGetDeck() {
	tempID := "3cdc5e5a-8f56-4f70-91e6-bd564d04ce79"

	s.Run("success", func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/decks/%s", tempID), nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().GetDeck(r.Context(), tempID).Return(defaultDeck, nil)

		h := rest.NewHandler(s.svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /decks/{id}", h.GetDeck)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		mux.ServeHTTP(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusOK, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)
		expected, err := json.Marshal(&defaultDeck)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("failed - service layer returns unexpected error", func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/decks/%s", tempID), nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().GetDeck(r.Context(), tempID).Return(nil, errors.New("unknown error"))

		h := rest.NewHandler(s.svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /decks/{id}", h.GetDeck)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		mux.ServeHTTP(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusInternalServerError, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expectedError := entity.NewError(entity.ErrInternal, entity.ErrMsgInternal)
		expected, err := json.Marshal(&expectedError)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("failed - deck not found", func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/decks/%s", tempID), nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().GetDeck(r.Context(), tempID).Return(nil, entity.NewError(entity.ErrDeckNotFound, entity.ErrMsgDeckNotFound))

		h := rest.NewHandler(s.svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /decks/{id}", h.GetDeck)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		mux.ServeHTTP(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusNotFound, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expectedError := entity.NewError(entity.ErrDeckNotFound, entity.ErrMsgDeckNotFound)
		expected, err := json.Marshal(&expectedError)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})
}

func (s *HandlerTestSuite) TestDrawCards() {
	tempID := "3cdc5e5a-8f56-4f70-91e6-bd564d04ce79"
	var tempCount int64 = 2

	s.Run("success", func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/decks/%s/cards?count=%d", tempID, tempCount), nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().DrawCards(r.Context(), tempID, tempCount).Return(&defaultCards, nil)

		h := rest.NewHandler(s.svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /decks/{id}/cards", h.DrawCards)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		mux.ServeHTTP(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusOK, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expected, err := json.Marshal(&defaultDrawCardResponse)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("failed - service layer returns unexpected error", func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/decks/%s/cards?count=%d", tempID, tempCount), nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().DrawCards(r.Context(), tempID, tempCount).Return(nil, errors.New("unknown error"))

		h := rest.NewHandler(s.svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /decks/{id}/cards", h.DrawCards)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		mux.ServeHTTP(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusInternalServerError, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expectedError := entity.NewError(entity.ErrInternal, entity.ErrMsgInternal)
		expected, err := json.Marshal(&expectedError)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("failed - deck not found", func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/decks/%s/cards?count=%d", tempID, tempCount), nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().DrawCards(r.Context(), tempID, tempCount).Return(nil, entity.NewError(entity.ErrDeckNotFound, entity.ErrMsgDeckNotFound))

		h := rest.NewHandler(s.svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /decks/{id}/cards", h.DrawCards)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		mux.ServeHTTP(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusNotFound, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expectedError := entity.NewError(entity.ErrDeckNotFound, entity.ErrMsgDeckNotFound)
		expected, err := json.Marshal(&expectedError)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("failed - service layer returns deck insufficient error", func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/decks/%s/cards?count=%d", tempID, tempCount), nil)
		w := httptest.NewRecorder()

		s.svc.EXPECT().DrawCards(r.Context(), tempID, tempCount).Return(nil, entity.NewError(entity.ErrDeckCardInsufficient, entity.ErrMsgDeckCardInsufficient))

		h := rest.NewHandler(s.svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /decks/{id}/cards", h.DrawCards)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		mux.ServeHTTP(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusUnprocessableEntity, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expectedError := entity.NewError(entity.ErrDeckCardInsufficient, entity.ErrMsgDeckCardInsufficient)
		expected, err := json.Marshal(&expectedError)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})

	s.Run("failed - count parameter invalid", func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/decks/%s/cards?count=%s", tempID, "not_a_number"), nil)
		w := httptest.NewRecorder()

		h := rest.NewHandler(s.svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /decks/{id}/cards", h.DrawCards)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		mux.ServeHTTP(w, r)
		response := w.Result()

		assert.Equal(s.T(), http.StatusBadRequest, response.StatusCode)

		rawResponseBody, err := io.ReadAll(response.Body)
		assert.NoError(s.T(), err)

		expectedError := entity.NewError(entity.ErrParamInvalid, entity.ErrMsgParamInvalid)
		expectedError.AddDetail(entity.NewErrorDetail("count", "count parameter is invalid"))
		expected, err := json.Marshal(&expectedError)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), string(expected), strings.TrimSuffix(string(rawResponseBody), "\n"))
	})
}
