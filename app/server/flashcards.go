package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var server *Flashcards

// Flashcards Structure
type Flashcards struct {
	Database  *mongo.Database
	Router    *gin.Engine
	Version   string
	Port      string
	TokenKey  string
	Origin    string
	LogFormat string
	Mode      string
	DBHost    string
}

func (fc *Flashcards) ParseParameters() {
	fc.LogFormat = os.Getenv("LOG_FORMAT")
	fc.Version = os.Getenv("API_VERSION")
	fc.Port = os.Getenv("API_PORT")
	fc.TokenKey = os.Getenv("TOKEN_KEY")
	fc.Origin = os.Getenv("ALLOW_ORIGIN")
	fc.Mode = os.Getenv("MODE")
	fc.DBHost = os.Getenv("DB_HOST")
}

// ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections.
// https://github.com/gin-gonic/gin
func (fc *Flashcards) ListenAndServe() error {
	srv := &http.Server{
		Addr:              fc.Port,
		Handler:           fc.Router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// start
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Msgf("Unable to listen and serve: %v", err)
		return err
	}
	return nil
}

// SetServer init mongo database
func SetServer(s *Flashcards) {
	server = s
}

// GetServer Flashcards
func GetServer() *Flashcards {
	return server
}
