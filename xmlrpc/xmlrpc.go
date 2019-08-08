package xmlrpc

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	rpc "github.com/gorilla/rpc/v2"
	log "github.com/sirupsen/logrus"
)

type Codec struct {
	AutoCapitalizeMethodName bool
}

func NewCodec() *Codec {
	return &Codec{}
}

func (c *Codec) NewRequest(req *http.Request) rpc.CodecRequest {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.WithError(err).Error("xmlrpc: error reading request body")
		return nil
	}

	var methodCall struct {
		MethodName string `xml:"methodName"`
	}

	if err = xml.Unmarshal(data, &methodCall); err != nil {
		log.WithError(err).Error("xmlrpc: error unmarshalling request method name")
		return nil
	}

	return &CodecRequest{
		req:                      req,
		body:                     data,
		method:                   methodCall.MethodName,
		autoCapitalizeMethodName: c.AutoCapitalizeMethodName,
	}
}

type CodecRequest struct {
	req                      *http.Request
	body                     []byte
	method                   string
	autoCapitalizeMethodName bool
}

func (c *CodecRequest) Method() (string, error) {
	if !c.autoCapitalizeMethodName {
		return c.method, nil
	}

	toks := strings.Split(c.method, ".")
	toks[len(toks)-1] = strings.Title(toks[len(toks)-1])
	return strings.Join(toks, "."), nil
}

func (c *CodecRequest) ReadRequest(args interface{}) error {
	mc := MethodCall{}
	if err := xml.Unmarshal(c.body, &mc); err != nil {
		log.WithError(err).Error("xmlrpc: error unmarshalling request body")
		return err
	}

	numArgs := reflect.TypeOf(args).Elem().NumField()
	if numArgs != len(mc.Params) {
		log.Errorf("xmlrpc: wrong number of arguments (expected %d, got %d)", numArgs, len(mc.Params))
		return fmt.Errorf("wrong number of arguments")
	}

	for i, param := range mc.Params {
		field := reflect.ValueOf(args).Elem().Field(i)
		if err := mapValueToField(param.Value, &field); err != nil {
			log.WithError(err).Error("error mapping incoming param to field")
			return err
		}
	}

	return nil
}

func (c *CodecRequest) WriteResponse(w http.ResponseWriter, reply interface{}) {

	data, err := marshalReply(reply)
	if err != nil {
		c.WriteError(w, 500, err)
		return
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write(data)
}

func (c *CodecRequest) WriteError(w http.ResponseWriter, status int, err error) {
	log.WithField("code", status).Errorf("write error: %s", err.Error())
	http.Error(w, http.StatusText(status), status)
}

func marshalReply(reply interface{}) ([]byte, error) {
	replyValue := reflect.ValueOf(reply).Elem()
	if replyValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("xmlrpc reply must be a struct")
	}

	paramsXML, err := marshalReplyParams(&replyValue)
	if err != nil {
		return nil, err
	}

	buf := fmt.Sprintf("<methodResponse><params>%s</params></methodResponse>", paramsXML)

	return []byte(buf), nil
}

func marshalReplyParams(value *reflect.Value) (string, error) {
	buf := ""
	numFields := value.NumField()
	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		paramXML, err := marshalReplyParam(&field)
		if err != nil {
			return "", err
		}
		buf += "<param><value>" + paramXML + "</value></param>"
	}
	return buf, nil
}

func marshalReplyParam(value *reflect.Value) (string, error) {
	log.Infof("marshalling reply value of kind %v", value.Kind())
	switch value.Kind() {
	case reflect.String:
		// FIXME: This needs to be escaped.
		return fmt.Sprintf("<string>%s</string>", value.String()), nil
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return fmt.Sprintf("<int>%d</int>", value.Int()), nil
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return fmt.Sprintf("<int>%d</int>", value.Uint()), nil
	case reflect.Bool:
		return fmt.Sprintf("<boolean>%t</boolean>", value.Bool()), nil
	case reflect.Slice, reflect.Array:
		buf := ""
		length := value.Len()
		for i := 0; i < length; i++ {
			item := value.Index(i)
			itemXML, err := marshalReplyParam(&item)
			if err != nil {
				return "", err
			}

			buf += itemXML
		}
		return fmt.Sprintf("<array><data>%s</data></array>", buf), nil
	case reflect.Struct:
		if reflect.TypeOf(value).String() == "time.Time" {
			t := value.Interface().(time.Time)
			return fmt.Sprintf(
				"<dateTime.iso8601>%04d%02d%02dT%02d:%02d:%02d</dateTime.iso8601>",
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second(),
			), nil
		}

		buf := ""
		numFields := value.NumField()
		for i := 0; i < numFields; i++ {
			fieldType := value.Type().Field(i)
			name := fieldType.Name
			xmlTag := fieldType.Tag.Get("xml")
			if xmlTag != "" {
				name = xmlTag
			}

			field := value.Field(i)
			memberXML, err := marshalReplyParam(&field)
			if err != nil {
				return "", err
			}

			buf += fmt.Sprintf(
				"<member><name>%s</name><value>%s</value></member>",
				name,
				memberXML,
			)
		}
		return fmt.Sprintf("<struct>%s</struct>", buf), nil
	case reflect.Ptr:
		return fmt.Sprintf("<nil/>"), nil
	default:
		return "", fmt.Errorf("unknown reply value type '%v'", value.Kind())
	}
	return "", nil
}

func mapValueToField(value interface{}, field *reflect.Value) error {
	valueKind := reflect.TypeOf(value).Kind()
	fieldKind := field.Kind()

	if valueKind != fieldKind {
		return fmt.Errorf("value type mismatch: (%v; %v)", valueKind, fieldKind)
	}

	if valueKind == reflect.Struct {
		xs := value.(XMLRPCStruct)
		fieldType := field.Type()

		for i := 0; i < len(xs.Members); i++ {
			member := xs.Members[i]
			for j := 0; j < fieldType.NumField(); j++ {
				structField := fieldType.Field(j)
				xmlTag := structField.Tag.Get("xml")
				log.Infof("checking struct field '%s' (tag: '%s') against '%s'...", structField.Name, xmlTag, member.Name)
				if xmlTag == member.Name || strings.Title(member.Name) == structField.Name {
					log.Info("yay struct field name match!")
					targetField := field.FieldByName(structField.Name)
					if err := mapValueToField(member.Value.Value, &targetField); err != nil {
						return err
					}
					break
				}
			}
		}
	} else if valueKind == reflect.Slice {
		xs := value.([]XMLRPCValue)
		slice := reflect.MakeSlice(reflect.TypeOf(field.Interface()), len(xs), len(xs))
		for i, xv := range xs {
			item := slice.Index(i)
			if err := mapValueToField(xv.Value, &item); err != nil {
				log.WithError(err).Error("error mapping element of slice")
				return err
			}
		}
		field.Set(reflect.AppendSlice(*field, slice))
	} else {
		field.Set(reflect.ValueOf(value))
	}

	return nil
}
