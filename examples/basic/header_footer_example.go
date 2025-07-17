package main

import (
	"fmt"
	"log"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/wml/stypes"
)

func main() {
	// Create a new document
	doc, err := godocx.NewDocument()
	if err != nil {
		log.Fatalf("Failed to create new document: %v", err)
	}

	// Add document content
	doc.AddParagraph("This is the main document content.")
	doc.AddParagraph("Headers and footers will appear on every page.")

	// Add a default header (no manual save needed!)
	header := doc.AddHeader(stypes.HdrFtrDefault)
	header.AddParagraph("This is the default header - appears on all pages")

	// Add a default footer (no manual save needed!)
	footer := doc.AddFooter(stypes.HdrFtrDefault)
	footer.AddParagraph("This is the default footer - appears on all pages")

	// Save the document (headers and footers are automatically included)
	if err := doc.SaveTo("header_footer_example.docx"); err != nil {
		log.Fatalf("Failed to save document: %v", err)
	}

	fmt.Println("Document with header and footer saved successfully!")
	fmt.Println("Note: Headers and footers are automatically saved with the document!")
}
