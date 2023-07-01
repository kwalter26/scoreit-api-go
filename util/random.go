package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var result strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len(alphabet))]
		result.WriteByte(c)
	}

	return result.String()
}

func RandomName() string {
	return RandomString(6)
}

func RandomBalance() int64 {
	return rand.Int63n(1000)
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomBaseballPosition returns a random baseball position from a list of positions from positions.go
func RandomBaseballPosition() BaseballPosition {
	list := []BaseballPosition{Pitcher, Catcher, FirstBase, SecondBase, ThirdBase, ShortStop, LeftField, CenterField, RightField, DesignatedHitter, RightCenterField, LeftCenterField}
	return list[rand.Intn(len(list))]
}
