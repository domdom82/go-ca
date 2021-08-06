package main

import (
	"flag"
	"fmt"
	"github.com/domdom82/go-ca/ca"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	roots := flag.Int("roots", 1, "Number of root CAs to generate")
	chainsPerRoot := flag.Int("chains", 1, "Number of chains (intermediate CAs) per root CA")
	chainLength := flag.Int("chainlen", 1, "Path length per chain")
	clientCertsPerChain := flag.Int("clients", 1, "Number of valid client certificates per chain")
	serverCertsPerChain := flag.Int("servers", 1, "Number of valid server certificates per chain")
	expiredClientCertsPerChain := flag.Int("expiredClients", 0, "Number of expired client certificates per chain")
	expiredServerCertsPerChain := flag.Int("expiredServers", 0, "Number of expired server certificates per chain")

	flag.Parse()

	//TODO: validate input

	for r := 1; r <= *roots; r++ {
		rootCa := ca.New(fmt.Sprintf("%d", r))
		ca.DumpCertPEM(rootCa.Cert.Subject.CommonName, ca.STATUS_VALID, rootCa.CertPEM, rootCa.KeyPEM)
		for c := 1; c <= *chainsPerRoot; c++ {
			intermediateCa := rootCa.IntermediateCA(fmt.Sprintf("%d-%d", c, 1))
			ca.DumpCertPEM(intermediateCa.Cert.Subject.CommonName, ca.STATUS_VALID, intermediateCa.CertPEM, intermediateCa.KeyPEM)
			for i := 2; i <= *chainLength; i++ {
				intermediateCa = intermediateCa.IntermediateCA(fmt.Sprintf("%d-%d", c, i))
				ca.DumpCertPEM(intermediateCa.Cert.Subject.CommonName, ca.STATUS_VALID, intermediateCa.CertPEM, intermediateCa.KeyPEM)
			}
			for cc := 1; cc <= *clientCertsPerChain; cc++ {
				clientCert := intermediateCa.ClientCert()
				ca.DumpCertPEM(clientCert.Cert.Subject.CommonName, ca.STATUS_VALID, clientCert.CertPEM, clientCert.KeyPEM)
			}
			for cc := 1; cc <= *expiredClientCertsPerChain; cc++ {
				clientCert := intermediateCa.ExpiredClientCert()
				ca.DumpCertPEM(clientCert.Cert.Subject.CommonName, ca.STATUS_EXPIRED, clientCert.CertPEM, clientCert.KeyPEM)
			}
			for sc := 1; sc <= *serverCertsPerChain; sc++ {
				serverCert := intermediateCa.ServerCert()
				ca.DumpCertPEM(serverCert.Cert.Subject.CommonName, ca.STATUS_VALID, serverCert.CertPEM, serverCert.KeyPEM)
			}
			for sc := 1; sc <= *expiredServerCertsPerChain; sc++ {
				serverCert := intermediateCa.ExpiredServerCert()
				ca.DumpCertPEM(serverCert.Cert.Subject.CommonName, ca.STATUS_EXPIRED, serverCert.CertPEM, serverCert.KeyPEM)
			}
		}
	}

}
