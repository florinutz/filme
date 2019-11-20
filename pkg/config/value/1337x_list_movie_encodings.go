package value

import "fmt"

type LeetxListEncoding int

const (
	// ids from the site
	EncodingDivxXvid LeetxListEncoding = 2
	EncodingDvd      LeetxListEncoding = 1
	EncodingH264X264 LeetxListEncoding = 54
	EncodingHD       LeetxListEncoding = 42
	EncodingUHD      LeetxListEncoding = 76
	EncodingHevcX265 LeetxListEncoding = 70
	EncodingMp4      LeetxListEncoding = 55
	EncodingSvcdVcd  LeetxListEncoding = 3
)

var encodings = map[LeetxListEncoding]string{
	EncodingDivxXvid: "divx-xvid",
	EncodingDvd:      "dvd",
	EncodingH264X264: "h264-x264",
	EncodingHD:       "hd",
	EncodingUHD:      "uhd",
	EncodingHevcX265: "hevc-x265",
	EncodingMp4:      "mp4",
	EncodingSvcdVcd:  "svcd-vcd",
}

func (v *LeetxListEncoding) String() string {
	return encodings[*v]
}

func (v *LeetxListEncoding) Set(value string) (err error) {
	for id, val := range encodings {
		if val == value {
			*v = id
			return
		}
	}
	return fmt.Errorf("value '%s' is not a valid encoding", value)
}

func (*LeetxListEncoding) Type() string {
	return "encoding"
}

func GetAllLeetxListMovieEncodingValues() (values []string) {
	for _, possibleValue := range encodings {
		values = append(values, possibleValue)
	}

	return
}
