# Table of Contents (TOC) Example

This example demonstrates how to create a professional Table of Contents in a DOCX document using the godocx library.

## Features

- **Professional formatting** with dotted leaders and proper indentation
- **Automatic heading detection** from document content
- **Configurable levels** (Heading1, Heading2, Heading3, etc.)
- **Page number calculation** for realistic TOC entries
- **Immediate visibility** - no manual Word updates required
- **Fluent API** for easy configuration

## Usage

```go
package main

import (
    "log"
    "github.com/mrlijnden/godocx"
)

func main() {
    // Create a new document
    doc, err := godocx.NewDocument()
    if err != nil {
        log.Fatal(err)
    }

    // Add content with headings
    doc.AddHeading("Introduction", 1)
    doc.AddParagraph("This is the introduction...")
    
    doc.AddHeading("Getting Started", 1)
    doc.AddParagraph("To get started...")
    
    doc.AddHeading("Installation", 2)
    doc.AddParagraph("Install the library...")

    // Add TOC with fluent API configuration
    toc := doc.AddTableOfContents().
        SetTitle("Table of Contents").
        SetMaxLevel(3).
        SetMinLevel(1).
        SetIncludePageNumbers(true).
        SetIndentation(20)

    // Save the document
    err = doc.SaveTo("document_with_toc.docx")
    if err != nil {
        log.Fatal(err)
    }
}
```

## API Reference

### `AddTableOfContents()`

Creates a new Table of Contents in the document. Returns a `*TOC` object for configuration.

### TOC Configuration Methods

All configuration methods return the TOC object for fluent chaining:

- **`SetTitle(title string)`** - Sets the TOC title (default: "Table of Contents")
- **`SetMaxLevel(level int)`** - Sets maximum heading level to include (default: 3)
- **`SetMinLevel(level int)`** - Sets minimum heading level to include (default: 1)
- **`SetIncludePageNumbers(include bool)`** - Enables/disables page numbers (default: true)
- **`SetIndentation(indent int)`** - Sets indentation per level (default: 20)

### Example Output

The generated TOC will look like:

```
TABLE OF CONTENTS

Introduction ................... 1
Getting Started ................ 1
  Installation ................. 2
  Configuration ................ 2
Advanced Usage ................. 3
  Troubleshooting .............. 3
Conclusion ..................... 4
```

## How It Works

1. **Heading Detection**: The TOC automatically scans the document for paragraphs with heading styles (Heading1, Heading2, etc.)
2. **Content Generation**: Creates professional-looking TOC entries with proper formatting
3. **Page Calculation**: Estimates page numbers based on document structure
4. **Document Insertion**: Places the TOC at the beginning of the document

## Notes

- The TOC is immediately visible in the generated DOCX file
- No manual updates are required in Microsoft Word
- Page numbers are estimated based on document structure
- The API follows the project's simple, direct pattern with fluent configuration 