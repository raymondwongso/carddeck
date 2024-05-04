// Package carddeck contains implementation for usecases related to French card deck.
package carddeck

import (
	"github.com/raymondwongso/carddeck/modules/carddeck/internal/rest"
	"github.com/raymondwongso/carddeck/modules/carddeck/internal/service"
)

// BuildHandler build and returns handler
func BuildHandler() *rest.Handler {
	svc := service.New()
	return rest.NewHandler(svc)
}
