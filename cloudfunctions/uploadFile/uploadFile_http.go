package uploadFile

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func init() {
	funcframework.RegisterHTTPFunction("/upload", UploadFiletoStorageBucket)
	funcframework.Start("8080")
}

type uploadInput struct {
	FileName string `json:"fileName"`
	Url      string `json:"url"`
}

type Response struct {
	Message string `json:"result"`
}

func UploadFiletoStorageBucket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req uploadInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Process the request and create a response
	result := uploadFile(req.FileName, req.Url)

	res := &Response{
		Message: result,
	}

	// Encode the response as JSON and write it to the response writer
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error encoding response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UploadFile uploads an object
func uploadFile(fileName string, url string) string {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
		CheckError(err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad status: %s", resp.Status)
		CheckError(err)
	}

	// Set up Google Cloud Storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		CheckError(err)
	}
	defer client.Close()

	// Create a new bucket handle
	bucket := client.Bucket("mtgjson-jy-bucket")

	// Create a new object in the bucket
	obj := bucket.Object(fileName)
	writer := obj.NewWriter(ctx)

	// Copy the contents of the CSV file to the object
	if _, err := io.Copy(writer, resp.Body); err != nil {
		log.Fatalf("Failed to upload file: %v", err)
		CheckError(err)
	}

	// Close the object writer
	if err := writer.Close(); err != nil {
		log.Fatalf("Failed to close writer: %v", err)
		CheckError(err)
	}

	return "File successfully uploaded"
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
