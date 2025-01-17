package student_service

import (
	"Flashcards/app/functions"
	"Flashcards/app/models"
	"Flashcards/app/mongodb"
	"Flashcards/app/server"
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentService struct {
	validate *validator.Validate
}

func New() *StudentService {
	return &StudentService{
		validate: validator.New(),
	}
}

// Get services to get list of student on db
func (s *StudentService) Get(queryParams models.QueryParams) ([]models.Student, error) {
	var (
		err      error
		students []models.Student
		student  models.Student
		cursor   *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(student.Collection())

	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		// A new result variable should be declared for each document.
		var student models.Student
		err = cursor.Decode(&student)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		students = append(students, student)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return students, err
}

// Create new student on db
func (s *StudentService) Create(in *models.StudentInput) (*models.Student, error) {
	var student models.Student

	srv := server.GetServer()
	collection := srv.Database.Collection(student.Collection())

	// Check input fields
	err := s.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	err = functions.ConvertInputStructToDataStruct(in, &student)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	student.CustomID = functions.NewUUID()
	student.CreatedAt = time.Now()
	student.Suspended = false

	_, err = collection.InsertOne(context.TODO(), student)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &student, nil
}

// GetByID controller to get one student by ID
func (s *StudentService) GetByID(id string) (models.Student, error) {
	var (
		err         error
		student     models.Student
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(student.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	err = collection.FindOne(context.TODO(), filter).Decode(&student)
	if err == nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Error().Err(err).Msg("")
			return student, err
		}

	}
	return student, err
}

// Update controller to update a student
func (s *StudentService) Update(id string, in *models.StudentInput) error {
	var (
		doc         interface{}
		result      *mongo.UpdateResult
		err         error
		queryParams models.QueryParams
		student     models.Student
	)

	srv := server.GetServer()

	// Check input fields
	err = s.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	student, err = s.GetByID(id)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	err = functions.ConvertInputStructToDataStruct(in, &student)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	collection := srv.Database.Collection(student.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	if doc, err = mongodb.ToDoc(student); err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	update := bson.M{"$set": doc}
	result, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		err = errors.New("StudentService to be modified was not found")
	}

	if err == nil && result.ModifiedCount == 0 {
		err = errors.New("StudentService could not be updated")
	}
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	return err
}

// Suspend controller to suspend a student
func (s *StudentService) Suspend(id string) error {
	var (
		err         error
		queryParams models.QueryParams
		student     models.Student
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(student.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "suspended", Value: true},
		}},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)

	return err
}

// GetByIds controller to get list of student by Ids
func (s *StudentService) GetByIds(ids []string) ([]models.Student, error) {
	var students []models.Student
	for _, id := range ids {
		student, err := s.GetByID(id)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}
