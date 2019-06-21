// Package foam is a SOAP 1.1 client for Go which implements the
// WSS BinarySecurityToken and XML Digital Signature standards.
//
// Due to limitations in Go abilities to handle XML, it uses CGO and
// depends on xmlsec and LibXML2 to sign the generated XML documents.
// For simpler use-cases that don't require signed documents, a different
// library is recommended.
//
// To compile the package, you must allow some CGO flags:
//
//     CGO_CFLAGS_ALLOW="-w|-UXMLSEC_CRYPTO_DYNAMIC_LOADING"
//
package foam
