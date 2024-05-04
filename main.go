/*
Carddeck is an application that serves API for card games.
It is configured by .env file in root directy by default.
See .env.sample for more information

Usage:

	carddeck [command]

Available Commands:

	completion  Generate the autocompletion script for the specified shell
	help        Help about any command
	server      Spin up HTTP Server

Flags:

	-h, --help   help for carddeck

Use "carddeck [command] --help" for more information about a command.
*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raymondwongso/carddeck/config"
	_ "github.com/raymondwongso/carddeck/docs"
	"github.com/raymondwongso/carddeck/modules/carddeck"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	root := &cobra.Command{
		Use:   "carddeck",
		Short: "Carddeck root command",
	}

	root.AddCommand(func() *cobra.Command {
		serverCmd := &cobra.Command{
			Use:   "server",
			Short: "Spin up HTTP Server",
			RunE: func(cmd *cobra.Command, args []string) error {
				return server()
			},
		}

		return serverCmd
	}())

	if err := root.Execute(); err != nil {
		log.Fatal().Err(err).Msg("error executing root command")
	}
}

func server() error {
	config, err := config.Load(".env")
	if err != nil {
		return err
	}

	handler := carddeck.BuildHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /decks", handler.CreateDeck)
	mux.HandleFunc("GET /decks/{id}", handler.GetDeck)
	mux.HandleFunc("GET /decks/{id}/cards", handler.DrawCards)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: mux,
	}

	intrCh := make(chan os.Signal, 1)
	signal.Notify(intrCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("error running server")
		}
	}()
	log.Info().Msg("server started")

	<-intrCh

	log.Info().Msg("stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("error stopping server")
		return err
	}
	log.Info().Msg("server stopped gracefully")

	return nil
}
