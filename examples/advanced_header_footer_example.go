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

	// Add a first page header (appears only on first page)
	firstHeader := doc.AddHeader(stypes.HdrFtrFirst)
	firstHeader.AddParagraph("FIRST PAGE HEADER - Only on page 1")

	// Add a default header (appears on all pages except first if first page header exists)
	defaultHeader := doc.AddHeader(stypes.HdrFtrDefault)
	defaultHeader.AddParagraph("DEFAULT HEADER - Appears on all other pages")

	// Add an even page header (appears only on even numbered pages)
	evenHeader := doc.AddHeader(stypes.HdrFtrEven)
	evenHeader.AddParagraph("EVEN PAGE HEADER - Only on even pages")

	// Add a first page footer
	firstFooter := doc.AddFooter(stypes.HdrFtrFirst)
	firstFooter.AddParagraph("FIRST PAGE FOOTER - Only on page 1")

	// Add a default footer
	defaultFooter := doc.AddFooter(stypes.HdrFtrDefault)
	defaultFooter.AddParagraph("DEFAULT FOOTER - Appears on all other pages")

	// Add an even page footer
	evenFooter := doc.AddFooter(stypes.HdrFtrEven)
	evenFooter.AddParagraph("EVEN PAGE FOOTER - Only on even pages")

	// Save all headers and footers to the document
	if err := doc.SaveHeader(firstHeader); err != nil {
		log.Fatalf("Failed to save first header: %v", err)
	}
	if err := doc.SaveHeader(defaultHeader); err != nil {
		log.Fatalf("Failed to save default header: %v", err)
	}
	if err := doc.SaveHeader(evenHeader); err != nil {
		log.Fatalf("Failed to save even header: %v", err)
	}

	if err := doc.SaveFooter(firstFooter); err != nil {
		log.Fatalf("Failed to save first footer: %v", err)
	}
	if err := doc.SaveFooter(defaultFooter); err != nil {
		log.Fatalf("Failed to save default footer: %v", err)
	}
	if err := doc.SaveFooter(evenFooter); err != nil {
		log.Fatalf("Failed to save even footer: %v", err)
	}

	// Save the document
	if err := doc.SaveTo("advanced_header_footer_example.docx"); err != nil {
		log.Fatalf("Failed to save document: %v", err)
	}

	fmt.Println("Advanced document with multiple headers and footers saved successfully!")
}