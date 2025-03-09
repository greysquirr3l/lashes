package lashes

import (
	"crypto/tls"
)

// TLSMinVersion defines the minimum TLS version to use
type TLSMinVersion int

const (
	// TLSv12 requires TLS 1.2 or higher
	TLSv12 TLSMinVersion = iota
	// TLSv13 requires TLS 1.3 or higher
	TLSv13
)

// SecurityOptions configures security-related settings
type SecurityOptions struct {
	// VerifyTLS determines whether SSL/TLS certificates are verified
	VerifyTLS bool

	// MinTLSVersion sets the minimum TLS version to accept
	MinTLSVersion TLSMinVersion

	// PreferServerCipherSuites lets the server choose the cipher suite
	// Deprecated: This field is ignored in recent Go versions
	PreferServerCipherSuites bool

	// AllowInsecure enables insecure connections (not recommended)
	// This bypasses all TLS security checks
	AllowInsecure bool
}

// DefaultSecurityOptions returns secure default options
func DefaultSecurityOptions() SecurityOptions {
	return SecurityOptions{
		VerifyTLS:                true,
		MinTLSVersion:            TLSv12,
		PreferServerCipherSuites: true,
		AllowInsecure:            false,
	}
}

// GetTLSConfig returns a TLS configuration based on the security options
func GetTLSConfig(opts SecurityOptions) *tls.Config {
	var minVersion uint16

	// Handle all possible TLSMinVersion values
	switch opts.MinTLSVersion {
	case TLSv13:
		minVersion = tls.VersionTLS13
	case TLSv12:
		fallthrough
	default:
		minVersion = tls.VersionTLS12
	}

	config := &tls.Config{
		MinVersion: minVersion,
		// #nosec G402 -- InsecureSkipVerify is explicitly controlled by config options
		// This is intentional to allow users to bypass certificate validation if needed
		InsecureSkipVerify: !opts.VerifyTLS || opts.AllowInsecure,
	}

	return config
}
