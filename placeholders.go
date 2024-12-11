package beo

import (
	"regexp"
	"strings"
	"time"
)

// processPlaceholders memproses placeholder seperti %date% dan %time%
// Mengambil format dari knowledge base jika tersedia
func processPlaceholders(answer string, kb KnowledgeBase) string {
	formats := kb.Formats
	placeholders := kb.Placeholders

	// Periksa apakah ada placeholder %date% atau %time% dalam string
	if regexp.MustCompile(`(%date%|%time%)`).MatchString(answer) {
		zone, _ := time.LoadLocation(formats.TimeZone)
		if zone == nil {
			zone = time.Local
		}
		currentTime := time.Now().In(zone)

		// Proses placeholder %date%
		if formats.Date == "" {
			formats.Date = "2006-01-02"
		}
		answer = strings.ReplaceAll(answer, "%date%", currentTime.Format(formats.Date))

		// Proses placeholder %time%
		if formats.Time == "" {
			formats.Time = "15:04:05"
		}
		answer = strings.ReplaceAll(answer, "%time%", currentTime.Format(formats.Time))
	}

	answer = strings.ReplaceAll(answer, "%ainame%", kb.AIName)
	answer = strings.ReplaceAll(answer, "%model%", kb.Model)
	answer = strings.ReplaceAll(answer, "%trainer%", kb.Trainer)

	re := regexp.MustCompile(`%(\w+)%`)
	return re.ReplaceAllStringFunc(answer, func(match string) string {
		key := match[1 : len(match)-1]
		if value, exists := placeholders[key]; exists {
			return value
		}
		return match
	})
}
