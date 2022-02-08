package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var (
	serveDir   string = "./upload"
	serverAddr string = "localhost:8000"
)

type FileInfo struct {
	Name string
	Extr string
	Size int64
}

func uploadHandleFunc(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Print("can't read file")
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	filePath := serveDir + "/" + header.Filename
	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fileLink := "http://" + serverAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink)
}

func listHandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not GET method", http.StatusBadRequest)
		return
	}

	files, err := ioutil.ReadDir(serveDir)
	if err != nil {
		log.Printf("can't read file from derectory %s", serveDir)
		http.Error(w, "Unable to read directory", http.StatusBadRequest)
		return
	}

	ext := r.FormValue("ext")
	var filesInfo []FileInfo
	for _, someFile := range files {
		if !someFile.IsDir() {
			fileAttr := FileInfo{
				Name: someFile.Name(),
				Extr:  filepath.Ext(someFile.Name()),
				Size: someFile.Size(),
			}
			if ext == "" || fileAttr.Extr == ext {
				filesInfo = append(filesInfo, fileAttr)
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(filesInfo)
	if err != nil {
		log.Printf("can't convert data into json %v", err)
		http.Error(w, "Unable to read directory", http.StatusBadRequest)
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/upload", uploadHandleFunc)
	router.HandleFunc("/list", listHandleFunc)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Print(err)
	}
}
