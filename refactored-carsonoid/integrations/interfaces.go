package integrations

import (
	"context"
	"io"
	"os"
)

// A DocumentDownloader should download a document using metadata from the Document and return a file handle
type DocumentDownloader interface {
	DownloadDocument(ctx context.Context, doc *Document) (*os.File, error)
}

// A DocumentImporter should import a document using metadata from the Document and the file handle
type DocumentImporter interface {
	ImportDocument(ctx context.Context, doc *Document, file io.Reader) error
}
