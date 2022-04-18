package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
)

var UUID = uuid.NewString()
var console = bufio.NewReader(os.Stdin)

func ReadLine() (str string) {
	str, _ = console.ReadString('\n')
	str = strings.TrimSpace(str)
	return
}

func ReadNumber() (int, error) {
	str, _ := console.ReadString('\n')
	str = strings.TrimSpace(str)
	r, err := strconv.Atoi(str)
	if err == nil {
		return r, nil
	}
	return 0, err
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func DS(saltType int) string {
	var salt string
	switch saltType {
	case 0:
		salt = "fd3ykrh7o1j54g581upo1tvpam0dsgtf"
	case 1:
		salt = "14bmu1mz0yuljprsfgpvjh3ju2ni468r" // web
	case 2:
		salt = "h8w582wxwgqvahcdkpvdhbh2w9casgfl" // web_old
	}
	timeStamp := int(time.Now().Unix())
	r := RandStringBytes(6)
	ok := fmt.Sprintf("salt=%s&t=%d&r=%s", salt, timeStamp, r)
	return fmt.Sprintf("%d,%s,%s", timeStamp, r, GetMD5Hash(ok))
}

func ParseCookie(s string) map[string]string {
	ret := make(map[string]string)
	s = strings.ReplaceAll(s, " ", "")
	sp := strings.Split(s, ";")
	for _, c := range sp {
		ds := strings.Split(c, "=")
		if len(ds) < 2 {
			ret[ds[0]] = ""
		} else {
			ret[ds[0]] = ds[1]
		}
	}
	return ret
}
