package agent

import (
	"crypto/rand"
	"math/big"
)

type GeoLocation struct {
	Latitude  float64
	Longitude float64
	Timezone  string
	Locale    string
}

var locations = []GeoLocation{
	{40.7128, -74.0060, "America/New_York", "en-US"},     // New York
	{51.5074, -0.1278, "Europe/London", "en-GB"},         // London
	{48.8566, 2.3522, "Europe/Paris", "fr-FR"},           // Paris
	{35.6762, 139.6503, "Asia/Tokyo", "ja-JP"},           // Tokyo
	{52.5200, 13.4050, "Europe/Berlin", "de-DE"},         // Berlin
	{-33.8688, 151.2093, "Australia/Sydney", "en-AU"},    // Sydney
	{55.7558, 37.6173, "Europe/Moscow", "ru-RU"},         // Moscow
	{1.3521, 103.8198, "Asia/Singapore", "en-SG"},        // Singapore
	{39.9042, 116.4074, "Asia/Shanghai", "zh-CN"},        // Beijing
	{19.4326, -99.1332, "America/Mexico_City", "es-MX"},  // Mexico City
	{37.7749, -122.4194, "America/Los_Angeles", "en-US"}, // San Francisco
	{-23.5505, -46.6333, "America/Sao_Paulo", "pt-BR"},   // São Paulo
	{28.6139, 77.2090, "Asia/Kolkata", "hi-IN"},          // New Delhi
	{34.0522, -118.2437, "America/Los_Angeles", "en-US"}, // Los Angeles
	{41.9028, 12.4964, "Europe/Rome", "it-IT"},           // Rome
	{25.2048, 55.2708, "Asia/Dubai", "ar-AE"},            // Dubai
	{-6.2088, 106.8456, "Asia/Jakarta", "id-ID"},         // Jakarta
	{59.3293, 18.0686, "Europe/Stockholm", "sv-SE"},      // Stockholm
	{37.5665, 126.9780, "Asia/Seoul", "ko-KR"},           // Seoul
	{43.6532, -79.3832, "America/Toronto", "en-CA"},      // Toronto
}

// GetRandomLocation returns a randomized location from a major city
// with slight coordinate variations for privacy
func GetRandomLocation() GeoLocation {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(locations))))
	if err != nil {
		n = big.NewInt(0)
	}
	loc := locations[n.Int64()]

	// Securely generate small variations to coordinates
	v1, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		v1 = big.NewInt(0)
	}
	v2, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		v2 = big.NewInt(0)
	}
	loc.Latitude += (float64(v1.Int64())/10000.0 - 0.5) * 0.02
	loc.Longitude += (float64(v2.Int64())/10000.0 - 0.5) * 0.02

	return loc
}

// GetLocationByTimezone returns a location matching the specified timezone
func GetLocationByTimezone(timezone string) (GeoLocation, bool) {
	for _, loc := range locations {
		if loc.Timezone == timezone {
			return loc, true
		}
	}
	return GeoLocation{}, false
}

// GetLocationByLocale returns a location matching the specified locale
func GetLocationByLocale(locale string) (GeoLocation, bool) {
	for _, loc := range locations {
		if loc.Locale == locale {
			return loc, true
		}
	}
	return GeoLocation{}, false
}
