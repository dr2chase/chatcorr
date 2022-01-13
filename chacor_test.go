// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chatcorr_test

import (
	"fmt"
	"github.com/dr2chase/chatcorr"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"
)

type FP = chatcorr.Point[float64,float64]

func line(n int, dx float64) []FP {
	var fps []FP
	for i := 0.0; i < float64(n); i += 1.0 {
		fps = append(fps, FP{dx*i, i})
	}
	return fps
}

func fuzz(fps []FP, radius float64, r *rand.Rand) {
	for i := range fps {
		delx, dely := radius*(2*r.Float64()-1), radius*(2*r.Float64()-1)
		fps[i].X += delx
		fps[i].Y += dely
	}
}

func sin(fps []FP, n, ampl float64) {
	flen := float64(len(fps))
	for i := range fps {
		theta := n*2*math.Pi*(float64(i)/flen)
		dely := ampl*math.Sin(theta)
		fps[i].Y += dely
	}
}

func stepX(fps []FP) {
	for i := 3; i < len(fps); i += 4 {
		x := fps[i-3].X
		fps[i].X, fps[i-1].X, fps[i-2].X = x, x, x
	}
}

func stepY(fps []FP) {
	for i := 3; i < len(fps); i += 4 {
		y := fps[i-3].Y
		fps[i].Y, fps[i-1].Y, fps[i-2].Y = y, y, y
	}
}

func fcmp(x, y float64) int {
	if x < y {
		return -1
	}
	if x == y {
		return 0
	}
	return 1
}

func scmp(x, y string) int {
	if x < y {
		return -1
	}
	if x == y {
		return 0
	}
	return 1
}

func icmp(x, y int) int {
	if x < y {
		return -1
	}
	if x == y {
		return 0
	}
	return 1
}



func TestLineF64(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := line(20, 0.1)
	x0 := chatcorr.CCF64(l)
	x1 := chatcorr.CC(l)
	x2 := chatcorr.CCFn(l, fcmp)
	x3 := chatcorr.CCMixed(l, fcmp, fcmp)
	r.Shuffle(len(l), func(i, j int) { l[i], l[j] = l[j], l[i] })
	x4 := chatcorr.CCF64(l)
	fmt.Printf("x0, x1, x2, x3, x4 = %f, %f, %f, %f, %f\n", x0, x1, x2, x3, x4)
	if x0 != x1 {
		t.Fail()
	}
	if x0 != x2 {
		t.Fail()
	}
	if x0 != x3 {
		t.Fail()
	}
	if x0 != x4 {
		t.Fail()
	}
}

func TestFuzzLineF64(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := line(20, 4)
	x0 := chatcorr.CCF64(l)
	fuzz(l, 2, r) // Fuzz w/ dx < 1/2 the step size, else X might overlap and change results, but larger than 1.
	x1 := chatcorr.CC(l)
	r.Shuffle(len(l), func(i, j int) { l[i], l[j] = l[j], l[i] })
	x2 := chatcorr.CCF64(l)
	fmt.Printf("x0, x1, x2 = %f, %f, %f\n", x0, x1, x2)
	if x0 <= x1 {
		t.Fail()
	}
	if x1 != x2 {
		t.Fail()
	}
}

func TestSinLineF64(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := line(20, 1.0)
	sin(l, 2, 3)
	x0 := chatcorr.CCF64(l)
	x1 := chatcorr.CC(l)
	r.Shuffle(len(l), func(i, j int) { l[i], l[j] = l[j], l[i] })
	x2 := chatcorr.CCF64(l)
	fmt.Printf("x0, x1, x2 = %f, %f, %f\n", x0, x1, x2)
	if x0 != x1 {
		t.Fail()
	}
	if x0 != x2 {
		t.Fail()
	}
}

func TestStepLineF64(t *testing.T) {
	l := line(1000, 1.0)
	x0 := chatcorr.CCF64(l)
	fmt.Printf("x0 = %f\n", x0)
	stepX(l)
	var x []float64
	const N=100
	for i := 0; i <= N; i++ {
	  x = append(x, chatcorr.CCF64(l))
	}
	sort.Slice(x, func(i,j int) bool { return x[i] < x[j]})
	fmt.Printf("x[0,25,50,75,100] = %f, %f, %f, %f, %f\n", x[0], x[N/4], x[N/2], x[3*N/4], x[N])
}

func TestFuzzStepLineF64(t *testing.T) {
	l := line(1000, 1.0)
	x0 := chatcorr.CCF64(l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fuzz(l, 4, r)
	fuzz(l, 4, r)
	fuzz(l, 4, r)
	fuzz(l, 4, r)
	x1 := chatcorr.CCF64(l)
	fmt.Printf("x0, x1 = %f, %f\n", x0, x1)
	stepX(l) // because L is sorted this increases the correlation
	var x []float64
	const N=100
	for i := 0; i <= N; i++ {
	  x = append(x, chatcorr.CCF64(l))
	}
	sort.Slice(x, func(i,j int) bool { return x[i] < x[j]})
	fmt.Printf("x[min, 25%%, 50%%, 75%%, max] = %f, %f, %f, %f, %f\n", x[0], x[N/4], x[N/2], x[3*N/4], x[N])
	if x0 <= x1 {
		t.Fail()
	}
	if x1 >= x[N/2] {
		t.Fail()
	}
	if x[0] == x[N] {
		t.Fail()
	}
}

func TestCorners(t *testing.T) {
	l := line(3, 1)
	x0 := chatcorr.CCF64(l[:0])
	x1 := chatcorr.CC(l[:1])
	x2 := chatcorr.CCFn(l[:2], fcmp)
	x3 := chatcorr.CCMixed(l[:3], fcmp, fcmp)
	fmt.Printf("x0, x1, x2, x3 = %f, %f, %f, %f\n", x0, x1, x2, x3)
	// Don't know right answers for sure for x0, x1, x2, but x3 is darn sure larger than zero
	// currently observed values are NaN, NaN, 0, 0.25
	if x2 == 0 && x3 > 0 { // writing it like this ensures x2, x3 are not NaN
		return
	}
	t.Fail()
}

func TestMixed(t *testing.T) {
	l := []chatcorr.Point[string, int]{{"ant", 1},{"bat", 2},{"cat", 3},{"dog", 4}}
	x0 := chatcorr.CC(l)
	x1 := chatcorr.CCMixed(l, scmp, icmp)
	fmt.Printf("x0, x1 = %f, %f\n", x0, x1)
	if x0 != x1 {
		t.Fail()
	}
}

func TestRepeatable(t *testing.T) {
	l := line(1000, 1.0)
	x0 := chatcorr.CCF64(l)
	fmt.Printf("x0 = %f\n", x0)
	stepX(l)
	var x []float64
	const N=100
	now := time.Now().UnixNano()
	for i := 0; i <= N; i++ {
	  x = append(x, chatcorr.CCF64Rand(l, rand.New(rand.NewSource(now))))
	}
	sort.Slice(x, func(i,j int) bool { return x[i] < x[j]})
	fmt.Printf("x[0,25,50,75,100] = %f, %f, %f, %f, %f\n", x[0], x[N/4], x[N/2], x[3*N/4], x[N])
	if x[0] != x[N] {
		t.Fail()
	}

}
