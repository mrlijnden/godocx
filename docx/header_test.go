package docx

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/gomutex/godocx/wml/ctypes"
	"github.com/gomutex/godocx/wml/stypes"
)

func TestHeader_MarshalXML(t *testing.T) {
	tests := []struct {
		name     string
		input    Header
		expected []string
	}{
		{
			name: "Empty Header",
			input: Header{
				Children: []DocumentChild{},
			},
			expected: []string{
				`<w:hdr`,
				`xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"`,
				`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`,
				`</w:hdr>`,
			},
		},
		{
			name: "Header with Paragraph",
			input: Header{
				Children: []DocumentChild{
					{Para: newTestParagraph("Test Header")},
				},
			},
			expected: []string{
				`<w:hdr`,
				`xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"`,
				`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`,
				`</w:hdr>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result strings.Builder
			encoder := xml.NewEncoder(&result)
			start := xml.StartElement{Name: xml.Name{Local: "w:hdr"}}

			err := tt.input.MarshalXML(encoder, start)
			if err != nil {
				t.Fatalf("Error marshaling XML: %v", err)
			}

			encoder.Flush()
			actual := result.String()

			for _, exp := range tt.expected {
				if !strings.Contains(actual, exp) {
					t.Errorf("Expected XML part not found in actual XML:\nExpected part: %s\nActual XML: %s", exp, actual)
				}
			}
		})
	}
}

func TestRootDoc_AddHeader(t *testing.T) {
	rd := NewRootDoc()
	rd.Document = &Document{
		Root: rd,
		Body: NewBody(rd),
	}
	rd.Document.Body.SectPr = nil // Start with no section properties

	// Test adding a default header
	header := rd.AddHeader(stypes.HdrFtrDefault)

	if header == nil {
		t.Fatal("AddHeader returned nil")
	}

	if header.Root != rd {
		t.Error("Header root reference is incorrect")
	}

	if strings.Contains(header.filename, "header1.xml") == false {
		t.Errorf("Expected header filename to contain 'header1.xml', got: %s", header.filename)
	}

	// Check that section properties were created
	if rd.Document.Body.SectPr == nil {
		t.Error("Section properties were not created")
	}

	// Check that header reference was added
	if rd.Document.Body.SectPr.HeaderReference == nil {
		t.Error("Header reference was not added to section properties")
	}

	if rd.Document.Body.SectPr.HeaderReference.Type != stypes.HdrFtrDefault {
		t.Error("Header reference type is incorrect")
	}
}

func TestHeader_AddParagraph(t *testing.T) {
	rd := NewRootDoc()
	header := NewHeader(rd, "word/header1.xml")

	para := header.AddParagraph("Test Header Text")

	if para == nil {
		t.Fatal("AddParagraph returned nil")
	}

	if len(header.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(header.Children))
	}

	if header.Children[0].Para != para {
		t.Error("Paragraph not properly added to header children")
	}
}

func TestHeader_AddEmptyParagraph(t *testing.T) {
	rd := NewRootDoc()
	header := NewHeader(rd, "word/header1.xml")

	para := header.AddEmptyParagraph()

	if para == nil {
		t.Fatal("AddEmptyParagraph returned nil")
	}

	if len(header.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(header.Children))
	}

	if header.Children[0].Para != para {
		t.Error("Empty paragraph not properly added to header children")
	}
}

func TestRootDoc_SaveHeader(t *testing.T) {
	rd := NewRootDoc()
	header := NewHeader(rd, "word/header1.xml")
	header.AddParagraph("Test Header")

	err := rd.SaveHeader(header)
	if err != nil {
		t.Fatalf("SaveHeader failed: %v", err)
	}

	// Check that header was saved to FileMap
	value, ok := rd.FileMap.Load("word/header1.xml")
	if !ok {
		t.Error("Header was not saved to FileMap")
	}

	xmlData, ok := value.([]byte)
	if !ok {
		t.Error("Header data in FileMap is not []byte")
	}

	xmlString := string(xmlData)
	if !strings.Contains(xmlString, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("Header XML does not contain XML declaration")
	}

	if !strings.Contains(xmlString, "<w:hdr") {
		t.Error("Header XML does not contain header element")
	}
}

// Helper function to create a test paragraph
func newTestParagraph(text string) *Paragraph {
	return &Paragraph{
		ct: ctypes.Paragraph{},
	}
}