package docx

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/mrlijnden/godocx/wml/stypes"
)

func TestFooter_MarshalXML(t *testing.T) {
	tests := []struct {
		name     string
		input    Footer
		expected []string
	}{
		{
			name: "Empty Footer",
			input: Footer{
				Children: []DocumentChild{},
			},
			expected: []string{
				`<w:ftr`,
				`xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"`,
				`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`,
				`</w:ftr>`,
			},
		},
		{
			name: "Footer with Paragraph",
			input: Footer{
				Children: []DocumentChild{
					{Para: newTestParagraph("Test Footer")},
				},
			},
			expected: []string{
				`<w:ftr`,
				`xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"`,
				`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`,
				`</w:ftr>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result strings.Builder
			encoder := xml.NewEncoder(&result)
			start := xml.StartElement{Name: xml.Name{Local: "w:ftr"}}

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

func TestRootDoc_AddFooter(t *testing.T) {
	rd := NewRootDoc()
	rd.Document = &Document{
		Root: rd,
		Body: NewBody(rd),
	}
	rd.Document.Body.SectPr = nil // Start with no section properties

	// Test adding a default footer
	footer := rd.AddFooter(stypes.HdrFtrDefault)

	if footer == nil {
		t.Fatal("AddFooter returned nil")
	}

	if footer.Root != rd {
		t.Error("Footer root reference is incorrect")
	}

	if strings.Contains(footer.filename, "footer1.xml") == false {
		t.Errorf("Expected footer filename to contain 'footer1.xml', got: %s", footer.filename)
	}

	// Check that section properties were created
	if rd.Document.Body.SectPr == nil {
		t.Error("Section properties were not created")
	}

	// Check that footer reference was added
	if rd.Document.Body.SectPr.FooterReference == nil {
		t.Error("Footer reference was not added to section properties")
	}

	if rd.Document.Body.SectPr.FooterReference.Type != stypes.HdrFtrDefault {
		t.Error("Footer reference type is incorrect")
	}
}

func TestFooter_AddParagraph(t *testing.T) {
	rd := NewRootDoc()
	footer := NewFooter(rd, "word/footer1.xml")

	para := footer.AddParagraph("Test Footer Text")

	if para == nil {
		t.Fatal("AddParagraph returned nil")
	}

	if len(footer.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(footer.Children))
	}

	if footer.Children[0].Para != para {
		t.Error("Paragraph not properly added to footer children")
	}
}

func TestFooter_AddEmptyParagraph(t *testing.T) {
	rd := NewRootDoc()
	footer := NewFooter(rd, "word/footer1.xml")

	para := footer.AddEmptyParagraph()

	if para == nil {
		t.Fatal("AddEmptyParagraph returned nil")
	}

	if len(footer.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(footer.Children))
	}

	if footer.Children[0].Para != para {
		t.Error("Empty paragraph not properly added to footer children")
	}
}

func TestRootDoc_SaveFooter(t *testing.T) {
	rd := NewRootDoc()
	footer := NewFooter(rd, "word/footer1.xml")
	footer.AddParagraph("Test Footer")

	err := rd.SaveFooter(footer)
	if err != nil {
		t.Fatalf("SaveFooter failed: %v", err)
	}

	// Check that footer was saved to FileMap
	value, ok := rd.FileMap.Load("word/footer1.xml")
	if !ok {
		t.Error("Footer was not saved to FileMap")
	}

	xmlData, ok := value.([]byte)
	if !ok {
		t.Error("Footer data in FileMap is not []byte")
	}

	xmlString := string(xmlData)
	if !strings.Contains(xmlString, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("Footer XML does not contain XML declaration")
	}

	if !strings.Contains(xmlString, "<w:ftr") {
		t.Error("Footer XML does not contain footer element")
	}
}

func TestRootDoc_getFooterCount(t *testing.T) {
	rd := NewRootDoc()

	// Initially should be 0
	if count := rd.getFooterCount(); count != 0 {
		t.Errorf("Expected initial footer count to be 0, got %d", count)
	}

	// Add a footer file to FileMap
	rd.FileMap.Store("word/footer1.xml", []byte("<w:ftr></w:ftr>"))

	// Should now be 1
	if count := rd.getFooterCount(); count != 1 {
		t.Errorf("Expected footer count to be 1, got %d", count)
	}

	// Add another footer file
	rd.FileMap.Store("word/footer2.xml", []byte("<w:ftr></w:ftr>"))

	// Should now be 2
	if count := rd.getFooterCount(); count != 2 {
		t.Errorf("Expected footer count to be 2, got %d", count)
	}

	// Add a non-footer file (should not affect count)
	rd.FileMap.Store("word/document.xml", []byte("<w:document></w:document>"))

	// Should still be 2
	if count := rd.getFooterCount(); count != 2 {
		t.Errorf("Expected footer count to still be 2, got %d", count)
	}
}