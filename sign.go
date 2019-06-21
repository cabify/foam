package foam

import (
	"encoding/xml"
	"fmt"

	"github.com/crewjam/go-xmlsec"
)

// Sign uses xmlsec to sign the document.
// The returned payload must not be modified
func sign(key []byte, envelope interface{}) ([]byte, error) {
	buf, err := xml.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("encode document: %v", err)
	}

	signed, err := xmlsec.Sign(key, buf, xmlsec.SignatureOptions{})
	if err != nil {
		return nil, fmt.Errorf("sign with xmlsec: %v", err)
	}

	return signed, nil
}
