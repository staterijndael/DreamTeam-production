package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"time"
)

var (
	phone = regexp.
		MustCompile(`\+\s*\d\s*\(?\s*\d\s*\d\s*\d\s*\)?\s*\d\s*\d\s*\d\s*-?\s*\d\s*\d\s*-?\s*\d\s*\d`)
	nums           = []byte("0123456789")
	plus           = byte('+')
	openingBracket = byte('(')
	closingBracket = byte(')')
	dash           = byte('-')
	space          = byte(' ')
)

type Container struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func FileCheckSum(f *os.File) (*string, error) {
	digest := md5.New()
	if _, err := io.Copy(digest, f); err != nil {
		return nil, err
	}

	hash := digest.Sum(nil)
	checksumString := hex.EncodeToString(hash)
	return &checksumString, nil
}

func BytesChecksum(b []byte) string {
	digest := md5.New()
	digest.Write(b)
	return hex.EncodeToString(digest.Sum(nil))
}

func NewMediaFile(dir string) (*os.File, *string, error) {
	now := time.Now().UnixNano()
	fp := fmt.Sprintf("%s/%v", dir, now)
	f, err := os.Create(fp)

	if err != nil {
		return nil, nil, err
	}

	return f, &fp, err
}

func FormatPhone(s string) string {
	pn := make([]byte, 0)
	for _, ch := range []byte(s) {
		for _, n := range nums {
			if n == ch {
				pn = append(pn, ch)
				break
			}
		}
	}

	return string([]byte{
		plus, pn[0],
		space, openingBracket,
		pn[1], pn[2], pn[3],
		closingBracket, space,
		pn[4], pn[5], pn[6],
		dash, pn[7], pn[8],
		dash, pn[9], pn[10],
	})
}

func IsValidOrgTitle(s string) bool {
	return true
}

func IsValidOrgDescription(s string) bool {
	return true
}

func IsValidGroupTitle(s string) bool {
	return true
}

func IsPhone(b []byte) bool {
	return phone.Match(b)
}

func IntSliceSum(arr []int32) int32 {
	var sum int32
	for _, a := range arr {
		sum += a
	}
	return sum
}

func ConvertToScale(score, min, max, scale int32) int32 {
	scoreF := float64(score)
	minF := float64(min)
	maxF := float64(max)
	scaleF := float64(scale)

	return int32(math.Ceil((scoreF - minF) / ((maxF - minF) / scaleF)))
}
