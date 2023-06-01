package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	port       = 80
	listenport = ":" + strconv.Itoa(port)
	uploadDir  = "./html/upload"
)

func indexHandler() {
	// Serve static files from the "./html" directory
	fileServer := http.FileServer(http.Dir("./html"))
	http.Handle("/", fileServer)
}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	// Handle file uploads via HTTP POST method

	if request.Method != "POST" {
		// If the method is not POST, return an error
		http.Error(writer, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		// If there is an error reading the uploaded file, return an error
		http.Error(writer, "Failed to read uploaded file.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file in the upload directory with the same name as the uploaded file
	filePath := filepath.Join(uploadDir, fileHeader.Filename)
	newFile, err := os.Create(filePath)
	if err != nil {
		// If there is an error creating the file, return an error
		http.Error(writer, "Failed to create file.", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	// Copy the uploaded file data to the new file
	_, err = io.Copy(newFile, file)
	if err != nil {
		// If there is an error saving the file, return an error
		http.Error(writer, "Failed to save file.", http.StatusInternalServerError)
		return
	}

	//fmt.Fprintf(writer, "File uploaded successfully.")
	http.Redirect(writer, request, "/succes.html", http.StatusSeeOther)

}
func directoryChecker() {
	// Check if upload directory exists, create it if not
	_, err := os.Stat(uploadDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(uploadDir, 0755)
		if errDir != nil {
			log.Fatal("Failed to create upload directory:", errDir)
		}
	}
}

func clearDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		log.Printf("Failed to clear directory %s: %s\n", dirPath, err)
		return err
	}
	log.Printf("Directory %s cleared successfully.\n", dirPath)
	return nil
}

func main() {

	clearDirectory(uploadDir)
	err := clearDirectory(uploadDir)

	if err != nil {
		log.Fatal(err)
	}

	// Function for checking if the upload directory exists
	directoryChecker()

	// Register a handler for serving static files
	indexHandler()

	// Register a handler for file uploads
	http.HandleFunc("/upload", uploadHandler)

	fmt.Print("Starting server at port ", port, "\n")
	if err := http.ListenAndServe(listenport, nil); err != nil {
		log.Fatal(err)
	}
}
