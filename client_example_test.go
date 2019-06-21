package foam_test

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/cabify/foam"
)

type foo struct {
	XMLName xml.Name `xml:"foo"`
	ID      string   `xml:"id,attr"`
}

type baz struct {
	XMLName xml.Name `xml:"baz"`
	Value   string   `xml:"value,attr"`
}

func ExampleNewClient() {
	// Read the RSA certificate and private key
	cert, err := ioutil.ReadFile("my_server.crt")
	if err != nil {
		log.Fatalf("read certificate: %v", err)
	}
	key, err := ioutil.ReadFile("my_server.key")
	if err != nil {
		log.Fatalf("read key: %v", err)
	}

	// Create an HTTP client with a timeout
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}

	client, err := foam.NewClient("https://example.com/MyServer?wsdl",
		foam.WithBinarySecurityToken(cert, key),
		foam.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalf("read key: %v", err)
	}

	var res baz
	if err := client.Call(context.Background(), "MyEndpoint", &foo{ID: "1"}, &res); err != nil {
		log.Fatalf("make request: %v", err)
	}
}
