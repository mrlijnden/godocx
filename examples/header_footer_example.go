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

	// Add a default header
	header := doc.AddHeader(stypes.HdrFtrDefault)
	header.AddParagraph("This is the default header - appears on all pages")

	// Add a default footer
	footer := doc.AddFooter(stypes.HdrFtrDefault)
	footer.AddParagraph("This is the default footer - appears on all pages")

	// Save the header and footer to the document
	if err := doc.SaveHeader(header); err != nil {
		log.Fatalf("Failed to save header: %v", err)
	}

	if err := doc.SaveFooter(footer); err != nil {
		log.Fatalf("Failed to save footer: %v", err)
	}

	// Save the document
	if err := doc.SaveTo("header_footer_example.docx"); err != nil {
		log.Fatalf("Failed to save document: %v", err)
	}

	fmt.Println("Document with header and footer saved successfully!")
}