package flashcard_service

import (
	"Flashcards/app/functions"
	"Flashcards/app/models"
	"Flashcards/app/server"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"
)

type FlashcardService struct {
	validate *validator.Validate
}

func New() *FlashcardService {
	return &FlashcardService{
		validate: validator.New(),
	}
}

func (f *FlashcardService) Get() ([]models.Flashcard, error) {
	var (
		err        error
		flashcards []models.Flashcard
		flashcard  models.Flashcard
		cursor     *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(flashcard.Collection())

	cursor, err = collection.Find(context.TODO(), bson.D{})

	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var flashcard models.Flashcard
		err = cursor.Decode(&flashcard)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		flashcards = append(flashcards, flashcard)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return flashcards, nil
}

func (f *FlashcardService) Create(flashcardDTO models.FlashcardDTO) (*models.Flashcard, error) {
	var flashcard models.Flashcard

	srv := server.GetServer()
	collection := srv.Database.Collection(flashcard.Collection())

	err := f.validate.Struct(flashcardDTO)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	err = functions.ConvertInputStructToDataStruct(flashcardDTO, &flashcard)

	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	flashcard.ID = functions.NewUUID()
	flashcard.CreatedAt = time.Now()

	_, err = collection.InsertOne(context.TODO(), flashcard)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &flashcard, nil
}

func (f *FlashcardService) GetById(id string) (models.Flashcard, error) {
	var flashcard models.Flashcard

	srv := server.GetServer()
	collection := srv.Database.Collection(flashcard.Collection())

	err := collection.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&flashcard)

	if err != nil {
		log.Error().Err(err).Msg("")
		return models.Flashcard{}, err
	}

	return flashcard, nil
}

func (f *FlashcardService) GetByTag(tag string) ([]models.Flashcard, error) {
	var (
		err        error
		flashcards []models.Flashcard
		flashcard  models.Flashcard
		cursor     *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(flashcard.Collection())

	cursor, err = collection.Find(context.TODO(), bson.D{{"tags", tag}})

	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var flashcard models.Flashcard
		err = cursor.Decode(&flashcard)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		flashcards = append(flashcards, flashcard)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return flashcards, nil
}

func (f *FlashcardService) GetByIds(ids []string) (flashcards []models.Flashcard, err error) {
	for _, id := range ids {
		flashcard, err := f.GetById(id)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		flashcards = append(flashcards, flashcard)
	}

	return flashcards, err
}

func (f *FlashcardService) GetRandomsByTag(tag string, number int) ([]models.Flashcard, error) {
	flashcards, err := f.GetByTag(tag)
	if err != nil {
		return nil, err
	}

	if number >= len(flashcards) {
		return flashcards, nil
	}

	rand.Shuffle(len(flashcards), func(i, j int) {
		flashcards[i], flashcards[j] = flashcards[j], flashcards[i]
	})

	return flashcards[:number], nil
}
