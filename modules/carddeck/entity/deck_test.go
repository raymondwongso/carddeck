package entity_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
	"github.com/stretchr/testify/assert"
)

func Test_NewDeck(t *testing.T) {
	deck := entity.NewDeck(false, &entity.Cards{
		{Val: "ACE", Suit: "SPADE", Code: "AS"},
		{Val: "2", Suit: "SPADE", Code: "2S"},
		{Val: "3", Suit: "SPADE", Code: "3S"},
	})

	assert.Equal(t, false, deck.Shuffled)
	assert.Equal(t, &entity.Cards{
		{Val: "ACE", Suit: "SPADE", Code: "AS"},
		{Val: "2", Suit: "SPADE", Code: "2S"},
		{Val: "3", Suit: "SPADE", Code: "3S"},
	}, deck.Cards)
	assert.Equal(t, 3, deck.Remaining())
}

func Test_Cards_Len(t *testing.T) {
	cards := entity.Cards{
		{Val: "ACE", Suit: "SPADE", Code: "AS"},
		{Val: "2", Suit: "SPADE", Code: "2S"},
		{Val: "3", Suit: "SPADE", Code: "3S"},
	}

	assert.Equal(t, 3, cards.Len())
}

func Test_Cards_Draw(t *testing.T) {
	t.Run("success draw", func(t *testing.T) {
		cards := entity.Cards{
			{Val: "ACE", Suit: "SPADE", Code: "AS"},
			{Val: "2", Suit: "SPADE", Code: "2S"},
			{Val: "3", Suit: "SPADE", Code: "3S"},
		}

		drawed, remaining, err := cards.Draw(2)
		assert.Equal(t, entity.Cards{
			{Val: "ACE", Suit: "SPADE", Code: "AS"},
			{Val: "2", Suit: "SPADE", Code: "2S"},
		}, drawed)
		assert.Equal(t, entity.Cards{
			{Val: "3", Suit: "SPADE", Code: "3S"},
		}, remaining)
		assert.NoError(t, err)
	})

	t.Run("failed draw count is larger than remaining", func(t *testing.T) {
		cards := entity.Cards{
			{Val: "ACE", Suit: "SPADE", Code: "AS"},
			{Val: "2", Suit: "SPADE", Code: "2S"},
			{Val: "3", Suit: "SPADE", Code: "3S"},
		}

		_, _, err := cards.Draw(4)
		assert.Error(t, err)
	})
}

func Test_JSON(t *testing.T) {
	t.Run("success marshal", func(t *testing.T) {
		deck := entity.NewDeck(false, &entity.Cards{
			{Val: "ACE", Suit: "SPADE", Code: "AS"},
			{Val: "2", Suit: "SPADE", Code: "2S"},
			{Val: "3", Suit: "SPADE", Code: "3S"},
		})
		timeTemp := time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC)
		deck.ID = "some-uuid-abc-def"
		deck.CreatedAt = timeTemp
		deck.UpdatedAt = timeTemp

		marshaled, err := json.Marshal(&deck)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"id": "some-uuid-abc-def",
			"shuffled": false,
			"remaining": 3,
			"cards": [
				{
					"value": "ACE",
					"suit": "SPADE",
					"code": "AS"
				},
				{
					"value": "2",
					"suit": "SPADE",
					"code": "2S"
				},
				{
					"value": "3",
					"suit": "SPADE",
					"code": "3S"
				}
			],
			"created_at": "2022-01-01T01:00:00Z",
			"updated_at": "2022-01-01T01:00:00Z"
		}`, string(marshaled))
	})

	t.Run("success unmarshal", func(t *testing.T) {
		jsonString := `{
			"id": "some-uuid-abc-def",
			"shuffled": false,
			"remaining": 3,
			"cards": [
				{
					"value": "ACE",
					"suit": "SPADE",
					"code": "AS"
				},
				{
					"value": "2",
					"suit": "SPADE",
					"code": "2S"
				},
				{
					"value": "3",
					"suit": "SPADE",
					"code": "3S"
				}
			],
			"created_at": "2022-01-01T01:00:00Z",
			"updated_at": "2022-01-01T01:00:00Z"
		}`

		jsonData := []byte(jsonString)

		var deck entity.Deck
		err := entity.JSONUnmarshalDeck(jsonData, &deck)
		assert.NoError(t, err)

		expected := entity.NewDeck(false, &entity.Cards{
			{Val: "ACE", Suit: "SPADE", Code: "AS"},
			{Val: "2", Suit: "SPADE", Code: "2S"},
			{Val: "3", Suit: "SPADE", Code: "3S"},
		})
		timeTemp := time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC)
		expected.ID = "some-uuid-abc-def"
		expected.CreatedAt = timeTemp
		expected.UpdatedAt = timeTemp

		assert.Equal(t, expected.ID, deck.ID)
		assert.Equal(t, expected.CreatedAt, deck.CreatedAt)
		assert.Equal(t, expected.UpdatedAt, deck.UpdatedAt)
		assert.Equal(t, expected.Remaining(), deck.Remaining())
		assert.Equal(t, expected.Cards, deck.Cards)
	})

	t.Run("success unmarshal", func(t *testing.T) {
		jsonString := `{
			"id": "some-uuid-abc-def",
			"shuffled": false,
			"remaining": 3,
			"cards": [
				{
					"value": "ACE",
					"suit": "SPADE",
					"code": "AS"
				},
				{
					"value": "2",
					"suit": "SPADE",
					"code": "2S"
				},
				{
					"value": "3",
					"suit": "SPADE",
					"code": "3S"
				}
			],
			"created_at": "2022-01-01T01:00:00Z",
			"updated_at": "2022-01-01T01:00:00Z"
		}`

		jsonData := []byte(jsonString)

		err := entity.JSONUnmarshalDeck(jsonData, nil)
		assert.Error(t, err)
	})
}
