package foam

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
)

// Doer is the interface used to perform HTTP request.
// The stdlib http.Client implements this interface.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Options contains the options that can be set on the client. Options should
// only be modified by the provided setter functions.
type Options struct {
	binarySecurityToken string
	privateKey          []byte
	client              Doer
}

// Option is a setter for a client option
type Option func(*Options) error

// WithHTTPClient sets the client that will be used to send the HTTP requests.
func WithHTTPClient(client Doer) Option {
	return func(opt *Options) error {
		opt.client = client
		return nil
	}
}

// WithBinarySecurityToken adds a binary security token to every otugoing
// requests. The requests will also be signed with the provided private key.
func WithBinarySecurityToken(cert, key []byte) Option {
	return func(opt *Options) error {
		opt.binarySecurityToken = base64.StdEncoding.EncodeToString(cert)
		opt.privateKey = key
		return nil
	}
}

// A Client is a SOAP client. The zero value is not useful,
// you should instead call NewClient to get an initialized client.
type Client struct {
	endpointUrl string
	client      Doer
	bst         string // Binary Security Token
	privateKey  []byte
}

// NewClients creates a new SOAP client with the provided options
func NewClient(wsdlUrl string, setters ...Option) (*Client, error) {
	endpointUrl, err := url.Parse(wsdlUrl)
	if err != nil {
		return nil, fmt.Errorf("parse wsdl url: %v", err)
	}
	endpointUrl.RawQuery = "" // Remove the ?wdsl query

	opts := Options{
		client: &http.Client{},
	}

	for _, setter := range setters {
		if err := setter(&opts); err != nil {
			setterName := runtime.FuncForPC(reflect.ValueOf(setter).Pointer()).Name()
			return nil, fmt.Errorf("set option: %s: %v", setterName, err)
		}
	}

	return &Client{
		client:      opts.client,
		bst:         opts.binarySecurityToken,
		privateKey:  opts.privateKey,
		endpointUrl: endpointUrl.String(),
	}, nil
}

// Call performs a SOAP 1.1 request to the specified endpoint.
//
// The payload will be seialized to XML and then included into a SOAP 1.1
// evelope, containing the BinarySecurityToken WSSE header.
// The generated XML body is then signled with xmlsec.
//
// The response body will be unmarshalled into the response interface.
func (c *Client) Call(ctx context.Context, endpoint string, payload, response interface{}) error {
	envelope := newEnvelope()
	envelope.Body = newBody(payload)

	if c.privateKey != nil {
		envelope.Header.Security = newSecurityHeader(c.bst, envelope.Body.WSUID)
	}

	signedBody, err := sign(c.privateKey, &envelope)
	if err != nil {
		return fmt.Errorf("sign document: %v", err)
	}

	bodyReader := bytes.NewReader(signedBody)
	req, err := http.NewRequest(http.MethodPost, c.endpointUrl, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %v", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("SOAPAction", endpoint)
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("call %s: %v", endpoint, err)
	}
	defer res.Body.Close()

	if err := xml.NewDecoder(res.Body).Decode(response); err != nil {
		return fmt.Errorf("decode response: %v", err)
	}

	return nil
}
