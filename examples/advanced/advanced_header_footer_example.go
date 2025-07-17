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
	doc.AddParagraph("This document demonstrates different header and footer types.")
	doc.AddParagraph("Notice how the first page, even pages, and default pages can have different headers and footers.")

	// Add more content to create multiple pages
	for i := 1; i <= 5; i++ {
		doc.AddParagraph(fmt.Sprintf("This is paragraph %d to create multiple pages for testing.", i))
	}

	// Add page break
	doc.AddPageBreak()

	// Add more content on second page
	doc.AddParagraph("This is content on the second page.")
	doc.AddParagraph("The header and footer should appear here too.")

	// Add a first page header (no manual save needed!)
	firstHeader := doc.AddHeader(stypes.HdrFtrFirst)
	firstHeader.AddParagraph("FIRST PAGE HEADER - Only on page 1")

	// Add a default header (no manual save needed!)
	defaultHeader := doc.AddHeader(stypes.HdrFtrDefault)
	defaultHeader.AddParagraph("DEFAULT HEADER - Appears on all other pages")

	// Add an even page header (no manual save needed!)
	evenHeader := doc.AddHeader(stypes.HdrFtrEven)
	evenHeader.AddParagraph("EVEN PAGE HEADER - Only on even pages")

	// Add a first page footer (no manual save needed!)
	firstFooter := doc.AddFooter(stypes.HdrFtrFirst)
	firstFooter.AddParagraph("FIRST PAGE FOOTER - Only on page 1")

	// Add a default footer (no manual save needed!)
	defaultFooter := doc.AddFooter(stypes.HdrFtrDefault)
	defaultFooter.AddParagraph("DEFAULT FOOTER - Appears on all other pages")

	// Add an even page footer (no manual save needed!)
	evenFooter := doc.AddFooter(stypes.HdrFtrEven)
	evenFooter.AddParagraph("EVEN PAGE FOOTER - Only on even pages")

	// Save the document (all headers and footers are automatically included)
	if err := doc.SaveTo("advanced_header_footer_example.docx"); err != nil {
		log.Fatalf("Failed to save document: %v", err)
	}

	fmt.Println("Advanced document with multiple headers and footers saved successfully!")
	fmt.Println("Note: All headers and footers are automatically saved with the document!")
}
