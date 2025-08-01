package docx

import (
	"encoding/xml"
	"sync"

	"github.com/mrlijnden/godocx/wml/ctypes"
)

// RootDoc represents the root document of an Office Open XML (OOXML) document.
// It contains information about the document path, file map, the document structure,
// and relationships with other parts of the document.
type RootDoc struct {
	Path        string        // Path represents the path of the document.
	FileMap     sync.Map      // FileMap is a synchronized map for managing files related to the document.
	RootRels    Relationships // RootRels represents relationships at the root level.
	ContentType ContentTypes
	Document    *Document      // Document is the main document structure.
	DocStyles   *ctypes.Styles // Document styles

	rID        int // rId is used to generate unique relationship IDs.
	ImageCount uint

	// Headers and footers storage
	Headers []*Header // Headers stores all headers for automatic serialization
	Footers []*Footer // Footers stores all footers for automatic serialization
}

// NewRootDoc creates a new instance of the RootDoc structure.
func NewRootDoc() *RootDoc {
	rd := &RootDoc{}
	rd.Document = &Document{
		Root: rd,
		Body: NewBody(rd),
	}
	rd.DocStyles = &ctypes.Styles{}
	rd.ContentType = ContentTypes{}
	rd.RootRels = Relationships{}
	rd.Document.DocRels = Relationships{}
	return rd
}

// LoadDocXml decodes the provided XML data and returns a Document instance.
// It is used to load the main document structure from the document file.
//
// Parameters:
//   - fileName: The name of the document file.
//   - fileBytes: The XML data representing the main document structure.
//
// Returns:
//   - doc: The Document instance containing the decoded main document structure.
//   - err: An error, if any occurred during the decoding process.
func LoadDocXml(rd *RootDoc, fileName string, fileBytes []byte) (*Document, error) {
	doc := Document{
		Root: rd,
	}
	err := xml.Unmarshal(fileBytes, &doc)
	if err != nil {
		return nil, err
	}

	doc.relativePath = fileName
	return &doc, nil
}

// Load styles.xml into Styles struct
func LoadStyles(fileName string, fileBytes []byte) (*ctypes.Styles, error) {
	styles := ctypes.Styles{}
	err := xml.Unmarshal(fileBytes, &styles)
	if err != nil {
		return nil, err
	}

	styles.RelativePath = fileName
	return &styles, nil
}
