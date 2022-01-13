// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chatcorr

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

// chatcorr provides functions that implement the Chatterjee Correlation 
// (see https://arxiv.org/pdf/1909.10140.pdf) for a variety of types.  
// This is also an exercise in programming with Go generics.

type Lessable interface {
	~int | ~uint | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64 | ~float32 | ~float64 | ~string
}

type Point[T, U any] struct {
	X T
	Y U
}

func finish(perm, r, l []int) float64 {
	numerator := 0.0
	denominator := 0.0
	fLenV := float64(len(perm))
	for i := 0; i < len(perm)-1; i++ {
		numerator += math.Abs(float64(r[perm[i+1]]) - float64(r[perm[i]]))
	}
	numerator *= fLenV

	for i := 0; i < len(l); i++ {
		li := float64(l[i])
		denominator += li * (fLenV - li)
	}
	denominator *= 2

	return 1 - numerator/denominator
}

func makePerm(l int) []int {
	perm := make([]int, l, l)
	for i := range perm {
		perm[i] = i
	}
	return perm
}

func shuffleX(rng *rand.Rand, last_i, i int, perm []int) int {
	if last_i+1 < i {
			rng.Shuffle(i-last_i, func(i, j int) {
				perm[last_i+i], perm[last_i+j] = perm[last_i+j], perm[last_i+i]
			})
		}
	return i
}

// CCF64 sorts v by y coordinate, and returns the Chatterjee Correlation of x and y.
// See https://arxiv.org/pdf/1909.10140.pdf
func CCF64(v []Point[float64, float64]) (float64) {
	return CCF64Rand(v, rand.New(rand.NewSource(int64(time.Now().Nanosecond()))))
}

// CCF64Rand is CCF64 with an explicit rng for breaking X ties in a repeatable way.
func CCF64Rand(v []Point[float64, float64], rng *rand.Rand) (float64) {
	sort.Slice(v, func(i, j int) bool { return v[i].Y < v[j].Y })

	// r[i] = | {j: v[j].Y <= v[i].Y} |
	// l[i] = | {j: v[j].Y >= v[i].Y} |
	r, l := make([]int, len(v), len(v)), make([]int, len(v), len(v))

	last_i := 0
	recordRL := func(i int) int {
		if i < len(v) && v[last_i].Y == v[i].Y {
			return last_i
		}
		for j := last_i; j < i; j++ {
			r[j] = i               // e.g. y0, y1;  y1 != y0; r[0] = 1
			l[j] = len(v) - last_i // e.g ynm2, ynm1, len(v); ynm2 != ynm1
		}
		last_i = i
		return i
	}
	for i := 1; i < len(v); i++ {
		last_i = recordRL(i)
	}
	// handle trailing case; i == len(v)
	recordRL(len(v))

	// perm will hold the X-sorted order of the Y-sorted V
	perm := makePerm(len(v))
	sort.Slice(perm, func(i, j int) bool { return v[perm[i]].X < v[perm[j]].X })

	// randomly break ties in x
	last_i = 0
	shuffleEQX := func(i int) int {
		if i < len(v) && v[perm[last_i]].X == v[perm[i]].X {
			return last_i
		}
		return shuffleX(rng,last_i, i, perm)
	}
	for i := 1; i < len(v); i++ {
		last_i = shuffleEQX(i)
	}
	shuffleEQX(len(v))

	return finish(perm, r, l)
}

// CC sorts v by y coordinate, and returns the Chatterjee Correlation of x and y.
// See https://arxiv.org/pdf/1909.10140.pdf
func CC[T, U Lessable](v []Point[T, U]) float64 {
	return CCRand(v, rand.New(rand.NewSource(int64(time.Now().Nanosecond()))))
}
// CCRand is CC with an explicit rng for breaking X ties in a repeatable way.
func CCRand[T, U Lessable](v []Point[T, U], rng *rand.Rand) float64 {
	sort.Slice(v, func(i, j int) bool { return v[i].Y < v[j].Y })

	// r[i] = | {j: v[j].Y <= v[i].Y} |
	// l[i] = | {j: v[j].Y >= v[i].Y} |
	r, l := make([]int, len(v), len(v)), make([]int, len(v), len(v))

	last_i := 0
	recordRL := func(i int) int {
		if i < len(v) && v[last_i].Y == v[i].Y {
			return last_i
		}
		for j := last_i; j < i; j++ {
			r[j] = i               // e.g. y0, y1;  y1 != y0; r[0] = 1
			l[j] = len(v) - last_i // e.g ynm2, ynm1, len(v); ynm2 != ynm1
		}
		last_i = i
		return i
	}
	for i := 1; i < len(v); i++ {
		last_i = recordRL(i)
	}
	// handle trailing case; i == len(v)
	recordRL(len(v))

	// perm will hold the X-sorted order of the Y-sorted V
	perm := makePerm(len(v))
	sort.Slice(perm, func(i, j int) bool { return v[perm[i]].X < v[perm[j]].X })

	// randomly break ties in x
	last_i = 0
	shuffleEQX := func(i int) int {
		if i < len(v) && v[perm[last_i]].X == v[perm[i]].X {
			return last_i
		}
		return shuffleX(rng,last_i, i, perm)
	}
	for i := 1; i < len(v); i++ {
		last_i = shuffleEQX(i)
	}
	shuffleEQX(len(v))

	return finish(perm, r, l)
}

// CCFn sorts v and returns its (x, y) Chatterjee correlation
// compare returns -1, 0, 1 as a is <, =, > b
// See https://arxiv.org/pdf/1909.10140.pdf
func CCFn[T any](v []Point[T, T], compare func (a,b T) int) float64 {
		return CCFnRand(v, compare, rand.New(rand.NewSource(int64(time.Now().Nanosecond()))))
}
// CCFnRand is CCFn with an explicit rng for breaking X ties in a repeatable way.
func CCFnRand[T any](v []Point[T, T], compare func (a,b T) int, rng *rand.Rand) float64 {
	sort.Slice(v, func(i, j int) bool { return compare(v[i].Y,v[j].Y) < 0 })

	// r[i] = | {j: v[j].Y <= v[i].Y} |
	// l[i] = | {j: v[j].Y >= v[i].Y} |
	r, l := make([]int, len(v), len(v)), make([]int, len(v), len(v))

	last_i := 0
	recordRL := func(i int) int {
		if i < len(v) && compare(v[last_i].Y, v[i].Y) == 0 {
			return last_i
		}
		for j := last_i; j < i; j++ {
			r[j] = i               // e.g. y0, y1;  y1 != y0; r[0] = 1
			l[j] = len(v) - last_i // e.g ynm2, ynm1, len(v); ynm2 != ynm1
		}
		last_i = i
		return i
	}
	for i := 1; i < len(v); i++ {
		last_i = recordRL(i)
	}
	// handle trailing case; i == len(v)
	recordRL(len(v))

	// perm will hold the X-sorted order of the Y-sorted V
	perm := makePerm(len(v))
	sort.Slice(perm, func(i, j int) bool { return compare(v[perm[i]].X, v[perm[j]].X) < 0})

	// randomly break ties in x
	last_i = 0
	shuffleEQX := func(i int) int {
		if i < len(v) && compare(v[perm[last_i]].X, v[perm[i]].X) == 0 {
			return last_i
		}
		return shuffleX(rng,last_i, i, perm)
	}
	for i := 1; i < len(v); i++ {
		last_i = shuffleEQX(i)
	}
	shuffleEQX(len(v))

	return finish(perm, r, l)

}

// CCMixed sorts v and returns its (x, y) Chatterjee correlation
// compareT and compareU return -1, 0, 1 as a is <, =, > b
// See https://arxiv.org/pdf/1909.10140.pdf
func CCMixed[T,U any](v []Point[T, U], compareT func (a,b T) int, compareU func (a,b U) int) float64 {
		return CCMixedRand(v, compareT, compareU, rand.New(rand.NewSource(int64(time.Now().Nanosecond()))))
}
// CCMixedRand is CCMixed with an explicit rng for breaking X ties in a repeatable way.
func CCMixedRand[T,U any](v []Point[T, U], compareT func (a,b T) int, compareU func (a,b U) int, rng *rand.Rand) float64 {
	sort.Slice(v, func(i, j int) bool { return compareU(v[i].Y,v[j].Y) < 0 })

	// r[i] = | {j: v[j].Y <= v[i].Y} |
	// l[i] = | {j: v[j].Y >= v[i].Y} |
	r, l := make([]int, len(v), len(v)), make([]int, len(v), len(v))

	last_i := 0
	recordRL := func(i int) int {
		if i < len(v) && compareU(v[last_i].Y, v[i].Y) == 0 {
			return last_i
		}
		for j := last_i; j < i; j++ {
			r[j] = i               // e.g. y0, y1;  y1 != y0; r[0] = 1
			l[j] = len(v) - last_i // e.g ynm2, ynm1, len(v); ynm2 != ynm1
		}
		last_i = i
		return i
	}
	for i := 1; i < len(v); i++ {
		last_i = recordRL(i)
	}
	// handle trailing case; i == len(v)
	recordRL(len(v))

	// perm will hold the X-sorted order of the Y-sorted V
	perm := makePerm(len(v))
	sort.Slice(perm, func(i, j int) bool { return compareT(v[perm[i]].X, v[perm[j]].X) < 0})

	// randomly break ties in x
	last_i = 0
	shuffleEQX := func(i int) int {
		if i < len(v) && compareT(v[perm[last_i]].X, v[perm[i]].X) == 0 {
			return last_i
		}
		return shuffleX(rng,last_i, i, perm)
	}
	for i := 1; i < len(v); i++ {
		last_i = shuffleEQX(i)
	}
	shuffleEQX(len(v))

	return finish(perm, r, l)
}


