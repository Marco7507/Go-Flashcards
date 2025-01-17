package main

import (
	"Flashcards/app/mongodb"
	flashcard "Flashcards/app/routes/flashcard"
	session "Flashcards/app/routes/session"
	student "Flashcards/app/routes/student"
	"Flashcards/app/server"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newFlashcardsServer() error {

	// loading .env files in dev mode
	if os.Getenv("MODE") == "" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	srv := &server.Flashcards{}

	srv.ParseParameters()

	// log format definition
	switch srv.LogFormat {
	case "HUMAN":
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	case "JSON":
		// Already default
	default:
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
	}

	// setup router
	srv.Router = setupRouter()

	// Init MongoDB
	client, err := mongodb.OpenMongoDB(srv.DBHost)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to Open database")
		return err
	}
	srv.Database = client.Database("flashcards")

	student.SetupRouter(srv.Router)
	flashcard.SetupRouter(srv.Router)
	session.SetupRouter(srv.Router)
	server.SetServer(srv)

	return nil
}
