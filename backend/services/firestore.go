package services

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"read-robin/models"
	"read-robin/utils"

	"cloud.google.com/go/firestore"
)

// FirestoreClient is a wrapper around the Firestore client
type FirestoreClient struct {
	Client *firestore.Client
}

// NewFirestoreClient creates a new Firestore client
func NewFirestoreClient(ctx context.Context) (*FirestoreClient, error) {
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT environment variable not set")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %v", err)
	}

	return &FirestoreClient{Client: client}, nil
}

// SaveQuizWithTitle saves a quiz and its title to Firestore
func (fc *FirestoreClient) SaveQuiz(ctx context.Context, url, title string, quiz models.Quiz) error {
	collection := "quizzes"

	contentID := utils.GenerateID(url)
	docRef := fc.Client.Collection(collection).Doc(contentID)

	// Get the existing document or create a new one
	doc, err := docRef.Get(ctx)
	var content models.Content
	if err == nil {
		if err := doc.DataTo(&content); err != nil {
			return fmt.Errorf("failed to parse existing content: %v", err)
		}
	} else {
		content = models.Content{
			URL:       url,
			Timestamp: time.Now(),
			ContentID: contentID,
			Title:     title,
			Quizzes:   []models.Quiz{},
		}
	}

	// Add logging to check existing quizzes
	fmt.Printf("Existing quizzes: %v\n", content.Quizzes)

	latestQuizID := GetLatestQuizID(content.Quizzes)
	quiz.QuizID = latestQuizID
	quiz.Timestamp = time.Now()

	content.Quizzes = append(content.Quizzes, quiz)
	content.Title = title // Update the title in case it was not previously set

	_, err = docRef.Set(ctx, content)
	if err != nil {
		return fmt.Errorf("failed adding quiz: %v", err)
	}
	return nil
}

// GetQuiz retrieves a quiz from Firestore by contentID and quizID
func (fc *FirestoreClient) GetQuiz(ctx context.Context, contentID, quizID string) (*models.Quiz, error) {
	collection := "quizzes"

	doc, err := fc.Client.Collection(collection).Doc(contentID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving quiz: %v", err)
	}

	var content models.Content
	if err := doc.DataTo(&content); err != nil {
		return nil, fmt.Errorf("dataTo: %v", err)
	}

	for _, quiz := range content.Quizzes {
		if quiz.QuizID == quizID {
			return &quiz, nil
		}
	}

	return nil, fmt.Errorf("no quiz found for quizID: %s", quizID)
}

// GetExistingQuizzes fetches existing quizzes from Firestore
func (fc *FirestoreClient) GetExistingQuizzes(ctx context.Context, contentID string) ([]models.Quiz, error) {
	doc, err := fc.Client.Collection("quizzes").Doc(contentID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var existingContent models.Content
	if err := doc.DataTo(&existingContent); err != nil {
		return nil, err
	}
	return existingContent.Quizzes, nil
}

// GetLatestQuizID generates the next sequential quiz ID for the given quizzes
func GetLatestQuizID(quizzes []models.Quiz) string {
	if len(quizzes) == 0 {
		return "0001"
	}
	maxID := 0
	for _, quiz := range quizzes {
		id, err := strconv.Atoi(quiz.QuizID)
		if err == nil && id > maxID {
			maxID = id
		}
	}
	nextID := fmt.Sprintf("%04d", maxID+1)
	fmt.Printf("Next quiz ID: %s\n", nextID) // Add logging to check next quiz ID
	return nextID
}
