package mitm

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/ethanvc/evol/base"
	"math"
	"math/big"
	"os"
	"time"
)

type CertManager struct {
	rootKey  *rsa.PrivateKey
	rootCert *x509.Certificate
}

func NewCertManager() (*CertManager, error) {
	mgr := &CertManager{}
	err := mgr.CreateRootCert()
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

func (mgr *CertManager) GetCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return mgr.CreateDomainCert(info.ServerName)
}

func (mgr *CertManager) generateSerialNumber() *big.Int {
	n, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return n
}

func (mgr *CertManager) CreateRootCert() error {
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return base.ErrWithCaller(err)
	}
	mgr.rootKey = caKey

	// 生成CA证书模板
	caTemplate := x509.Certificate{
		SerialNumber: mgr.generateSerialNumber(),
		Subject: pkix.Name{
			Organization: []string{"Mitm CA"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCert, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return base.ErrWithCaller(err)
	}
	mgr.rootCert = &caTemplate
	content := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCert})
	err = os.WriteFile("root.crt.pem", content, 0600)
	if err != nil {
		return base.ErrWithCaller(err)
	}
	content = x509.MarshalPKCS1PrivateKey(caKey)
	content = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: content})
	err = os.WriteFile("root.key.pem", content, 0600)
	if err != nil {
		return base.ErrWithCaller(err)
	}
	return nil
}

func (mgr *CertManager) CreateDomainCert(domain string) (*tls.Certificate, error) {
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, base.ErrWithCaller(err)
	}

	caTemplate := x509.Certificate{
		SerialNumber: mgr.generateSerialNumber(),
		Subject: pkix.Name{
			Organization: []string{"Domain Server"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * 365 * 5),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{domain},
		SubjectKeyId: []byte{1, 2, 3, 4},
	}

	caCert, err := x509.CreateCertificate(rand.Reader, &caTemplate, mgr.rootCert, &caKey.PublicKey, mgr.rootKey)
	if err != nil {
		return nil, base.ErrWithCaller(err)
	}
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCert})
	content := x509.MarshalPKCS1PrivateKey(caKey)
	pekKey := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: content})
	certPair, err := tls.X509KeyPair(pemCert, pekKey)
	if err != nil {
		return nil, err
	}
	return &certPair, nil
}
