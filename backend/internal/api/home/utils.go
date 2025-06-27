package home

import (
	"strings"
)

func SplitStatsText(raw string) StatsSections {
	// Replace literal `\n` with real newline characters
	processed := strings.ReplaceAll(raw, `\n`, "\n")
	processed = strings.ReplaceAll(processed, `\t`, "\t")
	processed = strings.ReplaceAll(processed, `\u003c`, "<") // decode \u003c as <

	lines := strings.Split(processed, "\n")

	var generalInfoLines, currentStatsLines []string

	section := 0 // 0 = none, 1 = general info, 2 = current stats
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "General info:") {
			section = 1
			continue
		}
		if strings.HasPrefix(line, "Current stats period:") {
			section = 2
			continue
		}

		switch section {
		case 1:
			generalInfoLines = append(generalInfoLines, line)
		case 2:
			currentStatsLines = append(currentStatsLines, line)
		}
	}

	return StatsSections{
		GeneralInfo:  strings.Join(generalInfoLines, "\n"),
		CurrentStats: strings.Join(currentStatsLines, "\n"),
	}
}
