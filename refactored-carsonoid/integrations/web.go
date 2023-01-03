package integrations

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// ensure that WebIntegration implements DocumentImporter
var _ DocumentDownloader = (*WebIntegration)(nil)

// WebIntegration is an integration that will download documents from a url
type WebIntegration struct {
	client *http.Client
}

// NewWebIntegration creates a new WebIntegration, it will use the
// provided http.Client if provided. Otherwise it will create a new one
//
// > If the client has no Timeout set, it will be set to 10 seconds
func NewWebIntegration(client *http.Client) *WebIntegration {
	if client == nil {
		client = &http.Client{}
	}

	if client.Timeout <= 0 {
		client.Timeout = 10 * time.Second
	}

	return &WebIntegration{
		client: client,
	}
}

// DownloadDocument will download the document from the provided url to a tmpdir and return the file handle
// it is the callers responsibility to close the file handle and delete the file
func (i *WebIntegration) DownloadDocument(ctx context.Context, doc *Document) (*os.File, error) {
	fmt.Println("Downloading a document with the web integration")
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.jpg", doc.PatientID))
	if err != nil {
		return nil, fmt.Errorf("unable to create tempFile. %w", err)
	}

	u, err := i.getDocumentURL(doc.DownloadURL)
	if err != nil {
		return nil, err
	}

	fmt.Println("Downloading file...")
	err = i.downloadFile(ctx, u, tmpFile)
	if err != nil {
		return nil, err
	}

	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("unable to seek to beginning of temp file. %w", err)
	}

	return tmpFile, nil
}

func (i *WebIntegration) downloadFile(ctx context.Context, u *url.URL, w io.Writer) error {
	req, err := http.NewRequest(http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("unable to create request to download file. %w", err)
	}

	resp, err := i.client.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("unable to make request to download file. %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("invalid status code received. %w StatusCode: %d", err, resp.StatusCode)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to copy payload to writer. %w", err)
	}

	return nil
}

func (i *WebIntegration) getDocumentURL(downloadURL string) (*url.URL, error) {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("invalid download url provided. %w", err)
	}
	return u, nil
}
