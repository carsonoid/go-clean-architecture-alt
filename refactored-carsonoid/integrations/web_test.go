package integrations

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestWebIntegration_DownloadDocument(t *testing.T) {
	validDownloadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("test file contents"))
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer validDownloadServer.Close()

	invalidDownloadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(300)
	}))
	defer validDownloadServer.Close()

	type fields struct {
		client *http.Client
	}
	type args struct {
		ctx context.Context
		doc *Document
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantBody []byte
		wantErr  bool
	}{
		{
			name: "bad url",
			fields: fields{
				&http.Client{Timeout: time.Second * 3},
			},
			args: args{
				ctx: context.Background(),
				doc: &Document{
					DownloadURL: "",
				},
			},
			wantBody: nil,
			wantErr:  true,
		},
		{
			name: "downloadFail",
			fields: fields{
				&http.Client{Timeout: time.Second * 3},
			},
			args: args{
				ctx: context.Background(),
				doc: &Document{
					DownloadURL: invalidDownloadServer.URL,
				},
			},
			wantBody: nil,
			wantErr:  true,
		},
		{
			name: "success",
			fields: fields{
				&http.Client{Timeout: time.Second * 3},
			},
			args: args{
				ctx: context.Background(),
				doc: &Document{
					DownloadURL: validDownloadServer.URL,
				},
			},
			wantBody: []byte("test file contents"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &WebIntegration{
				client: tt.fields.client,
			}
			gotFile, err := i.DownloadDocument(tt.args.ctx, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("WebIntegration.DownloadDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotFile == nil {
				if tt.wantBody != nil {
					t.Errorf("WebIntegration.DownloadDocument() = %v, want %v", gotFile, tt.wantBody)
				}
				return
			}

			defer gotFile.Close()
			defer os.Remove(gotFile.Name())

			// read the body
			body, err := io.ReadAll(gotFile)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(body, tt.wantBody) {
				t.Errorf("WebIntegration.DownloadDocument() = %v, want %v", string(body), string(tt.wantBody))
			}
		})
	}
}
