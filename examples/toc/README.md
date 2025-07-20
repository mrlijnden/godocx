# Table of Contents (TOC) Example

This directory contains a clean, simple example of how to create a Table of Contents in DOCX documents.

## 📁 Files

- `toc_example.go` - Simple TOC example with programmatic content generation

## 🎯 Usage

### Basic TOC Creation

```go
package main

import (
    "log"
    "github.com/mrlijnden/godocx"
    "github.com/mrlijnden/godocx/docx"
)

func main() {
    // Create a new document
    doc, err := godocx.NewDocument()
    if err != nil {
        log.Fatal(err)
    }

    // Add content with headings
    doc.AddHeading("Introduction", 1)
    doc.AddParagraph("This is the introduction section.")

    doc.AddHeading("Getting Started", 1)
    doc.AddParagraph("This section explains how to get started.")

    doc.AddHeading("Installation", 2)
    doc.AddParagraph("Installation instructions go here.")

    // Create TOC options
    tocOptions := docx.TOCOptions{
        IncludePageNumbers: true,
        MaxLevel:           3, // Include Heading1, Heading2, Heading3
        MinLevel:           1,
    }

    // Add programmatic TOC
    toc, err := doc.AddTableOfContentsProgrammatic(tocOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Save the document
    err = doc.SaveTo("document_with_toc.docx")
    if err != nil {
        log.Fatal(err)
    }
}
```

## ✨ Features

- **Professional formatting** with dotted leaders
- **Bold Level 1 headings** for emphasis
- **Proper indentation** for hierarchy
- **Realistic page numbers** based on content position
- **Immediate visibility** - no manual update needed
- **Clean, simple API** - easy to use

## 🎯 TOC Options

```go
type TOCOptions struct {
    IncludePageNumbers bool // Show page numbers
    MaxLevel           int  // Maximum heading level to include (1-9)
    MinLevel           int  // Minimum heading level to include (1-9)
}
```

## 🚀 Running the Example

```bash
go run toc_example.go
```

This will create a `toc_example.docx` file with a professional table of contents.

## 📋 Result

The generated TOC will look like:

```
TABLE OF CONTENTS

Introduction .................... 1
Getting Started ................ 1
  Installation ................. 1
  Configuration ................ 1
Advanced Usage ................. 1
  Troubleshooting ............. 1
Conclusion ..................... 1
```

## 🎉 Benefits

- **No manual updates** - TOC is visible immediately
- **Professional appearance** - Looks like a real table of contents
- **Automatic detection** - Scans document for heading styles
- **Flexible configuration** - Customize levels and options
- **Production ready** - Clean, simple implementation 