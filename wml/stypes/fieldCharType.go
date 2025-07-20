package stypes

import (
	"encoding/xml"
	"errors"
)

// FieldCharType represents the type of field character
type FieldCharType string

const (
	FieldCharTypeBegin    FieldCharType = "begin"    // Field start character
	FieldCharTypeSeparate FieldCharType = "separate" // Field separator character
	FieldCharTypeEnd      FieldCharType = "end"      // Field end character
)

func FieldCharTypeFromStr(value string) (FieldCharType, error) {
	switch value {
	case "begin":
		return FieldCharTypeBegin, nil
	case "separate":
		return FieldCharTypeSeparate, nil
	case "end":
		return FieldCharTypeEnd, nil
	default:
		return "", errors.New("Invalid FieldCharType value")
	}
}

func (fct *FieldCharType) UnmarshalXMLAttr(attr xml.Attr) error {
	val, err := FieldCharTypeFromStr(attr.Value)
	if err != nil {
		return err
	}
	*fct = val
	return nil
}
