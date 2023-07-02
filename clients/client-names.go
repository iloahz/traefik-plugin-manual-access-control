package clients

import (
	"math/rand"

	"github.com/lucasepe/codename"
)

var (
	rng *rand.Rand
)

func init() {
	var err error
	rng, err = codename.DefaultRNG()
	if err != nil {
		panic(err)
	}
}

func GenerateName() string {
	return codename.Generate(rng, 0)
}
