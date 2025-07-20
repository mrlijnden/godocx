package ctypes

import (
	"encoding/xml"

	"github.com/mrlijnden/godocx/wml/stypes"
)

// FieldChar represents a field character in a document
type FieldChar struct {
	FieldCharType stypes.FieldCharType `xml:"fldCharType,attr"`
}

func NewFieldChar(fieldCharType stypes.FieldCharType) *FieldChar {
	return &FieldChar{
		FieldCharType: fieldCharType,
	}
}

func (fc FieldChar) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "w:fldChar"
	start.Attr = append(start.Attr, xml.Attr{
		Name:  xml.Name{Local: "w:fldCharType"},
		Value: string(fc.FieldCharType),
	})
	return e.EncodeElement("", start)
}

func (fc *FieldChar) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "fldCharType":
			fieldCharType, err := stypes.FieldCharTypeFromStr(attr.Value)
			if err != nil {
				return err
			}
			fc.FieldCharType = fieldCharType
		}
	}
	return d.Skip()
}
