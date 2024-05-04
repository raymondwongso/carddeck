package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
	"github.com/rs/zerolog/log"
)

// Deck defines deck repository
type Deck struct {
	db *sqlx.DB
}

// NewDeck returns new deck repository
func NewDeck(db *sqlx.DB) *Deck {
	return &Deck{db: db}
}

// Insert insert new deck to database
func (d *Deck) Insert(ctx context.Context, deck *entity.Deck) (*entity.Deck, error) {
	query := `INSERT INTO public.decks (cards, shuffled) VALUES ($1, $2) RETURNING id, cards, shuffled, created_at, updated_at`

	row := d.db.QueryRowxContext(ctx, query, deck.Cards, deck.Shuffled)
	if err := row.Scan(&deck.ID, &deck.Cards, &deck.Shuffled, &deck.CreatedAt, &deck.UpdatedAt); err != nil {
		return nil, err
	}

	return deck, nil
}

// GetByID get deck by ID
func (d *Deck) GetByID(ctx context.Context, id string) (*entity.Deck, error) {
	query := `SELECT id, cards, shuffled, created_at, updated_at FROM public.decks WHERE id = $1`

	deck := entity.NewDeck(false, nil)
	row := d.db.QueryRowxContext(ctx, query, id)
	if err := row.Scan(&deck.ID, &deck.Cards, &deck.Shuffled, &deck.CreatedAt, &deck.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.NewError(entity.ErrDeckNotFound, entity.ErrMsgDeckNotFound)
		}
		return nil, err
	}

	return deck, nil
}

func (d *Deck) DrawCards(ctx context.Context, id string, count int64) (cards *entity.Cards, err error) {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("error rollbacking draw cards")
			}
		}
	}()

	selectForUpdateQuery := `SELECT id, cards, shuffled, created_at, updated_at FROM public.decks WHERE id = $1 FOR UPDATE`
	deck := entity.NewDeck(false, nil)
	row := tx.QueryRowContext(ctx, selectForUpdateQuery, id)
	if err := row.Scan(&deck.ID, &deck.Cards, &deck.Shuffled, &deck.CreatedAt, &deck.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.NewError(entity.ErrDeckNotFound, entity.ErrMsgDeckNotFound)
		}
		return nil, err
	}

	drawwed, remaining, err := deck.Cards.Draw(count)
	if err != nil {
		return nil, err
	}

	updateQuery := `UPDATE public.decks SET cards=$2, updated_at=NOW() WHERE id = $1 RETURNING id, cards, shuffled, created_at, updated_at`
	row = tx.QueryRowContext(ctx, updateQuery, id, &remaining)
	if err := row.Scan(&deck.ID, &deck.Cards, &deck.Shuffled, &deck.CreatedAt, &deck.UpdatedAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &drawwed, nil
}
