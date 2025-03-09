package lashes

import (
	"crypto/tls"
	"testing"
)

func TestDefaultSecurityOptions(t *testing.T) {
	opts := DefaultSecurityOptions()

	// Check default values
	if !opts.VerifyTLS {
		t.Error("Default VerifyTLS should be true")
	}

	if opts.MinTLSVersion != TLSv12 {
		t.Errorf("Default MinTLSVersion = %v, want %v", opts.MinTLSVersion, TLSv12)
	}

	if !opts.PreferServerCipherSuites {
		t.Error("Default PreferServerCipherSuites should be true")
	}

	if opts.AllowInsecure {
		t.Error("Default AllowInsecure should be false")
	}
}

func TestGetTLSConfig(t *testing.T) {
	testCases := []struct {
		name                   string
		opts                   SecurityOptions
		wantMinVersion         uint16
		wantInsecureSkipVerify bool
	}{
		{
			name:                   "Default secure options",
			opts:                   DefaultSecurityOptions(),
			wantMinVersion:         tls.VersionTLS12,
			wantInsecureSkipVerify: false,
		},
		{
			name: "TLS 1.3 minimum",
			opts: SecurityOptions{
				VerifyTLS:                true,
				MinTLSVersion:            TLSv13,
				PreferServerCipherSuites: true,
				AllowInsecure:            false,
			},
			wantMinVersion:         tls.VersionTLS13,
			wantInsecureSkipVerify: false,
		},
		{
			name: "Disabled certificate verification",
			opts: SecurityOptions{
				VerifyTLS:                false,
				MinTLSVersion:            TLSv12,
				PreferServerCipherSuites: true,
				AllowInsecure:            false,
			},
			wantMinVersion:         tls.VersionTLS12,
			wantInsecureSkipVerify: true,
		},
		{
			name: "Allow insecure connections",
			opts: SecurityOptions{
				VerifyTLS:                true,
				MinTLSVersion:            TLSv12,
				PreferServerCipherSuites: true,
				AllowInsecure:            true,
			},
			wantMinVersion:         tls.VersionTLS12,
			wantInsecureSkipVerify: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := GetTLSConfig(tc.opts)

			if config.MinVersion != tc.wantMinVersion {
				t.Errorf("MinVersion = %d, want %d", config.MinVersion, tc.wantMinVersion)
			}

			if config.InsecureSkipVerify != tc.wantInsecureSkipVerify {
				t.Errorf("InsecureSkipVerify = %v, want %v",
					config.InsecureSkipVerify, tc.wantInsecureSkipVerify)
			}

			// Remove test for deprecated field PreferServerCipherSuites
		})
	}
}
