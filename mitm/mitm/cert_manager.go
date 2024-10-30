package mitm

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/ethanvc/evol/base"
	"math/big"
	"time"
)

type CertManager struct {
}

func NewCertManager() *CertManager {
	return &CertManager{}
}

func (cm *CertManager) CreateRootCert() error {
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return base.ErrWithCaller(err)
	}

	// 生成CA证书模板
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Go CA"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCert, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return base.ErrWithCaller(err)
	}
	_ = caCert
	return nil
}
