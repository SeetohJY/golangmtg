package uploadFile

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	// "os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

const (
	projectID  = "jy-project-6cfb8"  // FILL IN WITH YOURS
	bucketName = "auto-expiring-mtgjsondata-bucket" // FILL IN WITH YOURS
)

func init() {
        // functions.HTTP("DownloadmtgcsvHTTP", DownloadmtgcsvHTTP)
}

func main() {
        // os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "jy-project-credentials.json")

	r := gin.Default()
	r.POST("/upload", upload)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}


type uploadInput struct {
        FileName string `json:"fileName"`
        Url string `json:"url"`
}

func upload(c *gin.Context) {
        var newUploadInput uploadInput

        if err := c.BindJSON(&newUploadInput); err != nil {
                c.JSON(400, gin.H{
                        "message": "Invalid Input",
                })
        }

        uploadFile(newUploadInput.FileName, newUploadInput.Url)

        c.JSON(200, gin.H{
                "message": "success",
        })
}

// UploadFile uploads an object
func uploadFile(fileName string, url string) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
                log.Fatalf("Bad status: %s", resp.Status)
	}

	// Set up Google Cloud Storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create a new bucket handle
	bucket := client.Bucket(bucketName)

	// Create a new object in the bucket
	obj := bucket.Object(fileName)
	writer := obj.NewWriter(ctx)

	// Copy the contents of the CSV file to the object
	if _, err := io.Copy(writer, resp.Body); err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	// Close the object writer
	if err := writer.Close(); err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}

	fmt.Println("File uploaded successfully!")
}


func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}