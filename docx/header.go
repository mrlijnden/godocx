package docx

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/mrlijnden/godocx/common/constants"
	"github.com/mrlijnden/godocx/wml/ctypes"
	"github.com/mrlijnden/godocx/wml/stypes"
)

var headerAttrs = map[string]string{
	"xmlns:w": "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
	"xmlns:r": "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
}

// Header represents a document header with its content and properties
type Header struct {
	Root     *RootDoc
	Children []DocumentChild
	rID      int
	filename string
}

// NewHeader creates a new header instance
func NewHeader(root *RootDoc, filename string) *Header {
	return &Header{
		Root:     root,
		filename: filename,
	}
}

// AddParagraph adds a paragraph to the header
func (h *Header) AddParagraph(text string) *Paragraph {
	p := newParagraph(h.Root)
	p.AddText(text)
	h.Children = append(h.Children, DocumentChild{Para: p})
	return p
}

// AddEmptyParagraph adds an empty paragraph to the header
func (h *Header) AddEmptyParagraph() *Paragraph {
	p := newParagraph(h.Root)
	h.Children = append(h.Children, DocumentChild{Para: p})
	return p
}

// Update saves the current header content to the document's FileMap
func (h *Header) Update() error {
	return h.Root.SaveHeader(h)
}

// MarshalXML implements the xml.Marshaler interface for the Header type
func (h Header) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "w:hdr"

	for key, value := range headerAttrs {
		attr := xml.Attr{Name: xml.Name{Local: key}, Value: value}
		start.Attr = append(start.Attr, attr)
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for _, child := range h.Children {
		if child.Para != nil {
			if err = child.Para.ct.MarshalXML(e, xml.StartElement{}); err != nil {
				return err
			}
		}
		if child.Table != nil {
			if err = child.Table.ct.MarshalXML(e, xml.StartElement{}); err != nil {
				return err
			}
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// AddHeader adds a header to the document
func (rd *RootDoc) AddHeader(hdrType stypes.HdrFtrType) *Header {
	// Generate unique filename for header
	headerCount := rd.getHeaderCount() + 1
	filename := fmt.Sprintf("word/header%d.xml", headerCount)

	// Create header instance
	header := NewHeader(rd, filename)

	// Add relationship to document
	rID := rd.Document.addRelation(constants.SourceRelationshipHeader, "header"+strconv.Itoa(headerCount)+".xml")

	// Add content type override
	contentType := "application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"
	rd.ContentType.AddOverride("/"+filename, contentType)

	// Update section properties
	rd.ensureSectionProperties()
	rd.Document.Body.SectPr.HeaderReference = &ctypes.HeaderReference{
		Type: hdrType,
		ID:   rID,
	}

	// Store header for automatic serialization
	rd.Headers = append(rd.Headers, header)

	return header
}

// getHeaderCount returns the current number of headers in the document
func (rd *RootDoc) getHeaderCount() int {
	count := 0
	rd.FileMap.Range(func(key, value interface{}) bool {
		keyStr, ok := key.(string)
		if ok && len(keyStr) > 11 && keyStr[:11] == "word/header" {
			count++
		}
		return true
	})
	return count
}

// ensureSectionProperties ensures the document has section properties
func (rd *RootDoc) ensureSectionProperties() {
	if rd.Document.Body.SectPr == nil {
		rd.Document.Body.SectPr = ctypes.NewSectionProper()
	}
}

// SaveHeader saves the header XML to the document's FileMap
func (rd *RootDoc) SaveHeader(header *Header) error {
	xmlData, err := xml.Marshal(header)
	if err != nil {
		return err
	}

	// Add XML header
	fullXML := append(constants.XMLHeader, xmlData...)
	rd.FileMap.Store(header.filename, fullXML)

	return nil
}
