package network

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/idna"
)

const (
	// SchemeHTTP defines the HTTP scheme.
	SchemeHTTP = "http"
	// SchemeHTTPS defines the HTTPS scheme.
	SchemeHTTPS = "https"
	// SchemePostgres defines the Postgres scheme.
	SchemePostgres = "postgres"
)

var (
	// DomainRegexp can matches a valid domain name.
	DomainRegexp = regexp.MustCompile(`^([a-zA-Z0-9-_]{1,63}\.)*([a-zA-Z0-9-]{1,63})$`)
	// IPRegexp matches a valid IPv4 address.
	IPRegexp = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
)

// ValidateHostname checks the given host string and return an error if it is not a valid hostname.
func ValidateHostname(host string) error {
	if host == "" {
		return errors.New("empty host")
	}
	host = strings.ToLower(host)
	if DomainRegexp.MatchString(host) {
		return nil
	}
	if punycode, err := idna.ToASCII(host); err != nil {
		return err
	} else if DomainRegexp.MatchString(punycode) {
		return nil
	}
	if IPRegexp.MatchString(host) {
		return nil
	}
	return errors.New("invalid host")
}
