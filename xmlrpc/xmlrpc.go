package xmlrpc

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

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
	data, err := xml.MarshalIndent(reply, "  ", "  ")
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
				if xmlTag == member.Name || strings.Title(member.Name) == structField.Name {
					targetField := field.FieldByName(structField.Name)
					return mapValueToField(member.Value.Value, &targetField)
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
