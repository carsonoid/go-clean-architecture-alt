package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bradleyshawkins/go-clean-architecture/refactored-carsonoid/integrations"
)

type fakeDocumentImporter struct {
	importErr error
}

func (f *fakeDocumentImporter) ImportDocument(ctx context.Context, doc *integrations.Document, r io.Reader) error {
	return f.importErr
}

type fakeDocumentDownloader struct {
	downloadResp *os.File
	downloadErr  error
}

func (f *fakeDocumentDownloader) DownloadDocument(ctx context.Context, doc *integrations.Document) (*os.File, error) {
	return f.downloadResp, f.downloadErr
}

func TestFileDownloader_handler(t *testing.T) {
	type fields struct {
		DocumentImporter   integrations.DocumentImporter
		DocumentDownloader integrations.DocumentDownloader
		insertsEnabled     bool
	}
	tests := []struct {
		name       string
		fields     fields
		newRequest func() *http.Request
		wantCode   int
	}{
		{
			name: "inserts are disabled",
			fields: fields{
				insertsEnabled: false,
			},
			newRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/download", nil)
			},
			wantCode: http.StatusServiceUnavailable,
		},
		{
			name: "bad request body",
			fields: fields{
				insertsEnabled: true,
			},
			newRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/download", nil)
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "download error",
			fields: fields{
				insertsEnabled:     true,
				DocumentImporter:   &fakeDocumentImporter{},
				DocumentDownloader: &fakeDocumentDownloader{downloadErr: errors.New("download error")},
			},
			newRequest: func() *http.Request {
				body, err := json.Marshal(integrations.Document{
					DownloadURL: "http://example.com",
				})
				if err != nil {
					t.Fatal(err)
				}

				return httptest.NewRequest(http.MethodPost, "/download", bytes.NewBuffer(body))
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "import error",
			fields: fields{
				insertsEnabled:     true,
				DocumentImporter:   &fakeDocumentImporter{importErr: errors.New("import error")},
				DocumentDownloader: &fakeDocumentDownloader{},
			},
			newRequest: func() *http.Request {
				body, err := json.Marshal(integrations.Document{
					DownloadURL: "http://example.com",
				})
				if err != nil {
					t.Fatal(err)
				}

				return httptest.NewRequest(http.MethodPost, "/download", bytes.NewBuffer(body))
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			fields: fields{
				insertsEnabled:     true,
				DocumentImporter:   &fakeDocumentImporter{},
				DocumentDownloader: &fakeDocumentDownloader{},
			},
			newRequest: func() *http.Request {
				body, err := json.Marshal(integrations.Document{
					DownloadURL: "http://example.com",
				})
				if err != nil {
					t.Fatal(err)
				}

				return httptest.NewRequest(http.MethodPost, "/download", bytes.NewBuffer(body))
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDownloader{
				DocumentImporter:   tt.fields.DocumentImporter,
				DocumentDownloader: tt.fields.DocumentDownloader,
				insertsEnabled:     tt.fields.insertsEnabled,
			}

			rec := httptest.NewRecorder()

			f.handler(rec, tt.newRequest())

			if rec.Code != tt.wantCode {
				t.Errorf("FileDownloader.handler() code = %v, want %v", rec.Code, tt.wantCode)
			}
		})
	}
}
