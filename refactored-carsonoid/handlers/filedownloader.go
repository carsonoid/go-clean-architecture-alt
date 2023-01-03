package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bradleyshawkins/go-clean-architecture/refactored-carsonoid/integrations"
)

// FileDownloader is a handler that will download a file from a url and import it into a system
type FileDownloader struct {
	integrations.DocumentImporter
	integrations.DocumentDownloader

	insertsEnabled bool
}

func NewFileDownloader(downloader integrations.DocumentDownloader, importer integrations.DocumentImporter) http.HandlerFunc {
	return (&FileDownloader{
		DocumentDownloader: downloader,
		DocumentImporter:   importer,
		insertsEnabled:     true,
	}).handler
}

func (f *FileDownloader) handler(w http.ResponseWriter, r *http.Request) {
	if !f.canInsertDocument() {
		// close signals that we are done with the request and we are not going to read it
		_ = r.Body.Close()
		http.Error(w, "inserts are disabled", http.StatusServiceUnavailable)
		return
	}

	doc := &integrations.Document{}
	err := json.NewDecoder(r.Body).Decode(doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received request to download file from:", doc.DownloadURL)

	reader, err := f.DownloadDocument(r.Context(), doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = f.ImportDocument(r.Context(), doc, reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func (f *FileDownloader) canInsertDocument() bool {
	return f.insertsEnabled
}
