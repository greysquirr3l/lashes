package agent

import (
	"math/rand"
	"time"
)

// Common browser versions
var (
	chromeVersions  = []string{"120.0.0", "121.0.0", "122.0.0"}
	firefoxVersions = []string{"122.0", "123.0", "124.0"}
	safariVersions  = []string{"17.2", "17.1", "16.6"}
	osVersions      = []string{
		"Windows NT 10.0; Win64; x64",
		"Macintosh; Intel Mac OS X 10_15_7",
		"Macintosh; Intel Mac OS X 11_6_0",
		"X11; Linux x86_64",
		"X11; Ubuntu; Linux x86_64",
	}
)

func GetRandomUserAgent() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	switch r.Intn(3) {
	case 0:
		return "Mozilla/5.0 (" + randomChoice(osVersions) + ") AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" +
			randomChoice(chromeVersions) + " Safari/537.36"
	case 1:
		return "Mozilla/5.0 (" + randomChoice(osVersions) + "; rv:" +
			randomChoice(firefoxVersions) + ") Gecko/20100101 Firefox/" + randomChoice(firefoxVersions)
	default:
		return "Mozilla/5.0 (" + randomChoice(osVersions) + ") AppleWebKit/605.1.15 (KHTML, like Gecko) Version/" +
			randomChoice(safariVersions) + " Safari/605.1.15"
	}
}

func randomChoice(choices []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return choices[r.Intn(len(choices))]
}
