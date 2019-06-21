# Foam: SOAP's best friend 

[![Build Status](https://travis-ci.com/cabify/foam.svg?token=seG66JiMjNYXrKzButB4&branch=master)](https://travis-ci.com/cabify/foam)
[![GoDoc](https://godoc.org/github.com/cabify/foam?status.svg)](https://godoc.org/github.com/cabify/foam)

Foam is a SOAP 1.1 client for Go which implements the [WSS BinarySecurityToken](https://www.oasis-open.org/committees/download.php/16790/wss-v1.1-spec-os-SOAPMessageSecurity.pdf)
and [XML Signature](https://www.w3.org/TR/xmldsig-core1/) standards.

Currently it supports only signing with X.509 certificates using the [`rsa-sha1`](https://www.w3.org/TR/xmldsig-core1/#sec-PKCS1)
algorithm.

#### Warning

This package uses `cgo` and has an hard-dependency on [XMLSec](https://www.aleksey.com/xmlsec/)
and [LibXML2](http://xmlsoft.org/), so don't skimp on the installation instructions.
Due to this dependency, this library is not recommended for simpler use-cases
where XML Digital Signature is not a requirement.

## Usage

```go
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

func main() {
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
```

## Installation

### macOS

The easiest way to install all the required dependencies is with Homebrew:

```shell
brew install libxmlsec1 libxml2 pkg-config
```

**Make sure to follow the post-install instructions** printed by homebrew,
otherwise the Go compiler won't be able to find the libraries on your machine.

### Linux

#### Debian

```
apt-get install -y libxml2-dev libxmlsec1-dev pkg-config
```

## Build

When building the package, you must have CGO enabled and set the `CGO_CFLAGS_ALLOW`
to `-w|-UXMLSEC_CRYPTO_DYNAMIC_LOADING`.
For example:

```
CGO_CFLAGS_ALLOW="-w|-UXMLSEC_CRYPTO_DYNAMIC_LOADING" go build
```

