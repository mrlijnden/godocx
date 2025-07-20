package docx

import (
	"fmt"
	"strings"

	"github.com/mrlijnden/godocx/wml/ctypes"
)

// TOCEntry represents a single entry in the table of contents
type TOCEntry struct {
	Text        string
	Level       int
	PageNumber  int
	BookmarkID  string
	Style       string
	Indentation int
}

// TOC represents a complete table of contents
type TOC struct {
	Title    string
	Entries  []TOCEntry
	Options  TOCOptions
	Style    TOCStyle
	Position TOCPosition
	root     *RootDoc
}

// TOCOptions defines the configuration for a TOC
type TOCOptions struct {
	IncludePageNumbers bool
	MaxLevel           int
	MinLevel           int
}

// TOCStyle defines the styling for a TOC
type TOCStyle struct {
	TitleStyle  string
	Indentation int
}

// TOCPosition defines where the TOC should be placed
type TOCPosition struct {
	Location    TOCLocation // Beginning, End, Custom
	CustomIndex int         // For custom placement
}

type TOCLocation string

const (
	TOCLocationBeginning TOCLocation = "beginning"
	TOCLocationEnd       TOCLocation = "end"
	TOCLocationCustom    TOCLocation = "custom"
)

// DefaultTOCOptions returns default TOC options
func DefaultTOCOptions() TOCOptions {
	return TOCOptions{
		IncludePageNumbers: true,
		MaxLevel:           3,
		MinLevel:           1,
	}
}

// DefaultTOCStyle returns default TOC styling
func DefaultTOCStyle() TOCStyle {
	return TOCStyle{
		TitleStyle:  "TOC Title",
		Indentation: 20,
	}
}

// SetTitle sets the title for the TOC
func (toc *TOC) SetTitle(title string) *TOC {
	toc.Title = title
	return toc
}

// SetStyle sets the styling for the TOC
func (toc *TOC) SetStyle(style TOCStyle) *TOC {
	toc.Style = style
	return toc
}

// SetOptions sets the options for the TOC
func (toc *TOC) SetOptions(options TOCOptions) *TOC {
	toc.Options = options
	return toc
}

// GetHeadingStructure scans the document and returns all headings as TOC entries
func (rd *RootDoc) GetHeadingStructure() ([]TOCEntry, error) {
	return rd.getHeadingStructure()
}

// getHeadingStructure scans the document and returns all headings as TOC entries
func (rd *RootDoc) getHeadingStructure() ([]TOCEntry, error) {
	var entries []TOCEntry

	if rd.Document == nil || rd.Document.Body == nil {
		return entries, nil
	}

	for _, child := range rd.Document.Body.Children {
		if child.Para != nil {
			entry, err := rd.extractHeadingFromParagraph(child.Para)
			if err != nil {
				return nil, err
			}
			if entry != nil {
				entries = append(entries, *entry)
			}
		}
	}

	return entries, nil
}

// extractHeadingFromParagraph extracts heading information from a paragraph
func (rd *RootDoc) extractHeadingFromParagraph(para *Paragraph) (*TOCEntry, error) {
	ct := para.GetCT()
	if ct == nil || ct.Property == nil || ct.Property.Style == nil {
		return nil, nil // Not a heading
	}

	style := ct.Property.Style.Val
	if !strings.HasPrefix(style, "Heading") {
		return nil, nil // Not a heading style
	}

	// Extract heading level
	level := rd.extractHeadingLevel(style)
	if level == 0 {
		return nil, nil // Invalid heading level
	}

	// Extract heading text
	text := rd.extractTextFromParagraph(para)
	if text == "" {
		return nil, nil // Empty heading
	}

	entry := &TOCEntry{
		Text:        text,
		Level:       level,
		PageNumber:  0, // Will be calculated later
		BookmarkID:  rd.generateBookmarkID(text),
		Style:       style,
		Indentation: level * 20, // Default indentation
	}

	return entry, nil
}

// extractHeadingLevel extracts the heading level from a style name
func (rd *RootDoc) extractHeadingLevel(style string) int {
	if !strings.HasPrefix(style, "Heading") {
		return 0
	}

	levelStr := strings.TrimPrefix(style, "Heading")
	if len(levelStr) == 0 {
		return 1 // Default to level 1
	}

	// Try to parse the level number
	var level int
	_, err := fmt.Sscanf(levelStr, "%d", &level)
	if err != nil {
		return 1 // Default to level 1 if parsing fails
	}

	if level < 1 || level > 9 {
		return 1 // Default to level 1 if out of range
	}

	return level
}

// extractTextFromParagraph extracts text from a paragraph
func (rd *RootDoc) extractTextFromParagraph(para *Paragraph) string {
	var textBuilder strings.Builder

	ct := para.GetCT()
	if ct == nil {
		return ""
	}

	for _, child := range ct.Children {
		if child.Run != nil {
			for _, runChild := range child.Run.Children {
				if runChild.Text != nil {
					textBuilder.WriteString(runChild.Text.Text)
				}
			}
		}
	}

	return textBuilder.String()
}

// generateBookmarkID generates a bookmark ID for a heading
func (rd *RootDoc) generateBookmarkID(text string) string {
	// Simple bookmark ID generation - can be enhanced later
	// Remove special characters and replace spaces with underscores
	cleanText := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		if r == ' ' {
			return '_'
		}
		return -1
	}, text)

	// Limit length and ensure uniqueness
	if len(cleanText) > 40 {
		cleanText = cleanText[:40]
	}

	return "_TOC_" + cleanText
}

// createTOCContentParagraphs creates multiple paragraphs for TOC content
func (rd *RootDoc) createTOCContentParagraphs(headings []TOCEntry, options TOCOptions) []DocumentChild {
	var paragraphs []DocumentChild

	// Create TOC title paragraph with professional styling
	tocPara := newParagraph(rd)
	tocPara.AddText("TABLE OF CONTENTS").Bold(true).Size(18)
	paragraphs = append(paragraphs, DocumentChild{Para: tocPara})

	// Add a blank line after title
	blankPara := newParagraph(rd)
	blankPara.AddText("")
	paragraphs = append(paragraphs, DocumentChild{Para: blankPara})

	// Add TOC entries as separate paragraphs
	for _, heading := range headings {
		// Skip if outside the level range
		if heading.Level < options.MinLevel || heading.Level > options.MaxLevel {
			continue
		}

		// Create TOC entry paragraph with professional formatting
		entryPara := newParagraph(rd)

		// Add proper indentation based on level (more realistic)
		indent := (heading.Level - options.MinLevel) * 36 // More realistic indentation
		if indent > 0 {
			entryPara.Indent(&ctypes.Indent{Left: &indent})
		}

		// Add heading text with appropriate styling
		if heading.Level == 1 {
			entryPara.AddText(heading.Text).Bold(true)
		} else {
			entryPara.AddText(heading.Text)
		}

		// Add professional dotted leader with tab stop
		if options.IncludePageNumbers {
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

// calculatePageNumber calculates the page number for a heading (enhanced)
func (rd *RootDoc) calculatePageNumber(heading TOCEntry) int {
	// Enhanced page number calculation
	// This is still simplified but more realistic than the previous version

	if rd.Document == nil || rd.Document.Body == nil {
		return 1
	}

	// Count paragraphs to estimate page breaks
	// In a real document, approximately 20-25 paragraphs per page
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

// AddTableOfContentsProgrammatic creates a clean, simple TOC with programmatic content
// This is the recommended method for most use cases
func (rd *RootDoc) AddTableOfContentsProgrammatic(options TOCOptions) (*TOC, error) {
	// Create TOC instance
	toc := &TOC{
		Title:    "Table of Contents",
		Entries:  []TOCEntry{},
		Options:  options,
		Style:    DefaultTOCStyle(),
		Position: TOCPosition{Location: TOCLocationBeginning},
		root:     rd,
	}

	// Get all headings
	headings, err := rd.getHeadingStructure()
	if err != nil {
		return nil, fmt.Errorf("failed to get heading structure: %w", err)
	}
	toc.Entries = headings

	// Create TOC content paragraphs
	tocParagraphs := rd.createTOCContentParagraphs(headings, options)

	// Insert at the beginning
	rd.Document.Body.Children = append(tocParagraphs, rd.Document.Body.Children...)

	return toc, nil
}
