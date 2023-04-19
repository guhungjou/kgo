package admin

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"math/rand"
)

const (
	PTypeMD5    = "md5"
	PTypeSHA1   = "sha1"
	PTypeSHA256 = "sha256"
)

var (
	PTypeList = []string{PTypeMD5, PTypeSHA1, PTypeSHA256}
)

func randomPType() string {
	i := rand.Intn(len(PTypeList))
	return PTypeList[i]
}

func md5Encrypt(salt, plain string) string {
	s := salt + "@" + plain + "#"

	return fmt.Sprintf("%X", md5.Sum([]byte(s)))
}

func sha1Encrypt(salt, plain string) string {
	s := salt + "$" + plain + "-"
	return fmt.Sprintf("%X", sha1.Sum([]byte(s)))
}

func sha256Encrypt(salt, plain string) string {
	s := salt + "^" + plain + "+"
	return fmt.Sprintf("%X", sha256.Sum256([]byte(s)))
}

func encrypt(salt, plain, ptype string) string {
	switch ptype {
	case PTypeMD5:
		return md5Encrypt(salt, plain)
	case PTypeSHA1:
		return sha1Encrypt(salt, plain)
	case PTypeSHA256:
		return sha256Encrypt(salt, plain)
	}
	return sha256Encrypt(salt, plain)
}

func (u *AdminUser) Auth(plain string) bool {
	if u.Salt == "" || u.Password == "" {
		return false
	}
	return encrypt(u.Salt, plain, u.PType) == u.Password
}
