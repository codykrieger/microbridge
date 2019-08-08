package xmlrpc

import (
	"encoding/xml"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type MethodCall struct {
	XMLName    xml.Name      `xml:"methodCall"`
	MethodName string        `xml:"methodName"`
	Params     []XMLRPCValue `xml:"params>param>value"`
}

type XMLRPCValue struct {
	XMLName xml.Name `xml:"value"`
	Value   interface{}
}

type XMLRPCStruct struct {
	XMLName xml.Name             `xml:"struct"`
	Members []XMLRPCStructMember `xml:"member"`
}

type XMLRPCStructMember struct {
	XMLName xml.Name    `xml:"member"`
	Name    string      `xml:"name"`
	Value   XMLRPCValue `xml:"value"`
}

func (v *XMLRPCValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			log.WithError(err).Error("error retrieving first token while unmarshalling XMLRPCValue")
			return err
		}

		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "array":
				var s []XMLRPCValue
				if err := d.DecodeElement(&s, &el); err != nil {
					return err
				}
				v.Value = s
			case "struct":
				var s XMLRPCStruct
				if err := d.DecodeElement(&s, &el); err != nil {
					return err
				}
				v.Value = s
			case "string":
				var s string
				if err := d.DecodeElement(&s, &el); err != nil {
					return err
				}
				v.Value = s
			case "int":
				var i int
				if err := d.DecodeElement(&i, &el); err != nil {
					return err
				}
				v.Value = i
			case "boolean":
				var b bool
				if err := d.DecodeElement(&b, &el); err != nil {
					return err
				}
				v.Value = b
			case "dateTime.iso8601":
				var t time.Time
				if err := d.DecodeElement(&t, &el); err != nil {
					return err
				}
				v.Value = t
			default:
				return fmt.Errorf("unknown XMLRPCValue type '%s'", el.Name.Local)
			}
		case xml.EndElement:
			if el == start.End() {
				return nil
			}
		}
	}
	return nil
}
