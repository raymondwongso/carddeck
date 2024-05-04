package service_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
	"github.com/raymondwongso/carddeck/modules/carddeck/internal/service"
	mock_service "github.com/raymondwongso/carddeck/test/mock/modules/carddeck/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	defaultTime  = time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC)
	defaultCards = entity.Cards{
		{Val: "ACE", Suit: "SPADE", Code: "AS"},
		{Val: "2", Suit: "SPADE", Code: "2S"},
		{Val: "3", Suit: "SPADE", Code: "3S"},
	}
	defaultDeck = func() *entity.Deck {
		deck := entity.NewDeck(false, &defaultCards)
		deck.ID = "some-uuid-abc-def"
		deck.CreatedAt = defaultTime
		deck.UpdatedAt = defaultTime

		return deck
	}()
	shuffledDeck = func() *entity.Deck {
		deck := entity.NewDeck(false, &entity.Cards{
			{Val: "3", Suit: "SPADE", Code: "3S"},
			{Val: "ACE", Suit: "SPADE", Code: "AS"},
			{Val: "2", Suit: "SPADE", Code: "2S"},
		})
		deck.ID = "some-uuid-abc-def"
		deck.CreatedAt = defaultTime
		deck.UpdatedAt = defaultTime

		return deck
	}()
	shuffledCards = entity.Cards{
		{Val: "3", Suit: "SPADE", Code: "3S"},
		{Val: "ACE", Suit: "SPADE", Code: "AS"},
		{Val: "2", Suit: "SPADE", Code: "2S"},
	}
)

type ServiceTestSuite struct {
	suite.Suite
	deckRepo      *mock_service.MockDeckRepository
	randGenerator func() *rand.Rand
	cardShuffler  func(r *rand.Rand, cards []*entity.Card) []*entity.Card
}

func (s *ServiceTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.deckRepo = mock_service.NewMockDeckRepository(ctrl)
	s.randGenerator = func() *rand.Rand {
		return rand.New(rand.NewSource(defaultTime.Unix()))
	}
	s.cardShuffler = func(_ *rand.Rand, _ []*entity.Card) []*entity.Card {
		return shuffledCards
	}
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestCreateDeck() {
	ctx := context.Background()

	s.Run("success - no shuffle, no cards supplied", func() {
		s.deckRepo.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, deck *entity.Deck) (*entity.Deck, error) {
				assert.Equal(s.T(), false, deck.Shuffled)
				assert.Equal(s.T(), (*entity.Cards)(&service.CardArray), deck.Cards)

				return defaultDeck, nil
			})

		svc := service.New(s.deckRepo, s.randGenerator, s.cardShuffler)
		deck, err := svc.CreateDeck(ctx, false, nil)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), defaultDeck, deck)
	})

	s.Run("success - with shuffle, with cards supplied", func() {
		s.deckRepo.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, deck *entity.Deck) (*entity.Deck, error) {
				assert.Equal(s.T(), true, deck.Shuffled)
				assert.Equal(s.T(), (*entity.Cards)(&shuffledCards), deck.Cards)

				return defaultDeck, nil
			})

		svc := service.New(s.deckRepo, s.randGenerator, s.cardShuffler)
		deck, err := svc.CreateDeck(ctx, true, []string{"AS", "2S", "3S"})
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), shuffledDeck.ID, deck.ID)
		assert.Equal(s.T(), len(*shuffledDeck.Cards), len(*deck.Cards))
		assert.Equal(s.T(), shuffledDeck.CreatedAt, deck.UpdatedAt)
		assert.Equal(s.T(), shuffledDeck.Remaining(), deck.Remaining())
	})

	s.Run("failed - card codes invalid", func() {
		svc := service.New(s.deckRepo, s.randGenerator, s.cardShuffler)
		deck, err := svc.CreateDeck(ctx, false, []string{"XX", "YY", "ZZ"})
		assert.Nil(s.T(), deck)
		assert.Error(s.T(), err)

		perr, ok := err.(*entity.Error)
		assert.True(s.T(), ok)

		assert.Equal(s.T(), entity.ErrCardCodeInvalid, perr.Code)
		assert.Equal(s.T(), entity.ErrMsgCardCodeInvalid, perr.Message)
	})

	s.Run("failed - unexpected error", func() {
		s.deckRepo.EXPECT().Insert(ctx, gomock.Any()).Return(nil, errors.New("some error"))

		svc := service.New(s.deckRepo, s.randGenerator, s.cardShuffler)
		deck, err := svc.CreateDeck(ctx, false, nil)
		assert.Nil(s.T(), deck)
		assert.Error(s.T(), err)
	})
}

func (s *ServiceTestSuite) TestGetDeck() {
	ctx := context.Background()
	id := "some_id"

	s.Run("success", func() {
		s.deckRepo.EXPECT().GetByID(ctx, id).Return(defaultDeck, nil)

		svc := service.New(s.deckRepo, s.randGenerator, s.cardShuffler)
		deck, err := svc.GetDeck(ctx, id)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), defaultDeck, deck)
	})

	s.Run("failed - id empty", func() {
		svc := service.New(s.deckRepo, s.randGenerator, s.cardShuffler)
		deck, err := svc.GetDeck(ctx, "")
		assert.Nil(s.T(), deck)
		assert.Error(s.T(), err)

		perr, ok := err.(*entity.Error)
		assert.True(s.T(), ok)
		assert.Equal(s.T(), perr.Code, entity.ErrParamInvalid)
		assert.Equal(s.T(), perr.Message, entity.ErrMsgParamInvalid)
	})
}

func (s *ServiceTestSuite) TestDrawCards() {
	ctx := context.Background()
	id := "some_id"
	var n int64 = 2

	s.Run("success", func() {
		s.deckRepo.EXPECT().DrawCards(ctx, id, n).Return(&defaultCards, nil)

		svc := service.New(s.deckRepo, s.randGenerator, s.cardShuffler)
		cards, err := svc.DrawCards(ctx, id, n)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), &defaultCards, cards)
	})

	s.Run("failed - id empty", func() {
		svc := service.New(s.deckRepo, s.randGenerator, s.cardShuffler)
		cards, err := svc.DrawCards(ctx, "", n)
		assert.Nil(s.T(), cards)
		assert.Error(s.T(), err)

		perr, ok := err.(*entity.Error)
		assert.True(s.T(), ok)
		assert.Equal(s.T(), perr.Code, entity.ErrParamInvalid)
		assert.Equal(s.T(), perr.Message, entity.ErrMsgParamInvalid)
	})

	s.Run("failed - count is zero", func() {
		svc := service.New(s.deckRepo, s.randGenerator, s.cardShuffler)
		cards, err := svc.DrawCards(ctx, "", 0)
		assert.Nil(s.T(), cards)
		assert.Error(s.T(), err)

		perr, ok := err.(*entity.Error)
		assert.True(s.T(), ok)
		assert.Equal(s.T(), perr.Code, entity.ErrParamInvalid)
		assert.Equal(s.T(), perr.Message, entity.ErrMsgParamInvalid)
	})
}
