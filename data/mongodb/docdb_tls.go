package mongodb

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

const (
	// Path to the AWS CA file
	caFilePath = "/var/openfaas/secrets/rds-combined-ca-bundle.pem"

	// Timeout operations after N seconds
	connectTimeout = 5
	queryTimeout   = 30
)

func GetCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}
