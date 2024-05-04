// Package service implements usecases for card deck
package service

import (
	"context"
	"math/rand"

	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
)

var (
	CardMapping = map[string]entity.Card{
		"AS":  {Val: "ACE", Suit: "SPADE", Code: "AS"},
		"2S":  {Val: "2", Suit: "SPADE", Code: "2S"},
		"3S":  {Val: "3", Suit: "SPADE", Code: "3S"},
		"4S":  {Val: "4", Suit: "SPADE", Code: "4S"},
		"5S":  {Val: "5", Suit: "SPADE", Code: "5S"},
		"6S":  {Val: "6", Suit: "SPADE", Code: "6S"},
		"7S":  {Val: "7", Suit: "SPADE", Code: "7S"},
		"8S":  {Val: "8", Suit: "SPADE", Code: "8S"},
		"9S":  {Val: "9", Suit: "SPADE", Code: "9S"},
		"10S": {Val: "10", Suit: "SPADE", Code: "10S"},
		"JS":  {Val: "JACK", Suit: "SPADE", Code: "JS"},
		"QS":  {Val: "QUEEN", Suit: "SPADE", Code: "QS"},
		"KS":  {Val: "KING", Suit: "SPADE", Code: "KS"},

		"AD":  {Val: "ACE", Suit: "DIAMOND", Code: "AD"},
		"2D":  {Val: "2", Suit: "DIAMOND", Code: "2D"},
		"3D":  {Val: "3", Suit: "DIAMOND", Code: "3D"},
		"4D":  {Val: "4", Suit: "DIAMOND", Code: "4D"},
		"5D":  {Val: "5", Suit: "DIAMOND", Code: "5D"},
		"6D":  {Val: "6", Suit: "DIAMOND", Code: "6D"},
		"7D":  {Val: "7", Suit: "DIAMOND", Code: "7D"},
		"8D":  {Val: "8", Suit: "DIAMOND", Code: "8D"},
		"9D":  {Val: "9", Suit: "DIAMOND", Code: "9D"},
		"10D": {Val: "10", Suit: "DIAMOND", Code: "10D"},
		"JD":  {Val: "JACK", Suit: "DIAMOND", Code: "JD"},
		"QD":  {Val: "QUEEN", Suit: "DIAMOND", Code: "QD"},
		"KD":  {Val: "KING", Suit: "DIAMOND", Code: "KD"},

		"AC":  {Val: "ACE", Suit: "CLUB", Code: "AC"},
		"2C":  {Val: "2", Suit: "CLUB", Code: "2C"},
		"3C":  {Val: "3", Suit: "CLUB", Code: "3C"},
		"4C":  {Val: "4", Suit: "CLUB", Code: "4C"},
		"5C":  {Val: "5", Suit: "CLUB", Code: "5C"},
		"6C":  {Val: "6", Suit: "CLUB", Code: "6C"},
		"7C":  {Val: "7", Suit: "CLUB", Code: "7C"},
		"8C":  {Val: "8", Suit: "CLUB", Code: "8C"},
		"9C":  {Val: "9", Suit: "CLUB", Code: "9C"},
		"10C": {Val: "10", Suit: "CLUB", Code: "10C"},
		"JC":  {Val: "JACK", Suit: "CLUB", Code: "JC"},
		"QC":  {Val: "QUEEN", Suit: "CLUB", Code: "QC"},
		"KC":  {Val: "KING", Suit: "CLUB", Code: "KC"},

		"AH":  {Val: "ACE", Suit: "HEART", Code: "AH"},
		"2H":  {Val: "2", Suit: "HEART", Code: "2H"},
		"3H":  {Val: "3", Suit: "HEART", Code: "3H"},
		"4H":  {Val: "4", Suit: "HEART", Code: "4H"},
		"5H":  {Val: "5", Suit: "HEART", Code: "5H"},
		"6H":  {Val: "6", Suit: "HEART", Code: "6H"},
		"7H":  {Val: "7", Suit: "HEART", Code: "7H"},
		"8H":  {Val: "8", Suit: "HEART", Code: "8H"},
		"9H":  {Val: "9", Suit: "HEART", Code: "9H"},
		"10H": {Val: "10", Suit: "HEART", Code: "10H"},
		"JH":  {Val: "JACK", Suit: "HEART", Code: "JH"},
		"QH":  {Val: "QUEEN", Suit: "HEART", Code: "QH"},
		"KH":  {Val: "KING", Suit: "HEART", Code: "KH"},
	}

	CardArray = []*entity.Card{
		{Val: "ACE", Suit: "SPADE", Code: "AS"},
		{Val: "2", Suit: "SPADE", Code: "2S"},
		{Val: "3", Suit: "SPADE", Code: "3S"},
		{Val: "4", Suit: "SPADE", Code: "4S"},
		{Val: "5", Suit: "SPADE", Code: "5S"},
		{Val: "6", Suit: "SPADE", Code: "6S"},
		{Val: "7", Suit: "SPADE", Code: "7S"},
		{Val: "8", Suit: "SPADE", Code: "8S"},
		{Val: "9", Suit: "SPADE", Code: "9S"},
		{Val: "10", Suit: "SPADE", Code: "10S"},
		{Val: "JACK", Suit: "SPADE", Code: "JS"},
		{Val: "QUEEN", Suit: "SPADE", Code: "QS"},
		{Val: "KING", Suit: "SPADE", Code: "KS"},

		{Val: "ACE", Suit: "DIAMOND", Code: "AD"},
		{Val: "2", Suit: "DIAMOND", Code: "2D"},
		{Val: "3", Suit: "DIAMOND", Code: "3D"},
		{Val: "4", Suit: "DIAMOND", Code: "4D"},
		{Val: "5", Suit: "DIAMOND", Code: "5D"},
		{Val: "6", Suit: "DIAMOND", Code: "6D"},
		{Val: "7", Suit: "DIAMOND", Code: "7D"},
		{Val: "8", Suit: "DIAMOND", Code: "8D"},
		{Val: "9", Suit: "DIAMOND", Code: "9D"},
		{Val: "10", Suit: "DIAMOND", Code: "10D"},
		{Val: "JACK", Suit: "DIAMOND", Code: "JD"},
		{Val: "QUEEN", Suit: "DIAMOND", Code: "QD"},
		{Val: "KING", Suit: "DIAMOND", Code: "KD"},

		{Val: "ACE", Suit: "CLUB", Code: "AC"},
		{Val: "2", Suit: "CLUB", Code: "2C"},
		{Val: "3", Suit: "CLUB", Code: "3C"},
		{Val: "4", Suit: "CLUB", Code: "4C"},
		{Val: "5", Suit: "CLUB", Code: "5C"},
		{Val: "6", Suit: "CLUB", Code: "6C"},
		{Val: "7", Suit: "CLUB", Code: "7C"},
		{Val: "8", Suit: "CLUB", Code: "8C"},
		{Val: "9", Suit: "CLUB", Code: "9C"},
		{Val: "10", Suit: "CLUB", Code: "10C"},
		{Val: "JACK", Suit: "CLUB", Code: "JC"},
		{Val: "QUEEN", Suit: "CLUB", Code: "QC"},
		{Val: "KING", Suit: "CLUB", Code: "KC"},

		{Val: "ACE", Suit: "HEART", Code: "AH"},
		{Val: "2", Suit: "HEART", Code: "2H"},
		{Val: "3", Suit: "HEART", Code: "3H"},
		{Val: "4", Suit: "HEART", Code: "4H"},
		{Val: "5", Suit: "HEART", Code: "5H"},
		{Val: "6", Suit: "HEART", Code: "6H"},
		{Val: "7", Suit: "HEART", Code: "7H"},
		{Val: "8", Suit: "HEART", Code: "8H"},
		{Val: "9", Suit: "HEART", Code: "9H"},
		{Val: "10", Suit: "HEART", Code: "10H"},
		{Val: "JACK", Suit: "HEART", Code: "JH"},
		{Val: "QUEEN", Suit: "HEART", Code: "QH"},
		{Val: "KING", Suit: "HEART", Code: "KH"},
	}
)

// DeckRepository defines repository for accessing deck data
type DeckRepository interface {
	Insert(ctx context.Context, deck *entity.Deck) (*entity.Deck, error)
	GetByID(ctx context.Context, id string) (*entity.Deck, error)
	DrawCards(ctx context.Context, id string, count int64) (*entity.Cards, error)
}

type Service struct {
	deckRepository DeckRepository
	generateRandom RandomGenerator
	shuffleCard    CardShuffler
}

type RandomGenerator func() *rand.Rand
type CardShuffler func(*rand.Rand, []*entity.Card) []*entity.Card

// New creates new carddeck service layer (usecase)
func New(dr DeckRepository, randGenerator RandomGenerator, cardShuffler CardShuffler) *Service {
	return &Service{
		deckRepository: dr,
		generateRandom: randGenerator,
		shuffleCard:    cardShuffler,
	}
}

// CreateDeck create deck
func (s *Service) CreateDeck(ctx context.Context, shuffled bool, cardCodes []string) (*entity.Deck, error) {
	cards := CardArray

	if len(cardCodes) > 0 {
		cards = make([]*entity.Card, len(cardCodes))
		for i, code := range cardCodes {
			card, ok := CardMapping[code]
			if !ok {
				return nil, entity.NewError(entity.ErrCardCodeInvalid, entity.ErrMsgCardCodeInvalid)
			}
			cards[i] = &card
		}
	}

	if shuffled {
		cards = s.shuffleCard(s.generateRandom(), cards)
	}

	deck := entity.NewDeck(shuffled, (*entity.Cards)(&cards))
	return s.deckRepository.Insert(ctx, deck)
}

// GetDeck get deck by ID
// will return error when:
//
//	deck not found
func (s *Service) GetDeck(ctx context.Context, id string) (*entity.Deck, error) {
	if id == "" {
		err := entity.NewError(entity.ErrParamInvalid, entity.ErrMsgParamInvalid)
		err.AddDetail(entity.NewErrorDetail("id", "ID is empty"))
		return nil, err
	}

	return s.deckRepository.GetByID(ctx, id)
}

// DrawCards draw cards according to n parameter
// Will return error when:
//
//	deck not found
//	n is larger than remaining card in deck
func (s *Service) DrawCards(ctx context.Context, id string, n int64) (*entity.Cards, error) {
	if id == "" {
		err := entity.NewError(entity.ErrParamInvalid, entity.ErrMsgParamInvalid)
		err.AddDetail(entity.NewErrorDetail("id", "ID is empty"))
		return nil, err
	}

	if n <= 0 {
		err := entity.NewError(entity.ErrParamInvalid, entity.ErrMsgParamInvalid)
		err.AddDetail(entity.NewErrorDetail("count", "count must be bigger than 0"))
		return nil, err
	}

	return s.deckRepository.DrawCards(ctx, id, n)
}
