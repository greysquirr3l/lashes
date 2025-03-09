package agent

import (
	"crypto/rand"
	"math/big"
)

// Common browser versions
var (
	chromeVersions  = []string{"120.0.0", "121.0.0", "122.0.0"}
	firefoxVersions = []string{"122.0", "123.0", "124.0"}
	safariVersions  = []string{"17.2", "17.1", "16.6"}
	osVersions      = []string{
		"Windows NT 10.0; Win64; x64",
		"Macintosh; Intel Mac OS X 10_15_7",
		"X11; Linux x86_64",
	}
)

// GetRandomUserAgent returns a randomly constructed user agent using crypto/rand.
func GetRandomUserAgent() string {
	// Pick a random browser type securely
	browserType, err := secureRandInt(3)
	if err != nil {
		browserType = 0
	}
	os := secureRandomChoice(osVersions)

	switch browserType {
	case 0: // Chrome
		chromeVer := secureRandomChoice(chromeVersions)
		return "Mozilla/5.0 (" + os + ") AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" +
			chromeVer + " Safari/537.36"
	case 1: // Firefox
		firefoxVer := secureRandomChoice(firefoxVersions)
		return "Mozilla/5.0 (" + os + "; rv:" + firefoxVer + ") Gecko/20100101 Firefox/" +
			firefoxVer
	default: // Safari
		safariVer := secureRandomChoice(safariVersions)
		return "Mozilla/5.0 (" + os + ") AppleWebKit/605.1.15 (KHTML, like Gecko) Version/" +
			safariVer + " Safari/605.1.15"
	}
}

func secureRandomChoice(choices []string) string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(choices))))
	if err != nil {
		return choices[0]
	}
	return choices[n.Int64()]
}

func secureRandInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

// GetLocation is a convenience function that calls GetRandomLocation
func GetLocation() GeoLocation {
	return GetRandomLocation()
}
