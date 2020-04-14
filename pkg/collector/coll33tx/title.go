package coll33tx

import (
	"regexp"
	"strconv"
	"strings"
)

type TitleInfo struct {
	Title   string
	Year    int
	Quality string
}

// https://regex101.com/r/4TeA0N/3
var titleSplitRE = regexp.MustCompile(`(?i)(?P<title>(?:.+?))(?:[.\s\[({])*(?P<year>(?:19\d{2})|(?:20[012]\d{1})).*?(?P<quality>(?:\d{3,4}p)|(?:dvdrip)|(?:brrip)|(?:dvdscr)|(?:hdrip)|(?:hdtv))`)

func ParseTitleInfo(title string) TitleInfo {
	info := TitleInfo{}
	groupNames := titleSplitRE.SubexpNames()
	for _, match := range titleSplitRE.FindAllStringSubmatch(title, -1) {
		for groupIdx, group := range match {
			switch groupNames[groupIdx] {
			case "title":
				info.Title = strings.Replace(group, ".", " ", -1)
				firstInvalidCharPosition := strings.IndexAny(info.Title, `[({`)
				if firstInvalidCharPosition > 0 {
					info.Title = info.Title[0:firstInvalidCharPosition]
				}
				info.Title = strings.Trim(info.Title, " ")
			case "year":
				info.Year, _ = strconv.Atoi(group)
			case "quality":
				info.Quality = group
			}
		}
	}
	return info
}
