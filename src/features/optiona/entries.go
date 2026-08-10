package optiona

import (
	"fmt"
	"strings"
)

// entries is placeholder data so the paged list has enough rows to page and
// search through. Replace it with a real source when adapting this feature.
var entries = buildEntries()

type entry struct {
	Label  string
	Detail string
}

func buildEntries() []entry {
	words := []string{
		"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf",
		"Hotel", "India", "Juliett", "Kilo", "Lima", "Mike", "November",
		"Oscar", "Papa", "Quebec", "Romeo", "Sierra", "Tango", "Uniform",
		"Victor", "Whiskey", "Xray", "Yankee", "Zulu",
	}

	built := make([]entry, 0, len(words)*5)
	for round := 1; round <= 5; round++ {
		for _, word := range words {
			built = append(built, entry{
				Label:  fmt.Sprintf("%s %d", word, round),
				Detail: fmt.Sprintf("#%03d", len(built)+1),
			})
		}
	}
	return built
}

func matching(query string) []entry {
	if query == "" {
		return entries
	}

	needle := strings.ToLower(query)
	matches := make([]entry, 0, len(entries))
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Label), needle) {
			matches = append(matches, e)
		}
	}
	return matches
}
