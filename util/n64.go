package util

import (
	"math/rand"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
)

const n64Chars = "0123456789" +
	"abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"-_"

func FormatN64(num int64) string {
	return ""
}

func ParseN64(n64 string) int64 {
	return 0
}

type N64Generator struct {
	size int
	set  mapset.Set
}

// NewN64Generator create a new N64Generator with a new seed
func NewN64Generator(size int) *N64Generator {
	return &N64Generator{
		size: size,
		set:  mapset.NewSet(),
	}
}

func (g *N64Generator) Next() string {
	s := GenN64(g.size)
	for g.set.Contains(s) {
		s = GenN64(g.size)
	}
	g.set.Add(s)
	return s
}

func InitSeed() {
	rand.Seed(time.Now().UnixNano())
}

// GenN64 generate an N64 format ID
func GenN64(size int) string {
	var sb strings.Builder
	for i := 0; i < size; i++ {
		n := rand.Intn(64)
		sb.WriteByte(n64Chars[n])
	}
	return sb.String()
}
