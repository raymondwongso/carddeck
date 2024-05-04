package entity_test

import (
	"testing"

	"github.com/raymondwongso/carddeck/modules/carddeck/entity"
	"github.com/stretchr/testify/assert"
)

func Test_NewErrorResponse(t *testing.T) {
	err := entity.NewError("internal.error", "something goes wrong")
	assert.Equal(t, "internal.error", err.Code)
	assert.Equal(t, "something goes wrong", err.Message)
}

func Test_NewErrorDetail(t *testing.T) {
	errDetail := entity.NewErrorDetail("shuffled", "shuffled field is invalid")
	assert.Equal(t, "shuffled", errDetail.Field)
	assert.Equal(t, "shuffled field is invalid", errDetail.Message)
}

func Test_WithDetail(t *testing.T) {
	err := entity.NewError("internal.error", "Something goes wrong")
	err.AddDetail(entity.NewErrorDetail("shuffled", "shuffled is invalid"))

	assert.Equal(t, 1, len(err.Details))
}
