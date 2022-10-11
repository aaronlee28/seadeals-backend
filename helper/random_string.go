package helper

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var driver = []string{
	"Sena Narumi",
	"Mona Narumi",
	"Yuu Setoguchi",
	"Natsuki Enomoto",
	"Arisa Takamizawa",
	"Ken Shibasaki",
	"Hina Setoguchi",
	"Kotarou Enomoto",
	"Aizou Shibasaki",
	"Yuujirou Someya",
}

func RandomString(n int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r1.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandomDriver() string {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	return driver[rand.Intn(len(driver))]
}
