package util

import (
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func Timestamp() int64 {
	return time.Now().Unix()
}

func IncT(t int64, m int) int64 {
	return t + (int64)(m*60)
}

func ElapsedTimeSince(t time.Time) int64 {
	d := time.Since(t)
	return (int64)(d / time.Millisecond)
}

func TimestampToUTC(t int64) string {
	return time.Unix(t, 0).UTC().String()
}

func Random(max int) int {
	return rand.Intn(max)
}

func RandomPlusMinus(max int) int {
	p := rand.Intn(10)
	if p < 6 {
		return rand.Intn(max)
	} else {
		return 0 - rand.Intn(max)
	}
}

func ValidateUrl(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func Fingerprint(a string, b string) string {
	hash := md5.Sum([]byte(fmt.Sprint(a, b)))
	return hex.EncodeToString(hash[:])
}

func UID(a string) string {
	hash := md5.Sum([]byte(fmt.Sprint(a)))
	return hex.EncodeToString(hash[:])
}

// citation: http://play.golang.org/p/4FkNSiUDMg
func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(crand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GetJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
