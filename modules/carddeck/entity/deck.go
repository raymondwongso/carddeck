// Package entity contains business domain entity widely used by usecases
package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// remainingFunc is a custom type for compute-function that return
type remainingFunc func() int

// MarshalJSON marshal the return value of remainingFunc to json
func (f remainingFunc) MarshalJSON() ([]byte, error) {
	return json.Marshal(f())
}

// UnmarshalJSON unmarshal the value back when remainingFunc is called
func (f remainingFunc) UnmarshalJSON(b []byte) error {
	var i int
	err := json.Unmarshal(b, &i)

	// f = func() int { return i }
	return err
}

// Deck defines a deck of card
type Deck struct {
	ID        string        `json:"id"`
	Shuffled  bool          `json:"shuffled"`
	Remaining remainingFunc `json:"remaining"`
	Cards     *Cards        `json:"cards,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// Cards defines array of card
type Cards []*Card

// Scan implements scanner interface
func (c *Cards) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		_ = json.Unmarshal(v, &c)
		return nil
	case string:
		_ = json.Unmarshal([]byte(v), &c)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

// Value implements valuer interface
func (c *Cards) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Card defines a single card
type Card struct {
	Val  string `json:"value"`
	Suit string `json:"suit"`
	Code string `json:"code"`
}

// Len returns the number of card available
func (c Cards) Len() int { return len(c) }

// Draw draws n number of card, while also returning the remaining cards.
// Return error if n is larger than available cards.
func (c Cards) Draw(n int) (drawedCards Cards, remainingCards Cards, err error) {
	if n > c.Len() {
		return nil, nil, errors.New("count is bigger than card len")
	}

	return c[:n], c[n:], nil
}

// NewDeck returns new deck, also inject the RemainingFunc for Remaining field
func NewDeck(shuffled bool, cards *Cards) *Deck {
	d := Deck{
		Shuffled: shuffled,
		Cards:    cards,
	}
	d.Remaining = RemainingFunc(&d)

	return &d
}

// JSONUnmarshalDeck unmarshal json data to deck and inject RemainingFunc
func JSONUnmarshalDeck(data []byte, dst *Deck) error {
	if dst == nil {
		return errors.New("nil deck destination")
	}

	var deck Deck
	err := json.Unmarshal(data, &deck)
	if err != nil {
		return err
	}

	deck.Remaining = RemainingFunc(&deck)
	*dst = deck
	return nil
}

// RemainingFunc return a function that calculate the len of cards in deck.
func RemainingFunc(d *Deck) remainingFunc {
	return func() int {
		return d.Cards.Len()
	}
}
