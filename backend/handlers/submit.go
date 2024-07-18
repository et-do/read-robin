package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/models"
	"read-robin/services"
	"read-robin/utils"
	"strings"

	"golang.org/x/net/context"
)

// URLRequest is a struct to hold the URL and persona details submitted by the user
type URLRequest struct {
	URL         string         `json:"url"`
	Persona     models.Persona `json:"persona"`
	ContentType string         `json:"content_type"`
}

// SubmitResponse is a struct to hold the response to be sent back to the user
type SubmitResponse struct {
	Status      string `json:"status"`
	URL         string `json:"url"`
	ContentID   string `json:"content_id"`
	QuizID      string `json:"quiz_id"`
	Title       string `json:"title"` // Add this line
	IsFirstQuiz bool   `json:"is_first_quiz"`
}

// decodeURLRequest decodes the URL request from the HTTP request
func decodeURLRequest(r *http.Request) (URLRequest, error) {
	var urlRequest URLRequest
	if r.Header.Get("Content-Type") == "application/json" {
		err := utils.DecodeJSONBody(r, &urlRequest)
		return urlRequest, err
	} else {
		err := utils.DecodeFormBody(r, "url", &urlRequest.URL)
		return urlRequest, err
	}
}

// normalizeAndGenerateID normalizes the URL and generates a content ID
func normalizeAndGenerateID(url string) (string, string, error) {
	normalizedURL, err := utils.NormalizeURL(url)
	if err != nil {
		return "", "", err
	}
	contentID := services.GenerateID(normalizedURL)
	return normalizedURL, contentID, nil
}

// createFirestoreClient creates a new Firestore client
func createFirestoreClient(ctx context.Context) (*services.FirestoreClient, error) {
	return services.NewFirestoreClient(ctx)
}

// createGeminiClient creates a new Gemini client
func createGeminiClient(ctx context.Context) (*services.GeminiClient, error) {
	return services.NewGeminiClient(ctx)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	urlRequest, err := decodeURLRequest(r)
	if err != nil {
		log.Printf("SubmitHandler: Unable to parse request: %v", err)
		http.Error(w, "Unable to parse request", http.StatusBadRequest)
		return
	}

	log.Printf("SubmitHandler: Received Request: %s", urlRequest)

	htmlContent, err := utils.FetchHTML(urlRequest.URL)
	if err != nil {
		log.Printf("SubmitHandler: Error fetching HTML content: %v", err)
		http.Error(w, "Error fetching HTML content", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	geminiClient, err := createGeminiClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Gemini client: %v", err)
		http.Error(w, "Error creating Gemini client", http.StatusInternalServerError)
		return
	}

	firestoreClient, err := createFirestoreClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Firestore client: %v", err)
		http.Error(w, "Error creating Firestore client", http.StatusInternalServerError)
		return
	}

	contentType := r.Header.Get("Content-Type")
	var response models.SubmitResponse

	if contentType == "application/json" {
		urlRequest, err := decodeURLRequest(r)
		if err != nil {
			log.Printf("SubmitHandler: Unable to parse request: %v", err)
			http.Error(w, "Unable to parse request", http.StatusBadRequest)
			return
		}

		response, err = processURLSubmission(ctx, urlRequest, geminiClient, firestoreClient)
		if err != nil {
			log.Printf("SubmitHandler: Error processing URL submission: %v", err)
			http.Error(w, "Error processing URL submission", http.StatusInternalServerError)
			return
		}

	} else if strings.Contains(contentType, "multipart/form-data") {
		response, err = handleMultipartForm(r, ctx, geminiClient, firestoreClient)
		if err != nil {
			log.Printf("SubmitHandler: Error processing PDF submission: %v", err)
			http.Error(w, "Error processing PDF submission", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	log.Printf("SubmitHandler: Response - %v\n", response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("SubmitHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	log.Println("SubmitHandler: Response sent successfully")
}
