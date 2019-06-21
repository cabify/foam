package foam

import (
	"bytes"
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/crewjam/go-xmlsec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var rsaCertPEM = []byte(`-----BEGIN CERTIFICATE-----
MIIB0zCCAX2gAwIBAgIJAI/M7BYjwB+uMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTIwOTEyMjE1MjAyWhcNMTUwOTEyMjE1MjAyWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANLJ
hPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wok/4xIA+ui35/MmNa
rtNuC+BdZ1tMuVCPFZcCAwEAAaNQME4wHQYDVR0OBBYEFJvKs8RfJaXTH08W+SGv
zQyKn0H8MB8GA1UdIwQYMBaAFJvKs8RfJaXTH08W+SGvzQyKn0H8MAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQEFBQADQQBJlffJHybjDGxRMqaRmDhX0+6v02TUKZsW
r5QuVbpQhH6u+0UgcW0jp9QwpxoPTLTWGXEWBBBurxFwiCBhkQ+V
-----END CERTIFICATE-----
`)

var rsaKeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBANLJhPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wo
k/4xIA+ui35/MmNartNuC+BdZ1tMuVCPFZcCAwEAAQJAEJ2N+zsR0Xn8/Q6twa4G
6OB1M1WO+k+ztnX/1SvNeWu8D6GImtupLTYgjZcHufykj09jiHmjHx8u8ZZB/o1N
MQIhAPW+eyZo7ay3lMz1V01WVjNKK9QSn1MJlb06h/LuYv9FAiEA25WPedKgVyCW
SmUwbPw8fnTcpqDWE3yTO3vKcebqMSsCIBF3UmVue8YU3jybC3NxuXq3wNm34R8T
xVLHwDXh/6NJAiEAl2oHGGLz64BuAfjKrqwz7qMYr9HCLIe/YsoWq/olzScCIQDi
D2lWusoe2/nEqfDVVWGWlyJ7yOmqaVm/iNUN9B2N2g==
-----END RSA PRIVATE KEY-----
`)

type mockClient struct {
	mock.Mock
	T *testing.T
}

// Do provides a mock client which verifies that the request contains
// an XML with a valid signature.
func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	// Check that the request XML contains a valid signature
	buf, err := ioutil.ReadAll(req.Body)
	require.NoErrorf(m.T, err, "read request body: %v", err)

	err = xmlsec.Verify(rsaCertPEM, buf, xmlsec.SignatureOptions{})
	require.NoErrorf(m.T, err, "verify xml: %v", err)

	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type testRequest struct {
	XMLName xml.Name `xml:"request"`
	Bar     string   `xml:"bar,attr"`
}

type testResponse struct {
	XMLName xml.Name `xml:"test"`
	Foo     string   `xml:"foo,attr"`
}

func TestClient_Call(t *testing.T) {
	doer := mockClient{T: t}

	doer.On("Do", mock.Anything).Return(&http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(`<test foo="foo" />`))),
	}, nil)

	client, err := NewClient(
		"https://example.com/TestServer?wsdl",
		WithBinarySecurityToken(rsaCertPEM, rsaKeyPEM),
		WithHTTPClient(&doer),
	)
	require.NoError(t, err)

	req := testRequest{Bar: "bar"}

	var res testResponse
	err = client.Call(context.Background(), "DoSomething", &req, &res)
	require.NoError(t, err)

	doer.AssertExpectations(t)
	assert.Equal(t, "foo", res.Foo)
}
