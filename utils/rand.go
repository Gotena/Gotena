package utils

import (
	"math/rand"
	"time"
)

type UniqueRand struct {
	generated map[int]bool //keeps track of
	rng       *rand.Rand   //underlying random number generator
	scope     int          //scope of number to be generated
}

//Generating unique rand less than N
//If N is less or equal to 0, the scope will be unlimited
//If N is greater than 0, it will generate (-scope, +scope)
//If no more unique number can be generated, it will return -1 forwards
func NewUniqueRand(N int) *UniqueRand {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return &UniqueRand{
		generated: map[int]bool{},
		rng:       r1,
		scope:     N,
	}
}

func (u *UniqueRand) Int() int {
	if u.scope > 0 && len(u.generated) >= u.scope {
		return -1
	}
	for {
		var i int
		if u.scope > 0 {
			i = u.rng.Int() % u.scope
		} else {
			i = u.rng.Int()
		}
		if !u.generated[i] {
			u.generated[i] = true
			return i
		}
	}
}
