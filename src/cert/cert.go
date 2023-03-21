package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

func CreateCert(tls_path string) (err error) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Country:            []string{"CN"},
		Province:           []string{"BeiJing"},
		Organization:       []string{"SRCHTTP"},
		OrganizationalUnit: []string{"certDevops"},
		CommonName:         "src-http.com",
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("0.0.0.0")},
	}

	// 创建证书目录
	os.MkdirAll(tls_path, os.ModePerm)

	pk, err := rsa.GenerateKey(rand.Reader, 2048)

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	certOut, err := os.Create(tls_path + "/server.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.Create(tls_path + "/server.key")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()

	return
}
