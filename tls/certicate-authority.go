package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

type CertificateAuthority struct {
	Certificate []byte
	PrivateKey  *rsa.PrivateKey
	CertPEM     []byte
}

func CreateCertificateAuthority() (*CertificateAuthority, error) {
	keys, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("unable to genarate private keys, error: %s", err)
	}

	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, err
	}

	now := time.Now()

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:    "super.soft",
			Organization:  []string{"SuperSoft"},
			Country:       []string{"US"},
			Province:      []string{"Colorado"},
			Locality:      []string{"Boulder"},
			StreetAddress: []string{"123 Main Lane"},
			PostalCode:    []string{"80301"},
		},
		DNSNames:              []string{"super.soft"},
		NotBefore:             now.Add(-10 * time.Minute).UTC(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	certificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &keys.PublicKey, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to generate certificate, error: %s", err)
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate,
	})

	return &CertificateAuthority{
		Certificate: certificate,
		PrivateKey:  keys,
		CertPEM:     certPEM.Bytes(),
	}, nil
}
