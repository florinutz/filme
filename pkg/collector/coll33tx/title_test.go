package coll33tx

import (
	"testing"
)

func Test_ParseTitleInfo(t *testing.T) {
	tests := [...]struct {
		title string
		info  TitleInfo
	}{
		{
			title: "1917 (2019) [1080p] [WEBRip] [5.1] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "1917",
				Year:    2019,
				Quality: "1080p",
			},
		},
		{
			title: "Russell Peters: Deported (2020) [720p] [WEBRip] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Russell Peters: Deported",
				Year:    2020,
				Quality: "720p",
			},
		},
		{
			title: "Mortal.Kombat.Legends.Scorpions.Revenge.2020.720p.WEBRip.x264.AAC-ETRG",
			info: TitleInfo{
				Title:   "Mortal Kombat Legends Scorpions Revenge",
				Year:    2020,
				Quality: "720p",
			},
		},
		{
			title: "Mortal Kombat Legends: Scorpions Revenge (2020) [720p] [WEBRip] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Mortal Kombat Legends: Scorpions Revenge",
				Year:    2020,
				Quality: "720p",
			},
		},
		{
			title: "Blue.Story.2019.720p.WEBRip.800MB.x264-GalaxyRG ⭐",
			info: TitleInfo{
				Title:   "Blue Story",
				Year:    2019,
				Quality: "720p",
			},
		},
		{
			title: "David Blaine The Magic Way 2020 720p WEB-DL H264 BONE",
			info: TitleInfo{
				Title:   "David Blaine The Magic Way",
				Year:    2020,
				Quality: "720p",
			},
		},
		{
			title: "Killer.Tattooist.2020.720p.WEBRip.800MB.x264-GalaxyRG ⭐",
			info: TitleInfo{
				Title:   "Killer Tattooist",
				Year:    2020,
				Quality: "720p",
			},
		},
		{
			title: "Mortal Kombat Legends - Scorpions Revenge (2020) (1080p AMZN Webrip x265 10bit E...",
			info: TitleInfo{
				Title:   "Mortal Kombat Legends - Scorpions Revenge",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Mortal Kombat Legends: Scorpions Revenge (2020) [1080p] [WEBRip] [5.1] [YTS] [YI...",
			info: TitleInfo{
				Title:   "Mortal Kombat Legends: Scorpions Revenge",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "The Rhythm Section (2020) (1080p NF WEB-DL x265 HEVC 10bit AAC 5.1 Q22 Joy) [UTR...",
			info: TitleInfo{
				Title:   "The Rhythm Section",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "The Wave (2019) [1080p] [WEBRip] [5.1] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "The Wave",
				Year:    2019,
				Quality: "1080p",
			},
		},
		{
			title: "Fractured (2019) [720p] [WEBRip] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Fractured",
				Year:    2019,
				Quality: "720p",
			},
		},
		{
			title: "Trolls World Tour (2020) ITA-ENG Ac3 5.1 WEBRip 1080p H264 [ArMor]",
			info: TitleInfo{
				Title:   "Trolls World Tour",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Alien Expedition (2018).1080p.H264.ita.eng.Ac3-5.1.sub.eng-MIRCrew",
			info: TitleInfo{
				Title:   "Alien Expedition",
				Year:    2018,
				Quality: "1080p",
			},
		},
		{
			title: "Russell Peters: Deported (2020) [1080p] [WEBRip] [5.1] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Russell Peters: Deported",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Blue.Story.2019.HDRip.XviD.AC3-EVO[TGx] ⭐",
			info: TitleInfo{
				Title:   "Blue Story",
				Year:    2019,
				Quality: "HDRip",
			},
		},
		{
			title: "Superman: Red Son (2020) + Extras (1080p BluRay x265 HEVC 10bit DTS 5.1 SAMPA) [...",
			info: TitleInfo{
				Title:   "Superman: Red Son",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Sound & Fury (2019) [720p] [WEBRip] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Sound & Fury",
				Year:    2019,
				Quality: "720p",
			},
		},
		{
			title: "Enemy Mine (1985) + Extras (1080p BluRay x265 HEVC 10bit AAC 5.0 r00t) [QxR]",
			info: TitleInfo{
				Title:   "Enemy Mine",
				Year:    1985,
				Quality: "1080p",
			},
		},
		{
			title: "28 Hotel Rooms 2012 1080p WEBRip x265-RARBG",
			info: TitleInfo{
				Title:   "28 Hotel Rooms",
				Year:    2012,
				Quality: "1080p",
			},
		},
		{
			title: "Ford v Ferrari (2019) 1080p 5.1 - 2.0 x264 Phun Psyz",
			info: TitleInfo{
				Title:   "Ford v Ferrari",
				Year:    2019,
				Quality: "1080p",
			},
		},
		{
			title: "The.Clark.Sisters.First.Ladies.of.Gospel.2020.HDTV.x264-CRiMSON[TGx] ⭐",
			info: TitleInfo{
				Title:   "The Clark Sisters First Ladies of Gospel",
				Year:    2020,
				Quality: "HDTV",
			},
		},
		{
			title: "Jason Bourne (2016) 720p BluRay 10bit x265 HEVC Dual Audio Hindi 5.1 English 5.1...",
			info: TitleInfo{
				Title:   "Jason Bourne",
				Year:    2016,
				Quality: "720p",
			},
		},
		{
			title: "Fractured (2019) [1080p] [WEBRip] [5.1] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Fractured",
				Year:    2019,
				Quality: "1080p",
			},
		},
		{
			title: "Blue.Story.2019.1080p.WEB-DL.H264.AC3-EVO[EtHD]",
			info: TitleInfo{
				Title:   "Blue Story",
				Year:    2019,
				Quality: "1080p",
			},
		},
		{
			title: "Girls.Of.the.Sun.2018.DVDRip.x264-PFa[TGx] ⭐",
			info: TitleInfo{
				Title:   "Girls Of the Sun",
				Year:    2018,
				Quality: "DVDRip",
			},
		},
		{
			title: "Sound & Fury (2019) [1080p] [WEBRip] [5.1] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Sound & Fury",
				Year:    2019,
				Quality: "1080p",
			},
		},
		{
			title: "Bad Boys for Life (2020) (1080p BluRay x265 HEVC 10bit AAC 5.1 Tigole) [QxR]",
			info: TitleInfo{
				Title:   "Bad Boys for Life",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Lost.Girls.2020.1080p.WEBRip.x264-WATCHER[TGx] ⭐",
			info: TitleInfo{
				Title:   "Lost Girls",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Pete Davidson: Alive from New York (2020) [1080p] [BluRay] [5.1] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Pete Davidson: Alive from New York",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Mortal Kombat Legends Scorpions Revenge (2020) 1080p 5.1 - 2.0 x264 Phun Psyz",
			info: TitleInfo{
				Title:   "Mortal Kombat Legends Scorpions Revenge",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Ski Patrol (1990) [720p] [WEBRip] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Ski Patrol",
				Year:    1990,
				Quality: "720p",
			},
		},
		{
			title: "Spider in the Web - Armi chimiche (2019).720p.H264.ita.eng.Ac3-5.1.sub.eng-MIRCr...",
			info: TitleInfo{
				Title:   "Spider in the Web - Armi chimiche",
				Year:    2019,
				Quality: "720p",
			},
		},
		{
			title: "Street Fighter (1994) [720p] [BluRay] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Street Fighter",
				Year:    1994,
				Quality: "720p",
			},
		},
		{
			title: "Like a Boss-Amiche in affari (2020) ITA-ENG Ac3 5.1 BDRip 1080p H264 [ArMor]",
			info: TitleInfo{
				Title:   "Like a Boss-Amiche in affari",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Fame (1980) [720p] [BluRay] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Fame",
				Year:    1980,
				Quality: "720p",
			},
		},
		{
			title: "Pete Davidson: Alive from New York (2020) [720p] [BluRay] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Pete Davidson: Alive from New York",
				Year:    2020,
				Quality: "720p",
			},
		},
		{
			title: "Zathura - A Space Adventure (2005) (1080p BluRay x265 HEVC 10bit AAC 5.1 Tigole)...",
			info: TitleInfo{
				Title:   "Zathura - A Space Adventure",
				Year:    2005,
				Quality: "1080p",
			},
		},
		{
			title: "Killer.Tattooist.2020.HDRip.XviD.AC3-EVO[TGx] ⭐",
			info: TitleInfo{
				Title:   "Killer Tattooist",
				Year:    2020,
				Quality: "HDRip",
			},
		},
		{
			title: "Blue.Story.2019.1080p.WEB-DL.H264.AC3-EVO[TGx] ⭐",
			info: TitleInfo{
				Title:   "Blue Story",
				Year:    2019,
				Quality: "1080p",
			},
		},
		{
			title: "Ski Patrol (1990) [1080p] [WEBRip] [2.0] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Ski Patrol",
				Year:    1990,
				Quality: "1080p",
			},
		},
		{
			title: "The Hunt (2020) (1080p WEB x265 HEVC 10bit AAC 5.1 Q22 Joy) [UTR]",
			info: TitleInfo{
				Title:   "The Hunt",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "The Bourne Identity (2002) 1080p BluRay 10bit x265 HEVC Dual Audio Hindi 5.1 Eng...",
			info: TitleInfo{
				Title:   "The Bourne Identity",
				Year:    2002,
				Quality: "1080p",
			},
		},
		{
			title: "Fame (1980) [1080p] [BluRay] [5.1] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Fame",
				Year:    1980,
				Quality: "1080p",
			},
		},
		{
			title: "Blue.Story.2019.1080p.WEBRip.1400MB.DD5.1.x264-GalaxyRG ⭐",
			info: TitleInfo{
				Title:   "Blue Story",
				Year:    2019,
				Quality: "1080p",
			},
		},
		{
			title: "The Fighting Renegade (Western 1939) Tim McCoy 720p",
			info: TitleInfo{
				Title:   "The Fighting Renegade",
				Year:    1939,
				Quality: "720p",
			},
		},
		{
			title: "The Bourne Legacy (2012) 720p BluRay 10bit x265 HEVC Dual Audio Hindi 5.1 Englis...",
			info: TitleInfo{
				Title:   "The Bourne Legacy",
				Year:    2012,
				Quality: "720p",
			},
		},
		{
			title: "Lost Girls (2020) [720p] [WEBRip] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Lost Girls",
				Year:    2020,
				Quality: "720p",
			},
		},
		{
			title: "The.Turning.2020.1080p.10bit.BluRay.8CH.x265.HEVC-PSA",
			info: TitleInfo{
				Title:   "The Turning",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "Lost Girls (2020) [1080p] [WEBRip] [5.1] [YTS] [YIFY]",
			info: TitleInfo{
				Title:   "Lost Girls",
				Year:    2020,
				Quality: "1080p",
			},
		},
		{
			title: "The Main Event-Sognando il ring (2020) ITA-ENG Ac3 5.1 WEBRip 1080p H264 [ArMor]",
			info: TitleInfo{
				Title:   "The Main Event-Sognando il ring",
				Year:    2020,
				Quality: "1080p",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			info := ParseTitleInfo(test.title)
			if info.Title != test.info.Title {
				t.Errorf("title = '%v', want '%v'", info.Title, test.info.Title)
			}
			if info.Quality != test.info.Quality {
				t.Errorf("quality = '%v', want '%v'", info.Quality, test.info.Quality)
			}
			if info.Year != test.info.Year {
				t.Errorf("year = '%v', want '%v'", info.Year, test.info.Year)
			}
		})
	}
}
