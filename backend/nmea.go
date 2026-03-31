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

// ----------------------------------------------------------------------------
// GNSS SENTENCES
// ----------------------------------------------------------------------------

func FormatGPGGA(lat, lon float64) string {
	now := time.Now().UTC()
	timeStr := now.Format("150405.00")
	latDeg, latMin, latHemi := splitCoord(lat, "N", "S")
	lonDeg, lonMin, lonHemi := splitCoord(lon, "E", "W")

	content := fmt.Sprintf("GPGGA,%s,%02d%07.4f,%s,%03d%07.4f,%s,1,08,1.0,6.0,M,46.9,M,,", 
		timeStr, latDeg, latMin, latHemi, lonDeg, lonMin, lonHemi)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatGPRMC(lat, lon float64, speed, course float64) string {
	now := time.Now().UTC()
	timeStr := now.Format("150405.00")
	dateStr := now.Format("020106")
	latDeg, latMin, latHemi := splitCoord(lat, "N", "S")
	lonDeg, lonMin, lonHemi := splitCoord(lon, "E", "W")

	content := fmt.Sprintf("GPRMC,%s,A,%02d%07.4f,%s,%03d%07.4f,%s,%.2f,%.1f,%s,,,A", 
		timeStr, latDeg, latMin, latHemi, lonDeg, lonMin, lonHemi, speed, course, dateStr)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatGPGLL(lat, lon float64) string {
	now := time.Now().UTC()
	timeStr := now.Format("150405.00")
	latDeg, latMin, latHemi := splitCoord(lat, "N", "S")
	lonDeg, lonMin, lonHemi := splitCoord(lon, "E", "W")

	content := fmt.Sprintf("GPGLL,%02d%07.4f,%s,%03d%07.4f,%s,%s,A,A", 
		latDeg, latMin, latHemi, lonDeg, lonMin, lonHemi, timeStr)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatGPGSA() string {
	content := "GPGSA,A,3,01,02,03,04,05,06,07,08,,,,1.5,1.0,1.2"
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatGPGSV() []string {
	// Usually multiple messages. We return 2 as an example.
	c1 := "GPGSV,2,1,08,01,40,083,46,02,17,308,41,03,39,081,44,04,35,133,40"
	c2 := "GPGSV,2,2,08,05,10,211,35,06,15,100,38,07,20,050,40,08,25,000,42"
	return []string{
		"$" + c1 + "*" + GenerateNMEAChecksum(c1),
		"$" + c2 + "*" + GenerateNMEAChecksum(c2),
	}
}

// ----------------------------------------------------------------------------
// INSTRUMENTS (Talker II)
// ----------------------------------------------------------------------------

func FormatIIMWV(angle float64, speed float64) string {
	content := fmt.Sprintf("IIMWV,%.1f,R,%.1f,N,A", angle, speed)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatIIDBT(depth float64) string {
	content := fmt.Sprintf("IIDBT,%.1f,f,%.1f,M,%.1f,F", depth*3.28084, depth, depth*0.546807)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatIIDPT(depth float64) string {
	content := fmt.Sprintf("IIDPT,%.1f,0.5,50.0", depth)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatIIVHW(course, speed float64) string {
	content := fmt.Sprintf("IIVHW,%.1f,T,%.1f,M,%.1f,N,%.1f,K", course, course, speed, speed*1.852)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatIIHDM(heading float64) string {
	content := fmt.Sprintf("IIHDM,%.1f,M", heading)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatIIHDT(heading float64) string {
	content := fmt.Sprintf("IIHDT,%.1f,T", heading)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatIIMTW(temp float64) string {
	content := fmt.Sprintf("IIMTW,%.1f,C", temp)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

// ----------------------------------------------------------------------------
// AUTOPILOT
// ----------------------------------------------------------------------------

func FormatGPAPB(xte float64) string {
	content := fmt.Sprintf("GPAPB,A,A,%.2f,L,N,V,V,225.0,T,DEST01,225.5,T,226.0,T", xte)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatGPBWC(lat, lon float64) string {
	now := time.Now().UTC()
	timeStr := now.Format("150405.00")
	latDeg, latMin, latHemi := splitCoord(lat, "N", "S")
	lonDeg, lonMin, lonHemi := splitCoord(lon, "E", "W")
	content := fmt.Sprintf("GPBWC,%s,%02d%07.4f,%s,%03d%07.4f,%s,225.0,T,226.0,M,1.5,N,DEST01",
		timeStr, latDeg, latMin, latHemi, lonDeg, lonMin, lonHemi)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatGPBOD() string {
	content := "GPBOD,225.0,T,226.0,M,DEST01,ORIG01"
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

func FormatGPXTE(xte float64) string {
	content := fmt.Sprintf("GPXTE,A,A,%.2f,L,N", xte)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

// ----------------------------------------------------------------------------
// AIS (Simplified)
// ----------------------------------------------------------------------------

func FormatAIVDM() string {
	// Static payload for a Position Report Type 1 (MMSI 123456789)
	// Normal AIS encoding is complex, we'll use a pre-encoded common message
	payload := "13u9P80000S6m88N;@S:00000000"
	content := fmt.Sprintf("AIVDM,1,1,,A,%s,0", payload)
	return "$" + content + "*" + GenerateNMEAChecksum(content)
}

// ----------------------------------------------------------------------------
// HELPERS
// ----------------------------------------------------------------------------

func splitCoord(val float64, pos, neg string) (int, float64, string) {
	hemi := pos
	if val < 0 {
		val = -val
		hemi = neg
	}
	deg := int(val)
	min := (val - float64(deg)) * 60
	return deg, min, hemi
}
