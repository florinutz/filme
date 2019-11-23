package encoding

import "fmt"

type ListEncoding int

const (
	// ids from the site
	EncDivxXvid ListEncoding = 2
	EncDvd      ListEncoding = 1
	EncH264X264 ListEncoding = 54
	EncHD       ListEncoding = 42
	EncUHD      ListEncoding = 76
	EncHevcX265 ListEncoding = 70
	EncMp4      ListEncoding = 55
	EncSvcdVcd  ListEncoding = 3
)

var encodings = map[ListEncoding]string{
	EncDivxXvid: "divx-xvid",
	EncDvd:      "dvd",
	EncH264X264: "h264-x264",
	EncHD:       "hd",
	EncUHD:      "uhd",
	EncHevcX265: "hevc-x265",
	EncMp4:      "mp4",
	EncSvcdVcd:  "svcd-vcd",
}

func (v *ListEncoding) String() string {
	return encodings[*v]
}

func (v *ListEncoding) Set(value string) (err error) {
	for id, val := range encodings {
		if val == value {
			*v = id
			return
		}
	}
	return fmt.Errorf("value '%s' is not a valid encoding", value)
}

func (*ListEncoding) Type() string {
	return "encoding"
}

func GetAll() (values []string) {
	for _, possibleValue := range encodings {
		values = append(values, possibleValue)
	}

	return
}
