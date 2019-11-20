package value

import "fmt"

type LeetxListSearchCategory uint8

const (
	SearchCategoryAll LeetxListSearchCategory = iota
	SearchCategoryMovies
	SearchCategoryTV
	SearchCategoryDocumentaries
	SearchCategoryAnime
	SearchCategoryXXX
)

var possibleCategoryValues = map[LeetxListSearchCategory]string{
	SearchCategoryAll:           "all",
	SearchCategoryMovies:        "movies",
	SearchCategoryTV:            "tv",
	SearchCategoryDocumentaries: "documentaries",
	SearchCategoryAnime:         "anime",
	SearchCategoryXXX:           "xxx",
}

func (v *LeetxListSearchCategory) TranslateToUrlParam() string {
	switch *v {
	case SearchCategoryMovies:
		return "Movies"
	case SearchCategoryTV:
		return "TV"
	case SearchCategoryDocumentaries:
		return "Documentaries"
	case SearchCategoryAnime:
		return "Anime"
	case SearchCategoryXXX:
		return "XXX"
	default:
		return ""
	}
}

func (v *LeetxListSearchCategory) String() string {
	if v, ok := possibleCategoryValues[*v]; ok {
		return v
	}
	return ""
}

func (v *LeetxListSearchCategory) Set(value string) (err error) {
	for key, val := range possibleCategoryValues {
		if val == value {
			*v = key
			return
		}
	}
	return fmt.Errorf("value '%s' is not a valid search category flag", value)
}

func (*LeetxListSearchCategory) Type() string {
	return "category"
}

func GetAllLeetxListMovieSearchCategoryValues() (values []string) {
	for _, possibleValue := range possibleCategoryValues {
		values = append(values, possibleValue)
	}
	return
}
