package coll33tx

import (
	"reflect"
	"sort"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestTorrent_fromResponse(t *testing.T) {
	var torrent Torrent

	r, err := mockResponse(TestPageDetail)
	if err != nil {
		t.Fatal(err)
	}

	torrent.fromResponse(r, log.NewEntry(log.New()))

	t.Run("Title", func(t *testing.T) {
		expected := "House Party (1990) [WEBRip] [1080p] [YTS] [YIFY]"
		got := torrent.Title
		if expected != got {
			t.Errorf("expected Title '%s', got '%s'", expected, got)
		}
	})

	t.Run("FilmCleanTitle", func(t *testing.T) {
		expected := "House Party"
		got := torrent.FilmCleanTitle
		if expected != got {
			t.Errorf("expected FilmCleanTitle '%s', got '%s'", expected, got)
		}
	})

	t.Run("FilmLink", func(t *testing.T) {
		if torrent.FilmLink == nil {
			t.Skip("FilmLink is nil")
		}
		expected := "https://1337x.to/movie/16094/House-Party-1990/"
		got := torrent.FilmLink.String()
		if expected != got {
			t.Errorf("expected FilmLink '%s', got '%s'", expected, got)
		}
	})

	t.Run("FilmDescription", func(t *testing.T) {
		expected := "Young Kid has been invited to a party at his friend Play's house. But after a fight at school, " +
			"Kid's father grounds him. None the less, Kid sneaks out when his father falls asleep. But Kid doesn't k" +
			"now that three of the thugs at school has decided to give him a lesson in behaviour"
		got := torrent.FilmDescription
		if expected != got {
			t.Errorf("expected FilmLink '%s', got '%s'", expected, got)
		}
	})

	t.Run("Magnet", func(t *testing.T) {
		expected := "magnet:?xt=urn:btih:7F736E2E527ADF321D94092AA3DDA2C326EF4F31&dn=House+Party+%281990%29+%5BWEBRi" +
			"p%5D+%5B1080p%5D+%5BYTS%5D+%5BYIFY%5D&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3" +
			"A%2F%2F9.rarbg.com%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337&tr=udp%3A%2F%2Ftracker.inter" +
			"netwarriors.net%3A1337&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.z" +
			"er0day.to%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F" +
			"%2Fcoppersurfer.tk%3A6969%2Fannounce"
		got := torrent.Magnet
		if expected != got {
			t.Errorf("expected Magnet '%s', got '%s'", expected, got)
		}
	})

	t.Run("FoundOn", func(t *testing.T) {
		if torrent.FoundOn == nil {
			t.Error("FoundOn is nil")
		}
		expected := TestPageDetail
		got := torrent.FoundOn.String()
		if expected != got {
			t.Errorf("expected FoundOn '%s', got '%s'", expected, got)
		}
	})

	t.Run("Image", func(t *testing.T) {
		if torrent.Image == nil {
			t.Skip("Image is nil")
		}
		expected := "//lx1.dyncdn.cc/cdn/b7/b731df6a1946886ab20a289df553380d.jpg"
		got := torrent.Image.String()
		if expected != got {
			t.Errorf("expected Image '%s', got '%s'", expected, got)
		}
	})

	t.Run("Description", func(t *testing.T) {
		expected := ""
		got := torrent.Description
		if expected != got {
			t.Errorf("expected Description '%s', got '%s'", expected, got)
		}
	})

	t.Run("Genres", func(t *testing.T) {
		expected := []string{"Comedy", "Drama"}
		got := torrent.Genres
		if !reflect.DeepEqual(expected, got) {
			sort.Strings(expected)
			sort.Strings(got)
			t.Errorf("expected FilmCategories '%s', got '%s'", expected, got)
		}
	})

	t.Run("Category", func(t *testing.T) {
		expected := "Movies"
		got := torrent.Category
		if expected != got {
			t.Errorf("expected Category '%s', got '%s'", expected, got)
		}
	})

	t.Run("Type", func(t *testing.T) {
		expected := "HD"
		got := torrent.Type
		if expected != got {
			t.Errorf("expected Type '%s', got '%s'", expected, got)
		}
	})

	t.Run("Language", func(t *testing.T) {
		expected := "English"
		got := torrent.Language
		if expected != got {
			t.Errorf("expected Language '%s', got '%s'", expected, got)
		}
	})

	t.Run("TotalSize", func(t *testing.T) {
		expected := "1.7 GB"
		got := torrent.TotalSize
		if expected != got {
			t.Errorf("expected TotalSize '%s', got '%s'", expected, got)
		}
	})

	t.Run("UploadedBy", func(t *testing.T) {
		expected := "YTSAGx"
		got := torrent.UploadedBy
		if expected != got {
			t.Errorf("expected UploadedBy '%s', got '%s'", expected, got)
		}
	})

	t.Run("Downloads", func(t *testing.T) {
		expected := 341
		got := torrent.Downloads
		if expected != got {
			t.Errorf("expected '%d' Downloads, got '%d'", expected, got)
		}
	})

	t.Run("LastChecked", func(t *testing.T) {
		expected := "1 week ago"
		got := torrent.LastChecked
		if expected != got {
			t.Errorf("expected LastChecked '%s', got '%s'", expected, got)
		}
	})

	t.Run("DateUploaded", func(t *testing.T) {
		expected := "2 weeks ago"
		got := torrent.DateUploaded
		if expected != got {
			t.Errorf("expected DateUploaded '%s', got '%s'", expected, got)
		}
	})

	t.Run("Seeders", func(t *testing.T) {
		expected := 158
		got := torrent.Seeders
		if expected != got {
			t.Errorf("expected '%d' Seeders, got '%d'", expected, got)
		}
	})

	t.Run("Leechers", func(t *testing.T) {
		expected := 38
		got := torrent.Leechers
		if expected != got {
			t.Errorf("expected '%d' Leechers, got '%d'", expected, got)
		}
	})
}
