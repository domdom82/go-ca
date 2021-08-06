package ca

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

const (
	STATUS_VALID   string = "valid"
	STATUS_EXPIRED string = "expired"
	STATUS_REVOKED string = "revoked"
)

type CA struct {
	company *company
	Cert    *x509.Certificate
	Key     *rsa.PrivateKey
	CertPEM *bytes.Buffer
	KeyPEM  *bytes.Buffer
}

type ClientCert struct {
	Cert    *x509.Certificate
	Key     *rsa.PrivateKey
	CertPEM *bytes.Buffer
	KeyPEM  *bytes.Buffer
}

type ServerCert struct {
	Cert    *x509.Certificate
	Key     *rsa.PrivateKey
	CertPEM *bytes.Buffer
	KeyPEM  *bytes.Buffer
}

func DumpCertPEM(prefix string, status string, certPEM *bytes.Buffer, keyPEM *bytes.Buffer) {
	fmt.Println(prefix + ".crt")
	fmt.Println(status)
	fmt.Println(certPEM)
	fmt.Println(prefix + ".key")
	fmt.Println(status)
	fmt.Println(keyPEM)
}

func certPEM(certBytes []byte) *bytes.Buffer {
	certPEM := new(bytes.Buffer)
	err := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		panic(err)
	}
	return certPEM
}

func keyPEM(key *rsa.PrivateKey) *bytes.Buffer {
	keyPEM := new(bytes.Buffer)
	err := pem.Encode(keyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		panic(err)
	}
	return keyPEM
}

func New(suffix string) *CA {

	company := randomCompany()

	ca := &CA{}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{company.name},
			Country:       []string{company.address.country},
			Province:      []string{""},
			Locality:      []string{company.address.city},
			StreetAddress: []string{fmt.Sprintf("%d %s", company.address.number, company.address.street)},
			PostalCode:    []string{fmt.Sprintf("%d", company.address.zip)},
			CommonName:    fmt.Sprintf("%s RootCA-%s", company.name, suffix),
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		panic(err)
	}

	caPEM := certPEM(caBytes)
	caPrivKeyPEM := keyPEM(caPrivKey)

	ca.Cert = cert
	ca.Key = caPrivKey
	ca.CertPEM = caPEM
	ca.KeyPEM = caPrivKeyPEM
	ca.company = company

	return ca
}

func (ca *CA) IntermediateCA(suffix string) *CA {
	intermediateCA := &CA{}
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  ca.Cert.Subject.Organization,
			Country:       ca.Cert.Subject.Country,
			Province:      ca.Cert.Subject.Province,
			Locality:      ca.Cert.Subject.Locality,
			StreetAddress: ca.Cert.Subject.StreetAddress,
			PostalCode:    ca.Cert.Subject.PostalCode,
			CommonName:    fmt.Sprintf("%s-%s", "IntermediateCA", suffix),
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &caPrivKey.PublicKey, ca.Key)
	if err != nil {
		panic(err)
	}

	caPEM := certPEM(caBytes)
	caPrivKeyPEM := keyPEM(caPrivKey)

	intermediateCA.Cert = cert
	intermediateCA.Key = caPrivKey
	intermediateCA.CertPEM = caPEM
	intermediateCA.KeyPEM = caPrivKeyPEM
	intermediateCA.company = ca.company

	return intermediateCA
}

func (ca *CA) ClientCert() *ClientCert {
	return ca.doClientCert(false, false)
}

func (ca *CA) ExpiredClientCert() *ClientCert {
	return ca.doClientCert(true, false)
}

func (ca *CA) RevokedClientCert() *ClientCert {
	return ca.doClientCert(false, true)
}

func (ca *CA) doClientCert(expired bool, revoked bool) *ClientCert {
	p := randomPerson(ca.company)
	notBefore := time.Now()
	notAfter := time.Now().AddDate(10, 0, 0)
	if expired {
		notBefore = time.Now().AddDate(-2, 0, 0)
		notAfter = time.Now().AddDate(-1, 0, 0)
	}
	clientCert := &ClientCert{}
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  ca.Cert.Subject.Organization,
			Country:       ca.Cert.Subject.Country,
			Province:      ca.Cert.Subject.Province,
			Locality:      ca.Cert.Subject.Locality,
			StreetAddress: ca.Cert.Subject.StreetAddress,
			PostalCode:    ca.Cert.Subject.PostalCode,
			CommonName:    fmt.Sprintf("%s %s", p.first, p.last),
		},
		EmailAddresses: []string{p.email},
		NotBefore:      notBefore,
		NotAfter:       notAfter,
		SubjectKeyId:   []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:       x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &certPrivKey.PublicKey, ca.Key)
	if err != nil {
		panic(err)
	}

	certPEM := certPEM(certBytes)
	keyPEM := keyPEM(certPrivKey)

	clientCert.CertPEM = certPEM
	clientCert.KeyPEM = keyPEM
	clientCert.Cert = cert
	clientCert.Key = certPrivKey

	return clientCert
}

func (ca *CA) ServerCert() *ServerCert {
	return ca.doServerCert(false, false)
}

func (ca *CA) ExpiredServerCert() *ServerCert {
	return ca.doServerCert(true, false)
}

func (ca *CA) RevokedServerCert() *ServerCert {
	return ca.doServerCert(false, true)
}

func (ca *CA) doServerCert(expired bool, revoked bool) *ServerCert {
	hostname := fmt.Sprintf("%s.%s", randomHostname(), ca.company.domain)
	notBefore := time.Now()
	notAfter := time.Now().AddDate(10, 0, 0)
	if expired {
		notBefore = time.Now().AddDate(-2, 0, 0)
		notAfter = time.Now().AddDate(-1, 0, 0)
	}
	serverCert := &ServerCert{}
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  ca.Cert.Subject.Organization,
			Country:       ca.Cert.Subject.Country,
			Province:      ca.Cert.Subject.Province,
			Locality:      ca.Cert.Subject.Locality,
			StreetAddress: ca.Cert.Subject.StreetAddress,
			PostalCode:    ca.Cert.Subject.PostalCode,
			CommonName:    hostname,
		},
		DNSNames:     []string{hostname},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &certPrivKey.PublicKey, ca.Key)
	if err != nil {
		panic(err)
	}

	certPEM := certPEM(certBytes)
	keyPEM := keyPEM(certPrivKey)

	serverCert.CertPEM = certPEM
	serverCert.KeyPEM = keyPEM
	serverCert.Cert = cert
	serverCert.Key = certPrivKey

	return serverCert
}
