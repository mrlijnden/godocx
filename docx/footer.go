package docx

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/mrlijnden/godocx/common/constants"
	"github.com/mrlijnden/godocx/wml/ctypes"
	"github.com/mrlijnden/godocx/wml/stypes"
)

var footerAttrs = map[string]string{
	"xmlns:w": "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
	"xmlns:r": "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
}

// Footer represents a document footer with its content and properties
type Footer struct {
	Root     *RootDoc
	Children []DocumentChild
	rID      int
	filename string
}

// NewFooter creates a new footer instance
func NewFooter(root *RootDoc, filename string) *Footer {
	return &Footer{
		Root:     root,
		filename: filename,
	}
}

// AddParagraph adds a paragraph to the footer
func (f *Footer) AddParagraph(text string) *Paragraph {
	p := newParagraph(f.Root)
	p.AddText(text)
	f.Children = append(f.Children, DocumentChild{Para: p})
	return p
}

// AddEmptyParagraph adds an empty paragraph to the footer
func (f *Footer) AddEmptyParagraph() *Paragraph {
	p := newParagraph(f.Root)
	f.Children = append(f.Children, DocumentChild{Para: p})
	return p
}

// Update saves the current footer content to the document's FileMap
func (f *Footer) Update() error {
	return f.Root.SaveFooter(f)
}

// MarshalXML implements the xml.Marshaler interface for the Footer type
func (f Footer) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "w:ftr"

	for key, value := range footerAttrs {
		attr := xml.Attr{Name: xml.Name{Local: key}, Value: value}
		start.Attr = append(start.Attr, attr)
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for _, child := range f.Children {
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

// AddFooter adds a footer to the document
func (rd *RootDoc) AddFooter(ftrType stypes.HdrFtrType) *Footer {
	// Generate unique filename for footer
	footerCount := rd.getFooterCount() + 1
	filename := fmt.Sprintf("word/footer%d.xml", footerCount)

	// Create footer instance
	footer := NewFooter(rd, filename)

	// Add relationship to document
	rID := rd.Document.addRelation(constants.SourceRelationshipFooter, "footer"+strconv.Itoa(footerCount)+".xml")

	// Add content type override
	contentType := "application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"
	rd.ContentType.AddOverride("/"+filename, contentType)

	// Update section properties
	rd.ensureSectionProperties()
	rd.Document.Body.SectPr.FooterReference = &ctypes.FooterReference{
		Type: ftrType,
		ID:   rID,
	}

	// Store footer for automatic serialization
	rd.Footers = append(rd.Footers, footer)

	return footer
}

// getFooterCount returns the current number of footers in the document
func (rd *RootDoc) getFooterCount() int {
	count := 0
	rd.FileMap.Range(func(key, value interface{}) bool {
		keyStr, ok := key.(string)
		if ok && len(keyStr) > 11 && keyStr[:11] == "word/footer" {
			count++
		}
		return true
	})
	return count
}

// SaveFooter saves the footer XML to the document's FileMap
func (rd *RootDoc) SaveFooter(footer *Footer) error {
	xmlData, err := xml.Marshal(footer)
	if err != nil {
		return err
	}

	// Add XML header
	fullXML := append(constants.XMLHeader, xmlData...)
	rd.FileMap.Store(footer.filename, fullXML)

	return nil
}
