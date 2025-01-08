// Package distance distance
package distance

import "math"

// CalcDistanceCheck 计算距离
func CalcDistanceCheck(lat, lng, orgLat, orgLng float64) int {
	if lat == 0 || lng == 0 || orgLat == 0 || orgLng == 0 {
		return 0
	}
	return CalcDistance(lat, lng, orgLat, orgLng)
}

// CalcDistance 计算距离
func CalcDistance(lat, lng, orgLat, orgLng float64) int {
	const earthRadius = 6378.138
	const pi = math.Pi
	latRad := lat * pi / 180
	lngRad := lng * pi / 180
	orgLatRad := orgLat * pi / 180
	orgLngRad := orgLng * pi / 180
	sinLat := math.Sin((latRad - orgLatRad) / 2)
	sinLng := math.Sin((lngRad - orgLngRad) / 2)
	cosLat := math.Cos(latRad)
	cosOrgLat := math.Cos(orgLatRad)
	a := sinLat*sinLat + cosLat*cosOrgLat*sinLng*sinLng
	c := 2 * math.Asin(math.Sqrt(a))
	distance := earthRadius * c * 1000
	return int(math.Round(distance))
}
