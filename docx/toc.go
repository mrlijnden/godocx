package docx

import (
	"fmt"
	"strings"

	"github.com/mrlijnden/godocx/wml/ctypes"
)

// TOC represents a Table of Contents in a document
type TOC struct {
	Title              string
	Entries            []TOCEntry
	MaxLevel           int
	MinLevel           int
	IncludePageNumbers bool
	Indentation        int
	root               *RootDoc
}

// TOCEntry represents a single entry in the table of contents
type TOCEntry struct {
	Text  string
	Level int
}

// AddTableOfContents creates a new Table of Contents in the document
// This follows the project's simple, direct API pattern
func (rd *RootDoc) AddTableOfContents() *TOC {
	toc := &TOC{
		Title:              "Table of Contents",
		MaxLevel:           3,
		MinLevel:           1,
		IncludePageNumbers: true,
		Indentation:        20,
		root:               rd,
	}

	// Get all headings
	headings, err := rd.getHeadingStructure()
	if err == nil {
		toc.Entries = headings
	}

	// Create TOC content paragraphs
	tocParagraphs := rd.createTOCContentParagraphs(toc)

	// Insert at the beginning
	rd.Document.Body.Children = append(tocParagraphs, rd.Document.Body.Children...)

	return toc
}

// SetTitle sets the title of the TOC
func (toc *TOC) SetTitle(title string) *TOC {
	toc.Title = title
	return toc
}

// SetMaxLevel sets the maximum heading level to include
func (toc *TOC) SetMaxLevel(level int) *TOC {
	toc.MaxLevel = level
	return toc
}

// SetMinLevel sets the minimum heading level to include
func (toc *TOC) SetMinLevel(level int) *TOC {
	toc.MinLevel = level
	return toc
}

// SetIncludePageNumbers sets whether to include page numbers
func (toc *TOC) SetIncludePageNumbers(include bool) *TOC {
	toc.IncludePageNumbers = include
	return toc
}

// SetIndentation sets the indentation for the TOC
func (toc *TOC) SetIndentation(indent int) *TOC {
	toc.Indentation = indent
	return toc
}

// getHeadingStructure retrieves all headings from the document
func (rd *RootDoc) getHeadingStructure() ([]TOCEntry, error) {
	var headings []TOCEntry

	if rd.Document == nil || rd.Document.Body == nil {
		return headings, nil
	}

	for _, child := range rd.Document.Body.Children {
		if child.Para != nil {
			heading := rd.extractHeadingFromParagraph(child.Para)
			if heading != nil {
				headings = append(headings, *heading)
			}
		}
	}

	return headings, nil
}

// extractHeadingFromParagraph extracts heading information from a paragraph
func (rd *RootDoc) extractHeadingFromParagraph(para *Paragraph) *TOCEntry {
	if para.ct.Property == nil || para.ct.Property.Style == nil {
		return nil
	}

	style := para.ct.Property.Style.Val
	text := rd.extractTextFromParagraph(para)

	// Check if it's a heading style
	if strings.HasPrefix(style, "Heading") {
		level := 1 // Default level
		if len(style) > 7 {
			levelStr := style[7:]
			if levelNum := int(levelStr[0] - '0'); levelNum >= 1 && levelNum <= 9 {
				level = levelNum
			}
		}
		return &TOCEntry{Text: text, Level: level}
	} else if style == "Title" {
		return &TOCEntry{Text: text, Level: 0}
	}

	return nil
}

// extractTextFromParagraph extracts all text from a paragraph
func (rd *RootDoc) extractTextFromParagraph(para *Paragraph) string {
	var text strings.Builder

	for _, child := range para.ct.Children {
		if child.Run != nil {
			for _, runChild := range child.Run.Children {
				if runChild.Text != nil {
					text.WriteString(runChild.Text.Text)
				}
			}
		}
	}

	return text.String()
}

// createTOCContentParagraphs creates the TOC content as paragraphs
func (rd *RootDoc) createTOCContentParagraphs(toc *TOC) []DocumentChild {
	var paragraphs []DocumentChild

	// Create TOC title paragraph with professional styling
	tocPara := newParagraph(rd)
	tocPara.AddText(toc.Title).Bold(true).Size(18)
	paragraphs = append(paragraphs, DocumentChild{Para: tocPara})

	// Add a blank line after title
	blankPara := newParagraph(rd)
	blankPara.AddText("")
	paragraphs = append(paragraphs, DocumentChild{Para: blankPara})

	// Add TOC entries as separate paragraphs
	for _, heading := range toc.Entries {
		// Skip if outside the level range
		if heading.Level < toc.MinLevel || heading.Level > toc.MaxLevel {
			continue
		}

		// Create TOC entry paragraph with professional formatting
		entryPara := newParagraph(rd)

		// Add proper indentation based on level
		indent := (heading.Level - toc.MinLevel) * toc.Indentation
		if indent > 0 {
			entryPara.Indent(&ctypes.Indent{Left: &indent})
		}

		// Add heading text with appropriate styling
		if heading.Level == 1 {
			entryPara.AddText(heading.Text).Bold(true)
		} else {
			entryPara.AddText(heading.Text)
		}

		// Add professional dotted leader with page number
		if toc.IncludePageNumbers {
			// Add space and dots for traditional TOC formatting
			entryPara.AddText(" ")

			// Add multiple dots for the leader
			for i := 0; i < 20; i++ {
				entryPara.AddText(".")
			}

			// Add page number
			pageNum := rd.calculatePageNumber(heading)
			entryPara.AddText(fmt.Sprintf(" %d", pageNum))
		}

		paragraphs = append(paragraphs, DocumentChild{Para: entryPara})
	}

	return paragraphs
}

// calculatePageNumber calculates the page number for a heading
func (rd *RootDoc) calculatePageNumber(heading TOCEntry) int {
	if rd.Document == nil || rd.Document.Body == nil {
		return 1
	}

	// Count paragraphs to estimate page breaks
	paragraphsPerPage := 20

	// Find the heading's position in the document
	headingIndex := rd.findHeadingIndex(heading.Text)
	if headingIndex == -1 {
		return 1
	}

	// Calculate page number based on position
	pageNumber := (headingIndex / paragraphsPerPage) + 1

	// Ensure minimum page number
	if pageNumber < 1 {
		pageNumber = 1
	}

	return pageNumber
}

// findHeadingIndex finds the index of a heading in the document
func (rd *RootDoc) findHeadingIndex(headingText string) int {
	if rd.Document == nil || rd.Document.Body == nil {
		return -1
	}

	for i, child := range rd.Document.Body.Children {
		if child.Para != nil {
			text := rd.extractTextFromParagraph(child.Para)
			if text == headingText {
				return i
			}
		}
	}

	return -1
}

// GetHeadingStructure returns the heading structure for external use
func (rd *RootDoc) GetHeadingStructure() ([]TOCEntry, error) {
	return rd.getHeadingStructure()
}
