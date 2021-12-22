package network

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ParseEndpoint turns the given URL into a url.URL. An optional scheme can be provided to be
// prepended to rawURL without a scheme prefix.
func ParseEndpoint(rawURL string, scheme string) (*url.URL, error) {
	if len(scheme) == 0 {
		scheme = SchemeHTTP
	}
	if !strings.Contains(rawURL, "://") {
		rawURL = scheme + "://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	parsedURL.Host = strings.ToLower(parsedURL.Host)
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
	host, _, err := SplitHostPort(parsedURL.Host)
	if err != nil {
		return nil, err
	}
	if err = ValidateHostname(host); err != nil {
		return nil, err
	}
	return parsedURL, nil
}

// SplitHostPort takes a valid network address in the format of "host:port" and return the host
// and the port separately.
func SplitHostPort(host string) (string, *uint16, error) {
	i := strings.LastIndex(host, ":")
	switch strings.LastIndex(host, ":") {
	case -1:
		return host, nil, nil
	case len(host) - 1:
		return "", nil, errors.New("empty port")
	}
	port, err := strconv.ParseUint(host[i+1:], 10, 16)
	if err != nil {
		return "", nil, errors.Wrap(err, "invalid port")
	}
	host = host[:i]
	port16 := uint16(port)
	return host, &port16, nil
}

// ParseBasicAuth parse and validates the given HTTP Basic Auth information and return their
// explicit values.
func ParseBasicAuth(username, password, passwordFile string) (string, string, error) {
	usernameSpecified := len(username) > 0
	passwordSpecified := len(password) > 0
	passwordFileSpecified := len(passwordFile) > 0
	if passwordSpecified && passwordFileSpecified {
		return "", "", errors.New("only one of password and password file can be used at the same time")
	}
	if !passwordFileSpecified && !passwordSpecified {
		if !usernameSpecified {
			return "", "", nil
		}
		return "", "", errors.New("empty password")
	}
	if !usernameSpecified {
		return "", "", errors.New("empty username")
	}
	if passwordFileSpecified {
		passwd, err := ioutil.ReadFile(passwordFile)
		if err != nil {
			return "", "", errors.Wrapf(err, "read password file %s", passwordFile)
		}
		return username, string(passwd), nil
	}
	return username, password, nil
}

// TLSClientConfig configures a TLS client.
type TLSClientConfig struct {
	CAFile             string
	CertFile           string
	KeyFile            string
	InsecureSkipVerify bool
	ServerName         string
}

// NewTLSClientConfig builds a tls.Config from a TLSClientConfig.
func NewTLSClientConfig(config TLSClientConfig) (*tls.Config, error) {
	ret := &tls.Config{
		InsecureSkipVerify: config.InsecureSkipVerify,
		ServerName:         config.ServerName,
	}
	if len(config.CAFile) > 0 {
		ca, err := ioutil.ReadFile(config.CAFile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read CA cert %s", config.CAFile)
		}
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(ca) {
			return nil, errors.Errorf("failed to parse CA cert %s", config.CAFile)
		}
		ret.RootCAs = caPool
	}
	certExists := len(config.CertFile) > 0
	keyExists := len(config.KeyFile) > 0
	if !certExists && !keyExists {
		return ret, nil
	}
	if !certExists {
		return nil, errors.Errorf("client cert file %q specified without client key file", config.CertFile)
	}
	if !keyExists {
		return nil, errors.Errorf("client key file %q specified without client cert file", config.KeyFile)
	}
	ret.GetClientCertificate = func(_ *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, err
		}
		return &cert, nil
	}
	return ret, nil
}

// HTTPOption provide a supplement way to configure the RoundTripper built from NewRoundTripper.
type HTTPOption func(*http.Transport)

// NewRoundTripper builds a http.RoundTripper with the default configs; the default configs can be
// overwritten with the given options.
func NewRoundTripper(options ...HTTPOption) http.RoundTripper {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          10000,
		MaxIdleConnsPerHost:   1000,
		IdleConnTimeout:       5 * time.Minute,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	for _, option := range options {
		option(transport)
	}
	return transport
}
