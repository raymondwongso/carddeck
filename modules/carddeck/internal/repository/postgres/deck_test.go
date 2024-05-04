package postgres_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
	"github.com/raymondwongso/carddeck/modules/carddeck/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	timeTemp    = time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC)
	defaultDeck = entity.NewDeck(false, &entity.Cards{
		{Val: "ACE", Suit: "SPADE", Code: "AS"},
		{Val: "2", Suit: "SPADE", Code: "2S"},
	})
	afterInsertDeck = func() *entity.Deck {
		temp := entity.NewDeck(false, &entity.Cards{
			{Val: "ACE", Suit: "SPADE", Code: "AS"},
			{Val: "2", Suit: "SPADE", Code: "2S"},
		})
		temp.ID = "temp-uuid-abc-def"
		temp.CreatedAt = timeTemp
		temp.UpdatedAt = timeTemp

		return temp
	}()
)

type DeckTestSuite struct {
	suite.Suite
	dbmock sqlmock.Sqlmock
	dbx    *sqlx.DB
}

func (s *DeckTestSuite) SetupSuite() {
	db, dbmock, err := sqlmock.New()
	assert.NoError(s.T(), err)
	s.dbmock = dbmock
	s.dbx = sqlx.NewDb(db, "sqlmock")
}

func TestDeckTestSuite(t *testing.T) {
	suite.Run(t, new(DeckTestSuite))
}

func (s *DeckTestSuite) TestInsert() {
	repo := postgres.NewDeck(s.dbx)
	returningCols := []string{"id", "cards", "shuffled", "created_at", "updated_at"}
	returningVals := []driver.Value{"temp-uuid-abc-def", []byte(`[{"value": "ACE", "suit": "SPADE", "code": "AS"},{"value": "2", "suit": "SPADE", "code": "2S"}]`), false, timeTemp, timeTemp}
	query := `INSERT INTO public.decks (cards, shuffled) VALUES ($1, $2) RETURNING id, cards, shuffled, created_at, updated_at`

	s.Run("success", func() {
		rows := sqlmock.NewRows(returningCols).AddRow(returningVals...)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		deck, err := repo.Insert(context.Background(), defaultDeck)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), afterInsertDeck.ID, deck.ID)
		assert.Equal(s.T(), afterInsertDeck.Cards, deck.Cards)
		assert.Equal(s.T(), afterInsertDeck.Remaining(), deck.Remaining())
		assert.Equal(s.T(), afterInsertDeck.CreatedAt, deck.CreatedAt)
		assert.Equal(s.T(), afterInsertDeck.UpdatedAt, deck.UpdatedAt)
	})

	s.Run("failed - unknown error from repository", func() {
		s.dbmock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("some error"))

		deck, err := repo.Insert(context.Background(), defaultDeck)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), deck)
	})
}
