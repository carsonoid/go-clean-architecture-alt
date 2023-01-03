package integrations

import (
	"context"
	"io"
	"os"
)

// ensure that NoopIntegration implements DocumentImporter and DocumentDownloader
var _ DocumentImporter = (*NoopIntegration)(nil)
var _ DocumentDownloader = (*NoopIntegration)(nil)

// NoopIntegration is an integration that will do nothing
type NoopIntegration struct{}

// NewNoopIntegration creates a new NoopIntegration
func NewNoopIntegration() *NoopIntegration {
	return &NoopIntegration{}
}

// DownloadDocument will return a file handle to a new, empty tempfile, the document is not actually
// inserted but the file handle is returned to simulate the download
func (i *NoopIntegration) DownloadDocument(ctx context.Context, doc *Document) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "noop-*.jpg")
	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}

// ImportDocument will do nothing, it will return nil to simulate the import as a noop
func (i *NoopIntegration) ImportDocument(ctx context.Context, doc *Document, file io.Reader) error {
	return nil
}
