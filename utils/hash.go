package utils

import (
	"crypto/md5"
	"time"
)

func TripHash(busNo string,date time.Time) string {
	// hash the bus number and date using md5 and return as string
	hash := md5.New()
	hash.Write([]byte(busNo))
	hash.Write([]byte(date.String()))
	return string(hash.Sum(nil))
}