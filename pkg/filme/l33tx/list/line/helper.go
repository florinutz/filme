package line

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ToBytes parses a string formatted by ByteSize as bytes
func ToBytes(s string) (uint64, error) {
	err := errors.New("invalid")

	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	i := strings.IndexFunc(s, unicode.IsLetter)
	if i == -1 {
		return 0, err
	}

	bytesString, multiple := s[:i], s[i:]
	bytesString = strings.Trim(bytesString, " ")
	bytesString = strings.Replace(bytesString, ",", ".", 1)
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil || bytes < 0 {
		return 0, err
	}

	const (
		BYTE = 1 << (10 * iota)
		KILOBYTE
		MEGABYTE
		GIGABYTE
		TERABYTE
		PETABYTE
		EXABYTE
	)

	switch multiple {
	case "E", "EB", "EIB":
		return uint64(bytes * EXABYTE), nil
	case "P", "PB", "PIB":
		return uint64(bytes * PETABYTE), nil
	case "T", "TB", "TIB":
		return uint64(bytes * TERABYTE), nil
	case "G", "GB", "GIB":
		return uint64(bytes * GIGABYTE), nil
	case "M", "MB", "MIB":
		return uint64(bytes * MEGABYTE), nil
	case "K", "KB", "KIB":
		return uint64(bytes * KILOBYTE), nil
	case "B":
		return uint64(bytes), nil
	default:
		return 0, err
	}
}
