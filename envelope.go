package foam

import (
	"encoding/xml"

	"github.com/rs/xid"
)

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	SOAPNS  string   `xml:"xmlns:soapenv,attr"`
	Header  *SOAPHeader
	Body    *SOAPBody
}

type SOAPBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	WSUNS   string   `xml:"xmlns:wsu,attr"`
	WSUID   string   `xml:"wsu:Id,attr"`
	XMLID   string   `xml:"xml:id,attr"`
	Payload interface{}
}

func newBody(payload interface{}) *SOAPBody {
	elementID := "Body-" + xid.New().String()

	return &SOAPBody{
		WSUNS:   "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
		WSUID:   elementID,
		XMLID:   elementID,
		Payload: payload,
	}
}

func newEnvelope() *SOAPEnvelope {
	return &SOAPEnvelope{
		SOAPNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Header: &SOAPHeader{},
	}
}
