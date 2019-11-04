package search

type FilmSources map[string][]BaseItem

func GetFilmSources(items map[int]BaseItem) FilmSources {
	sources := make(FilmSources)

	for _, item := range items {
		if item.Url() == nil {
			continue
		}

		var key string

		switch item.(type) {
		case *ItemNetflix:
			key = "Netflix"
		case *ItemPrime:
			key = "Prime"
		case *ItemImdb:
			key = "Imdb"
		case *ItemWikipedia:
			key = "Wikipedia"
		case *ItemMetacritic:
			key = "Metacritic"
		case *ItemRottenTomatoes:
			key = "Rotten Tomatoes"
		default:
			key = "other"
		}

		sources[key] = append(sources[key], item)
	}

	return sources
}
