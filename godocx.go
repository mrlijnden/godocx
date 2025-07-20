package godocx

import (
	"github.com/mrlijnden/godocx/docx"
)

// NewDocument creates a new DOCX document and returns a pointer to the RootDoc.
// This is the main entry point for creating new documents.
func NewDocument() (*docx.RootDoc, error) {
	return docx.NewRootDoc(), nil
}

// OpenDocument opens an existing DOCX document from the specified path.
// This function reads and parses an existing DOCX file.
// TODO: Implement OpenDocument functionality
func OpenDocument(path string) (*docx.RootDoc, error) {
	// TODO: Implement document opening functionality
	return nil, nil
}
