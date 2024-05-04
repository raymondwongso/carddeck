package postgres_test

import (
	"context"
	"database/sql"
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
	afterDrawCards = entity.Cards{
		{Val: "ACE", Suit: "SPADE", Code: "AS"},
	}
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

func (s *DeckTestSuite) TestGetByID() {
	repo := postgres.NewDeck(s.dbx)
	returningCols := []string{"id", "cards", "shuffled", "created_at", "updated_at"}
	returningVals := []driver.Value{"temp-uuid-abc-def", []byte(`[{"value": "ACE", "suit": "SPADE", "code": "AS"},{"value": "2", "suit": "SPADE", "code": "2S"}]`), false, timeTemp, timeTemp}
	query := `SELECT id, cards, shuffled, created_at, updated_at FROM public.decks WHERE id = $1`

	s.Run("success", func() {
		rows := sqlmock.NewRows(returningCols).AddRow(returningVals...)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		deck, err := repo.GetByID(context.Background(), "temp-uuid-abc-def")
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), afterInsertDeck.ID, deck.ID)
		assert.Equal(s.T(), afterInsertDeck.Cards, deck.Cards)
		assert.Equal(s.T(), afterInsertDeck.Remaining(), deck.Remaining())
		assert.Equal(s.T(), afterInsertDeck.CreatedAt, deck.CreatedAt)
		assert.Equal(s.T(), afterInsertDeck.UpdatedAt, deck.UpdatedAt)
	})

	s.Run("failed - unknown error from repository", func() {
		s.dbmock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("some error"))

		deck, err := repo.GetByID(context.Background(), "abc")
		assert.Error(s.T(), err)
		assert.Nil(s.T(), deck)
	})

	s.Run("failed - no rows result", func() {
		s.dbmock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(sql.ErrNoRows)

		deck, err := repo.GetByID(context.Background(), "abc")
		assert.Error(s.T(), err)
		assert.Nil(s.T(), deck)

		perr, ok := err.(*entity.Error)
		assert.True(s.T(), ok)
		assert.Equal(s.T(), entity.ErrDeckNotFound, perr.Code)
		assert.Equal(s.T(), entity.ErrMsgDeckNotFound, perr.Message)
	})
}

func (s *DeckTestSuite) TestDrawCards() {
	repo := postgres.NewDeck(s.dbx)
	returningCols := []string{"id", "cards", "shuffled", "created_at", "updated_at"}
	selectVals := []driver.Value{"temp-uuid-abc-def", []byte(`[{"value": "ACE", "suit": "SPADE", "code": "AS"},{"value": "2", "suit": "SPADE", "code": "2S"}]`), false, timeTemp, timeTemp}
	updateVals := []driver.Value{"temp-uuid-abc-def", []byte(`[{"value": "ACE", "suit": "SPADE", "code": "AS"}]`), false, timeTemp, timeTemp}
	selectForUpdateQuery := `SELECT id, cards, shuffled, created_at, updated_at FROM public.decks WHERE id = $1 FOR UPDATE`
	updateQuery := `UPDATE public.decks SET cards=$2, updated_at=NOW() WHERE id = $1 RETURNING id, cards, shuffled, created_at, updated_at`

	s.Run("success", func() {
		s.dbmock.ExpectBegin()

		selectRows := sqlmock.NewRows(returningCols).AddRow(selectVals...)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(selectForUpdateQuery)).WillReturnRows(selectRows)
		updateRows := sqlmock.NewRows(returningCols).AddRow(updateVals...)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(updateQuery)).WillReturnRows(updateRows)

		s.dbmock.ExpectCommit()

		cards, err := repo.DrawCards(context.Background(), "temp-uuid-abc-def", 1)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), &afterDrawCards, cards)
	})

	s.Run("failed - begin transaction failed", func() {
		s.dbmock.ExpectBegin().WillReturnError(errors.New("some error"))

		cards, err := repo.DrawCards(context.Background(), "temp-uuid-abc-def", 1)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), cards)
	})

	s.Run("failed - select for update failed", func() {
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(regexp.QuoteMeta(selectForUpdateQuery)).WillReturnError(errors.New("some error"))

		s.dbmock.ExpectRollback()

		cards, err := repo.DrawCards(context.Background(), "temp-uuid-abc-def", 1)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), cards)
	})

	s.Run("failed - update failed", func() {
		s.dbmock.ExpectBegin()

		selectRows := sqlmock.NewRows(returningCols).AddRow(selectVals...)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(selectForUpdateQuery)).WillReturnRows(selectRows)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(updateQuery)).WillReturnError(errors.New("some error"))

		s.dbmock.ExpectRollback()

		cards, err := repo.DrawCards(context.Background(), "temp-uuid-abc-def", 1)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), cards)
	})

	s.Run("failed - commit failed", func() {
		s.dbmock.ExpectBegin()

		selectRows := sqlmock.NewRows(returningCols).AddRow(selectVals...)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(selectForUpdateQuery)).WillReturnRows(selectRows)
		updateRows := sqlmock.NewRows(returningCols).AddRow(updateVals...)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(updateQuery)).WillReturnRows(updateRows)

		s.dbmock.ExpectCommit().WillReturnError(errors.New("some error"))

		cards, err := repo.DrawCards(context.Background(), "temp-uuid-abc-def", 1)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), cards)
	})

	s.Run("failed - rollback failed", func() {
		s.dbmock.ExpectBegin()

		selectRows := sqlmock.NewRows(returningCols).AddRow(selectVals...)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(selectForUpdateQuery)).WillReturnRows(selectRows)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(updateQuery)).WillReturnError(errors.New("some error"))

		s.dbmock.ExpectRollback().WillReturnError(errors.New("some error"))

		cards, err := repo.DrawCards(context.Background(), "temp-uuid-abc-def", 1)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), cards)
	})

	s.Run("failed - draw count is larger than available", func() {
		s.dbmock.ExpectBegin()

		selectRows := sqlmock.NewRows(returningCols).AddRow(selectVals...)
		s.dbmock.ExpectQuery(regexp.QuoteMeta(selectForUpdateQuery)).WillReturnRows(selectRows)

		s.dbmock.ExpectRollback()

		cards, err := repo.DrawCards(context.Background(), "temp-uuid-abc-def", 999)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), cards)

		perr, ok := err.(*entity.Error)
		assert.True(s.T(), ok)
		assert.Equal(s.T(), entity.ErrDeckCardInsufficient, perr.Code)
		assert.Equal(s.T(), entity.ErrMsgDeckCardInsufficient, perr.Message)
	})
}
