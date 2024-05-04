// Package carddeck contains implementation for usecases related to French card deck.
package carddeck

import "github.com/raymondwongso/carddeck/modules/carddeck/internal/rest"

// BuildHandler build and returns handler
func BuildHandler() *rest.Handler {
	return rest.NewHandler()
}
