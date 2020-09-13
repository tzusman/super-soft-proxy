package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type TLSCertificate struct {
	PrivateKey  rsa.PrivateKey
	Certificate tls.Certificate
	Template    x509.Certificate
	CertPEM     []byte
	KeyPEM      []byte
}

func CreateTLSCertificate(hostname string, caCert CertificateAuthority) (*TLSCertificate, error) {
	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:    hostname,
			Organization:  []string{"SuperSoft"},
			Country:       []string{"US"},
			Province:      []string{"Colorado"},
			Locality:      []string{"Boulder"},
			StreetAddress: []string{"123 Main Lane"},
			PostalCode:    []string{"80301"},
		},
		DNSNames:     []string{hostname},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	caCertCert, err := x509.ParseCertificate(caCert.Certificate)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCertCert, &certPrivKey.PublicKey, caCert.PrivateKey)
	if err != nil {
		return nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	tlsCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return nil, err
	}

	tls := TLSCertificate{
		PrivateKey:  *certPrivKey,
		Certificate: tlsCert,
		Template:    *cert,
		CertPEM:     certPEM.Bytes(),
		KeyPEM:      certPrivKeyPEM.Bytes(),
	}

	return &tls, nil
}
