package foam

import (
	"encoding/xml"

	"github.com/rs/xid"
)

type SOAPHeader struct {
	XMLName  xml.Name `xml:"soapenv:Header"`
	Security *SecurityHeader
}

type BinarySecurityTokenHeader struct {
	XMLName      xml.Name `xml:"wsse:BinarySecurityToken"`
	WSUID        string   `xml:"wsu:Id,attr"`
	EncodingType string   `xml:"EncodingType,attr"`
	ValueType    string   `xml:"ValueType,attr"`
	Token        string   `xml:",innerxml"`
}

type SecurityHeader struct {
	XMLName             xml.Name `xml:"wsse:Security"`
	XMLNSWSSE           string   `xml:"xmlns:wsse,attr"`
	XMLNSWSU            string   `xml:"xmlns:wsu,attr"`
	BinarySecurityToken *BinarySecurityTokenHeader
	Signature           *SignatureHeader
}

type SignatureHeader struct {
	XMLName        xml.Name `xml:"ds:Signature"`
	XMLNSDS        string   `xml:"xmlns:ds,attr"`
	ID             string   `xml:"Id,attr"`
	SignedInfo     SignedInfo
	SignatureValue SignatureValue
	KeyInfo        KeyInfo
}

type SignedInfo struct {
	XMLName                xml.Name `xml:"ds:SignedInfo"`
	CanonicalizationMethod CanonicalizationMethod
	SignatureMethod        SignatureMethod
	Reference              DSReference
}

type CanonicalizationMethod struct {
	XMLName   xml.Name `xml:"ds:CanonicalizationMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type SignatureMethod struct {
	XMLName   xml.Name `xml:"ds:SignatureMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type DSReference struct {
	XMLName      xml.Name `xml:"ds:Reference"`
	URI          string   `xml:"URI,attr"`
	Transforms   Transforms
	DigestMethod DigestMethod
	DigestValue  DigestValue
}

type Transforms struct {
	XMLName   xml.Name `xml:"ds:Transforms"`
	Transform Transform
}

type Transform struct {
	XMLName   xml.Name `xml:"ds:Transform"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type DigestMethod struct {
	XMLName   xml.Name `xml:"ds:DigestMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type DigestValue struct {
	XMLName xml.Name `xml:"ds:DigestValue"`
}

type SignatureValue struct {
	XMLName xml.Name `xml:"ds:SignatureValue"`
}

type KeyInfo struct {
	XMLName                xml.Name `xml:"ds:KeyInfo"`
	SecurityTokenReference SecurityTokenReference
}

type SecurityTokenReference struct {
	XMLName   xml.Name `xml:"wsse:SecurityTokenReference"`
	Reference WSSEReference
}

type WSSEReference struct {
	XMLName   xml.Name `xml:"wsse:Reference"`
	URI       string   `xml:"URI,attr"`
	ValueType string   `xml:"ValueType,attr"`
}

func newSecurityHeader(token string, bodyID string) *SecurityHeader {
	securityTokenID := "SecurityToken-" + xid.New().String()
	signatureID := "Signature-" + xid.New().String()

	return &SecurityHeader{
		XMLNSWSSE: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
		XMLNSWSU:  "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
		BinarySecurityToken: &BinarySecurityTokenHeader{
			WSUID:        securityTokenID,
			EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary",
			ValueType:    "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
			Token:        token,
		},
		Signature: &SignatureHeader{
			XMLNSDS: "http://www.w3.org/2000/09/xmldsig#",
			ID:      signatureID,
			SignedInfo: SignedInfo{
				CanonicalizationMethod: CanonicalizationMethod{Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#"},
				SignatureMethod:        SignatureMethod{Algorithm: "http://www.w3.org/2000/09/xmldsig#rsa-sha1"},
				Reference: DSReference{
					URI: "#" + bodyID,
					Transforms: Transforms{
						Transform: Transform{Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#"},
					},
					DigestMethod: DigestMethod{Algorithm: "http://www.w3.org/2000/09/xmldsig#sha1"},
				},
			},
			KeyInfo: KeyInfo{
				SecurityTokenReference: SecurityTokenReference{
					Reference: WSSEReference{
						URI:       "#" + securityTokenID,
						ValueType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
					},
				},
			},
		},
	}
}
