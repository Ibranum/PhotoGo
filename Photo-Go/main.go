package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)


func main() {
	http.HandleFunc("/fileDownload/", fileDownload)
	http.HandleFunc("/fileUpload", fileUpload)
	http.HandleFunc("/listUploads", listUploads)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func fileDownload(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /fileDownload request\n")
	// 
	guid := strings.TrimPrefix(r.URL.Path, "/fileDownload/")
	filePath := filepath.Join("./fileUploads/", guid)
	fmt.Println(guid, filePath)
	// 
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
    	http.NotFound(w, r)
    	return
	}
	// 
	w.Header().Set("Content-Disposition", "attachment; filename="+guid)
  	w.Header().Set("Content-Type", "application/octet-stream")
	//
	http.ServeFile(w, r, filePath)

	//io.WriteString(w, "This is file download!\n")
}
func fileUpload(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("got /hello request\n")
	// Parse multipart form data
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
   	// Handle errors for HTTP
	if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       return
   	}
	// Close the file that was opened above
   	defer file.Close()

	// Generate a GUID using the Google UUID library
	guid := uuid.New().String()
	// 
   	ext := path.Ext(handler.Filename)
	// Generate a new name for the file
   	newFileName := guid + ext
	// Create new GUID-named file
	dst, err := os.Create("./fileUploads/" + newFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Close the file creation handler.
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//io.WriteString(w, "This is file upload!\n")
}
func listUploads(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")

	//
	files, err := os.ReadDir("./fileUploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	if len(files) == 0 {
       fmt.Fprintln(w, "No files found.")
       return
   	}
	//
	for _, file := range files {
   		fmt.Fprintln(w, file.Name())
 	}


	//io.WriteString(w, "This is list uploads!\n")
}