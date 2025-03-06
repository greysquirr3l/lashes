package agent

import (
	"math/rand"
	"time"
)

type GeoLocation struct {
	Latitude  float64
	Longitude float64
	Timezone  string
	Locale    string
}

var locations = []GeoLocation{
	{40.7128, -74.0060, "America/New_York", "en-US"},  // New York
	{51.5074, -0.1278, "Europe/London", "en-GB"},      // London
	{48.8566, 2.3522, "Europe/Paris", "fr-FR"},        // Paris
	{35.6762, 139.6503, "Asia/Tokyo", "ja-JP"},        // Tokyo
	{52.5200, 13.4050, "Europe/Berlin", "de-DE"},      // Berlin
	{-33.8688, 151.2093, "Australia/Sydney", "en-AU"}, // Sydney
	{55.7558, 37.6173, "Europe/Moscow", "ru-RU"},      // Moscow
	{1.3521, 103.8198, "Asia/Singapore", "en-SG"},     // Singapore
}

func GetRandomLocation() GeoLocation {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	loc := locations[r.Intn(len(locations))]

	// Add some random noise to coordinates
	loc.Latitude += (r.Float64() - 0.5) * 0.01
	loc.Longitude += (r.Float64() - 0.5) * 0.01

	return loc
}
