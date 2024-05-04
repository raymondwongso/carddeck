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
		{"ACE", "SPADE", "AS"},
		{"2", "SPADE", "2S"},
		{"3", "SPADE", "3S"},
		{"4", "SPADE", "4S"},
		{"5", "SPADE", "5S"},
		{"6", "SPADE", "6S"},
		{"7", "SPADE", "7S"},
		{"8", "SPADE", "8S"},
		{"9", "SPADE", "9S"},
		{"10", "SPADE", "10S"},
		{"JACK", "SPADE", "JS"},
		{"QUEEN", "SPADE", "QS"},
		{"KING", "SPADE", "KS"},

		{"ACE", "DIAMOND", "AD"},
		{"2", "DIAMOND", "2D"},
		{"3", "DIAMOND", "3D"},
		{"4", "DIAMOND", "4D"},
		{"5", "DIAMOND", "5D"},
		{"6", "DIAMOND", "6D"},
		{"7", "DIAMOND", "7D"},
		{"8", "DIAMOND", "8D"},
		{"9", "DIAMOND", "9D"},
		{"10", "DIAMOND", "10D"},
		{"JACK", "DIAMOND", "JD"},
		{"QUEEN", "DIAMOND", "QD"},
		{"KING", "DIAMOND", "KD"},

		{"ACE", "CLUB", "AC"},
		{"2", "CLUB", "2C"},
		{"3", "CLUB", "3C"},
		{"4", "CLUB", "4C"},
		{"5", "CLUB", "5C"},
		{"6", "CLUB", "6C"},
		{"7", "CLUB", "7C"},
		{"8", "CLUB", "8C"},
		{"9", "CLUB", "9C"},
		{"10", "CLUB", "10C"},
		{"JACK", "CLUB", "JC"},
		{"QUEEN", "CLUB", "QC"},
		{"KING", "CLUB", "KC"},

		{"ACE", "HEART", "AH"},
		{"2", "HEART", "2H"},
		{"3", "HEART", "3H"},
		{"4", "HEART", "4H"},
		{"5", "HEART", "5H"},
		{"6", "HEART", "6H"},
		{"7", "HEART", "7H"},
		{"8", "HEART", "8H"},
		{"9", "HEART", "9H"},
		{"10", "HEART", "10H"},
		{"JACK", "HEART", "JH"},
		{"QUEEN", "HEART", "QH"},
		{"KING", "HEART", "KH"},
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
	panic("not implemented")
}

// GetDeck get deck by ID
// will return error when:
//
//	deck not found
func (s *Service) GetDeck(ctx context.Context, id string) (*entity.Deck, error) {
	panic("not implemented")
}

// DrawCards draw cards according to n parameter
// Will return error when:
//
//	deck not found
//	n is larger than remaining card in deck
func (s *Service) DrawCards(ctx context.Context, id string, n int64) (*entity.Cards, error) {
	panic("not implemented")
}
