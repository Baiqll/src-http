package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/baiqll/src-http/src/lib"
)

var CACERT = `-----BEGIN CERTIFICATE-----
MIIDgzCCAmugAwIBAgIEZMcmKDANBgkqhkiG9w0BAQsFADBgMRAwDgYDVQQGEwdT
UkNIVFRQMRAwDgYDVQQIEwdTUkNIVFRQMRAwDgYDVQQKEwdTUkNIVFRQMRMwEQYD
VQQLEwpTUkNIVFRQIENBMRMwEQYDVQQDEwpTUkNIVFRQIENBMB4XDTIzMDczMTAz
MTAzMloXDTMzMDczMTAzMTAzMlowYDEQMA4GA1UEBhMHU1JDSFRUUDEQMA4GA1UE
CBMHU1JDSFRUUDEQMA4GA1UEChMHU1JDSFRUUDETMBEGA1UECxMKU1JDSFRUUCBD
QTETMBEGA1UEAxMKU1JDSFRUUCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAM0sGzrsNOuLK6HrQqnE7G2T5nJffxGae3wcSi5KzkaHiNn7zrVNMzpt
qTbYQV5lvA4IXmZjHzYgqxSypzhOLHMS5ZlWbm6C1bN7IjtfGOkUhBIPrr4LYOm6
AKUXGQHPV/+ibeQtAxkWLVISrmGg6fUC1s79RKimwyW0gC61AO8mW3o38dvYxdFK
lOTkGaqEoUnjThDfwboKNSgE8HYK2aNB4JoHl6OkeZKhtT/q45VZ9QbK4TNq80Qa
QbUie9AZ6qsv1fsDp6Sw/ntz+J6b77erCGEoCzEqOj1FxReWXUmR+z2dedBRJCRf
mvthN890DM6X3RGozHlx4W1ZuKWn0FUCAwEAAaNFMEMwDgYDVR0PAQH/BAQDAgEG
MBIGA1UdEwEB/wQIMAYBAf8CAQEwHQYDVR0OBBYEFAKp3G97xnVBF7fuzpjBrXgA
owDcMA0GCSqGSIb3DQEBCwUAA4IBAQC1J8OEYOyqBnn2Qi+BKH/c2PbUfdVxWCi2
IrLz2YLSmMkpqU2VQB0qsz0/49oO5Fzc3thKh4DytUM3kP5iKTs4P7B9NrapCOAl
XijFYk5aAEDDZXEfWIJPB+6MtURnEZdvcsgUhqzWtwxK31j4yeeQgZiJ91B5g5la
Ry8tA8cSYSS+ppjAHGUEkKqkxenBylZ2z/lloQTZ6Y2+0X5lmpupCuXfDgSgmrgK
FShjbpWBzWPXwOCJiQZVc5z4MKTJxtzlGGcbPTtSt0oe60jzmDkXRu2h4uT1JiMY
6rqXSQtHmNbjCX/2sqq/88594lGUTiaNvYlnXQQxTjrzxaDCM2yY
-----END CERTIFICATE-----`

var CAKEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAzSwbOuw064sroetCqcTsbZPmcl9/EZp7fBxKLkrORoeI2fvO
tU0zOm2pNthBXmW8DgheZmMfNiCrFLKnOE4scxLlmVZuboLVs3siO18Y6RSEEg+u
vgtg6boApRcZAc9X/6Jt5C0DGRYtUhKuYaDp9QLWzv1EqKbDJbSALrUA7yZbejfx
29jF0UqU5OQZqoShSeNOEN/Bugo1KATwdgrZo0HgmgeXo6R5kqG1P+rjlVn1Bsrh
M2rzRBpBtSJ70Bnqqy/V+wOnpLD+e3P4npvvt6sIYSgLMSo6PUXFF5ZdSZH7PZ15
0FEkJF+a+2E3z3QMzpfdEajMeXHhbVm4pafQVQIDAQABAoIBAAxPRSL34QTwyKFi
WBGPew/n+7+I8zq/JgGAQQMeAdpBb3iEnxZJl3U99xUPTHy5ZdsBrYg/EjRRKXKI
dXfvWMNest/MS7vdpayrCpa9UeVKEdJzlmxYimv7eOZuyFVPd1wjBqzV9oWeywFN
laDN4ruMfA7XKzNjLfopJjenLHMsom9LVmbOAGYeRZyBNURiLlCK23Tlgk3qWigV
b6dRn7hnGPiKplUEuU3uUvPpmrkG7uH1qWkygBEehazvotBV9TYna8GgRDrBzKcC
BVlY9LZSAbKlULGF6w0Dqb1GiCeo/4nGj/s11MhUuJbXYxC/xslsGwtozs4NyfIB
sdWucdECgYEA2esHWmWeTu9c/L9UGParKCiVu/1igeKg0E57dE5cPDgBAp638YfJ
P8cl5hHjQL3hAhV2qR0Zz8IO0NJSxpTy4fsqe3SXoYjMUave1OSzKnQ5UgVowKjy
J7HXA7WQLQoA8g1XHSK174iThNKdVMGeuPagjRjJSZ3u2RzRZCY8masCgYEA8Qbe
964xqWXKdEmUDeF10HTS6NB/H9irzFdZlV4HOQXDgx22VJ3e+x6ViLgFHuToKEWs
BijQJRhnMJiBEz9E4eBrURTfyw44DNsgtJCcl3myTA/ESEOOfzZbTfPYhUPrMifN
7iQ+O8KAw18yLxcGDsCZhzbdSJiLGxQjCMEHPf8CgYEAjeOqdgGUgnD4atlpOJfj
+dHzLORfL5MQgpGXcLNU+yC8B6iwvGNddlmFI7ih75Wy3Fh9Wr/H/q6sVuubWhHB
08JmdtwDnvojj0oJXTVMM2hZqj47ZraadZ4mEhQ2PB03YGOvRRlEvSKAawt3xagM
YQK0pypsZbKfwl4xOLRs4OECgYBd9z3J9eFql0Kcn2rXFoTl5gWruk01TzV7Drrg
Hq5WLscQQO8qgfnCkSPfD07/wmI4ASGVrSeorqDcMzhvFoV2QhXUoHy3Hy3+5RcV
DiPechVuzd7KBXxyX/CsrVpGajoxbY89Pmf8yFGG2YApF6LG8ZNpQZx3hvEEd49J
BGgcZQKBgQCz9aU8ERchZAfcHsHNC+5gCZYVxmrTddClFoFRvwUfH2jsyO3Kh8m3
Vy+NzMRyvpdHkNx7KN7pCS0SqrRTaHwbdYuX4EVD2bTdKepXffyPXGQBjDt6YLfu
89pk9OcQQNukaCXWSR8LqTJ6K7wsEL9huXJBcrYeitwwldVxI+KDqw==
-----END RSA PRIVATE KEY-----`


type Cert struct {
	Cert       *x509.Certificate
	CertKey    *rsa.PrivateKey
	CertPem    *bytes.Buffer
	CertKeyPem *bytes.Buffer
	CSR        *x509.Certificate
}

func LoadCACert() (ca Cert) {

	CertBlock, _ := pem.Decode([]byte(CACERT))
	KeyBlock, _ := pem.Decode([]byte(CAKEY))

	ca.CertKey, _ = x509.ParsePKCS1PrivateKey(KeyBlock.Bytes)
	ca.Cert, _ = x509.ParseCertificate(CertBlock.Bytes)

	return

}

func CreateCert(tls_path string, domain_list []string, internet_ip string) (err error) {
	/*
		利用CA 证书颁发https 证书
	*/

	ca := LoadCACert()

	template := x509.Certificate{
		Version:      3,
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Country:            []string{"SRCHTTP"},
			Province:           []string{"SRCHTTP"},
			Locality:           []string{"SRCHTTP"},
			Organization:       []string{"SRCHTTP"},
			OrganizationalUnit: []string{"SRCHTTP"},
			CommonName:         "srchttp.com",
		},
		NotBefore:             time.Now().AddDate(0, 0, -2),
		NotAfter:              time.Now().AddDate(2, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP(internet_ip)},
		DNSNames: domain_list,
	}

	// 创建证书目录
	os.MkdirAll(tls_path, os.ModePerm)

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, ca.Cert, &pk.PublicKey, ca.CertKey)
	certOut, _ := os.Create(filepath.Join(tls_path, "server.pem"))
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, _ := os.Create(filepath.Join(tls_path, "server.key"))
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()

	return
}

func CreateTlsCert(tls_path string,  domain_list []string, internet_ip string,is_new bool) (err error) {

	// tls 证书路径
	if !lib.Exists(filepath.Join(tls_path, "server.pem")) || !lib.Exists(filepath.Join(tls_path, "server.key"))|| is_new {

		// 创建证书
		err = CreateCert(tls_path, domain_list, internet_ip)

	}
	return
}

