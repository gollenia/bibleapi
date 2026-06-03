package main

import (
	"strconv"
	"strings"
)

func parseBibleUri(path string) BibleRequest {
	req := BibleRequest{}
	parts := strings.SplitN(path, "/", 2)
	req.BookStr = strings.ToLower(parts[0])

	if len(parts) < 2 {
		return req
	}
	locationQuery := parts[1]

	if strings.Contains(locationQuery, "-") {
		req.IsRange = true
		rangeParts := strings.SplitN(locationQuery, "-", 2)
		
		req.StartChap, req.StartVerse = parseChapterVerse(rangeParts[0])
		
		if strings.Contains(rangeParts[1], ":") {
			req.EndChap, req.EndVerse = parseChapterVerse(rangeParts[1])
		} else {
			req.EndChap = req.StartChap
			req.EndVerse, _ = strconv.Atoi(rangeParts[1])
		}
	} else {
		req.IsRange = false
		req.StartChap, req.StartVerse = parseChapterVerse(locationQuery)
		req.EndChap = req.StartChap
		req.EndVerse = req.StartVerse
	}

	return req
}

func parseChapterVerse(s string) (int, int) {
	if strings.Contains(s, ":") {
		parts := strings.SplitN(s, ":", 2)
		chap, _ := strconv.Atoi(parts[0])
		verse, _ := strconv.Atoi(parts[1])
		return chap, verse
	}
	chap, _ := strconv.Atoi(s)
	return chap, 0
}