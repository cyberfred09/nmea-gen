package backend

import (
	"fmt"
	"time"
)

// GenerateNMEAChecksum computes the XOR checksum for an NMEA sentence string.
func GenerateNMEAChecksum(s string) string {
	checksum := 0
	for i := 0; i < len(s); i++ {
		checksum ^= int(s[i])
	}
	return fmt.Sprintf("%02X", checksum)
}

// FormatGPGGA generates a GPGGA sentence (Global Positioning System Fix Data).
func FormatGPGGA(lat, lon float64, fixQuality int, numSatellites int, hdop float64, altitude float64) string {
	now := time.Now().UTC()
	timeStr := now.Format("150405.00")
	
	latDeg := int(lat)
	latMin := (lat - float64(latDeg)) * 60
	latHemi := "N"
	if lat < 0 {
		latDeg = -latDeg
		latMin = -latMin
		latHemi = "S"
	}
	
	lonDeg := int(lon)
	lonMin := (lon - float64(lonDeg)) * 60
	lonHemi := "E"
	if lon < 0 {
		lonDeg = -lonDeg
		lonMin = -lonMin
		lonHemi = "W"
	}

	content := fmt.Sprintf("GPGGA,%s,%02d%07.4f,%s,%03d%07.4f,%s,%d,%02d,%.1f,%.1f,M,46.9,M,,", 
		timeStr, latDeg, latMin, latHemi, lonDeg, lonMin, lonHemi, fixQuality, numSatellites, hdop, altitude)
	
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

// FormatGPRMC generates a GPRMC sentence (Recommended Minimum Navigation Information).
func FormatGPRMC(lat, lon float64, speed float64, course float64) string {
	now := time.Now().UTC()
	timeStr := now.Format("150405.00")
	dateStr := now.Format("020106")
	
	latDeg := int(lat)
	latMin := (lat - float64(latDeg)) * 60
	latHemi := "N"
	if lat < 0 {
		latDeg = -latDeg
		latMin = -latMin
		latHemi = "S"
	}
	
	lonDeg := int(lon)
	lonMin := (lon - float64(lonDeg)) * 60
	lonHemi := "E"
	if lon < 0 {
		lonDeg = -lonDeg
		lonMin = -lonMin
		lonHemi = "W"
	}

	content := fmt.Sprintf("GPRMC,%s,A,%02d%07.4f,%s,%03d%07.4f,%s,%.2f,%.1f,%s,,,A", 
		timeStr, latDeg, latMin, latHemi, lonDeg, lonMin, lonHemi, speed, course, dateStr)
	
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

// FormatGPVTG generates a GPVTG sentence (Track Made Good and Ground Speed).
func FormatGPVTG(course float64, speedKnots float64) string {
	speedKmh := speedKnots * 1.852
	content := fmt.Sprintf("GPVTG,%.1f,T,,M,%.2f,N,%.2f,K,A", course, speedKnots, speedKmh)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}
