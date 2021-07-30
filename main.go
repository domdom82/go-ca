package main

import (
	"github.com/domdom82/go-ca/ca"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	rootCa := ca.New()

	ca.DumpCertPEM(rootCa.Cert.Subject.CommonName, rootCa.CertPEM, rootCa.KeyPEM)

	intermediateCa := rootCa.IntermediateCA("A")

	ca.DumpCertPEM(intermediateCa.Cert.Subject.CommonName, intermediateCa.CertPEM, intermediateCa.KeyPEM)

	clientCert := intermediateCa.ClientCert()

	ca.DumpCertPEM(clientCert.Cert.Subject.CommonName, clientCert.CertPEM, clientCert.KeyPEM)

	serverCert := intermediateCa.ServerCert()

	ca.DumpCertPEM(serverCert.Cert.Subject.CommonName, serverCert.CertPEM, serverCert.KeyPEM)

}
