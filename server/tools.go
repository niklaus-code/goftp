package server

import (
	"crypto/md5"
	"encoding/base64"
)

func encryption(pass string) string {
	base := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	ss := "123456"
	data := []byte(ss)
	has := md5.Sum(data)
	encodeString := base64.StdEncoding.EncodeToString(has[:15])

	s0 := has[15] >> 6
	s1 := has[15] & 0x3f
	rs := []rune(base)
	s2 := string(rs[s0]) + string(rs[s1])
	return encodeString + s2
}
