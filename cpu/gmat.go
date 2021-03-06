// Copyright 2018 kurosawa. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// =============================================================================
package cpu

import (
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync"
	//"fmt"
)

func numcpu() int {
	cpus := runtime.NumCPU()
	return cpus
}

type Tensor struct {
	CPU   [][]float64
	CPU4D [][][][]float64
	CPU6D [][][][][][]float64
	Shape []int
}

func Make(shape []int) Tensor {
	tensor := Tensor{Shape: shape}
	if len(shape) == 2 {
		tensor.CPU = make2D(shape[0], shape[1])
	} else if len(shape) == 4 {
		tensor.CPU4D = make4D(shape[0], shape[1], shape[2], shape[3])
	}
	return tensor
}

func make2D(n, m int) [][]float64 {
	z := make([][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([]float64, m)
	}
	return z
}

func make3D(n, c, h int) [][][]float64 {
	z := make([][][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([][]float64, c)
		for j := range z[i] {
			z[i][j] = make([]float64, h)
		}
	}
	return z
}

func make4D(n int, c int, h int, w int) [][][][]float64 {
	z := make([][][][]float64, n)
	for i := range z {
		z[i] = make([][][]float64, c)
		for j := range z[i] {
			z[i][j] = make([][]float64, h)
			for k := range z[i][j] {
				z[i][j][k] = make([]float64, w)
			}
		}
	}
	return z
}

func Make6D(n int, c int, h int, w int, x int, y int) [][][][][][]float64 {
	z := make([][][][][][]float64, n)
	for i := range z {
		z[i] = make([][][][][]float64, c)
		for j := range z[i] {
			z[i][j] = make([][][][]float64, h)
			for k := range z[i][j] {
				z[i][j][k] = make([][][]float64, w)
				for l := range z[i][j][k] {
					z[i][j][k][l] = make([][]float64, x)
					for m := range z[i][j][k][l] {
						z[i][j][k][l][m] = make([]float64, y)
					}
				}
			}
		}
	}
	return z
}

func Trans2D(input [][]float64, n int, c int) [][]float64 {
	if n >= 2 || c >= 2 {
		log.Fatal("need to set 2 below param")
	}
	inN := len(input)
	inC := len(input[0])
	tranmap := map[int]int{0: inN, 1: inC}
	z := make2D(tranmap[n], tranmap[c])
	for i := range z {
		for j := range z[i] {
			var amap = map[int]int{0: i, 1: j}
			z[i][j] = input[amap[n]][amap[c]]
		}
	}
	return z
}

func Trans4D(input [][][][]float64, n int, c int, h int, w int) [][][][]float64 {
	if n >= 4 || c >= 4 || h >= 4 || w >= 4 {
		log.Fatal("need to set 4 below param")
	}
	inN := len(input)
	inC := len(input[0])
	inH := len(input[0][0])
	inW := len(input[0][0][0])
	tranmap := map[int]int{0: inN, 1: inC, 2: inH, 3: inW}
	z := make4D(tranmap[n], tranmap[c], tranmap[h], tranmap[w])
	for i := range z {
		for j := range z[i] {
			for k := range z[i][j] {
				for l := range z[i][j][k] {
					var amap = map[int]int{0: i, 1: j, 2: k, 3: l}
					z[i][j][k][l] = input[amap[n]][amap[c]][amap[h]][amap[w]]
				}
			}
		}
	}
	return z
}

func Trans6D(input [][][][][][]float64, n int, c int, h int, w int, x int, y int) [][][][][][]float64 {
	if n >= 6 || c >= 6 || h >= 6 || w >= 6 || x >= 6 || y >= 6 {
		log.Fatal("need to set 6 below param")
	}
	inN := len(input)
	inC := len(input[0])
	inH := len(input[0][0])
	inW := len(input[0][0][0])
	inX := len(input[0][0][0][0])
	inY := len(input[0][0][0][0][0])
	tranmap := map[int]int{0: inN, 1: inC, 2: inH, 3: inW, 4: inX, 5: inY}
	z := Make6D(tranmap[n], tranmap[c], tranmap[h], tranmap[w], tranmap[x], tranmap[y])
	for i := range z {
		for j := range z[i] {
			for k := range z[i][j] {
				for l := range z[i][j][k] {
					for m := range z[i][j][k][l] {
						for o := range z[i][j][k][l][m] {
							var amap = map[int]int{0: i, 1: j, 2: k, 3: l, 4: m, 5: o}
							z[i][j][k][l][m][o] = input[amap[n]][amap[c]][amap[h]][amap[w]][amap[x]][amap[y]]
						}
					}
				}
			}
		}
	}
	return z
}

func Reshape2D(input [][]float64, reN int, reC int, reH int, reW int) [][][][]float64 {
	n := len(input)
	c := len(input[0])
	if reW == -1 {
		reW = n * c / (reN * reC * reH)
	}
	var input1D []float64
	tmp := 0
	for i := range input {
		for j := range input[i] {
			input1D[tmp] = input[i][j]
			tmp++
		}
	}
	result := make4D(reN, reC, reH, reW)
	tmp = 0
	for i := range result {
		for j := range result[i] {
			for k := range result[i][j] {
				for l := range result[i][j][k] {
					result[i][j][k][l] = input1D[tmp]
					tmp++
				}
			}
		}
	}
	return result
}

func Reshape2D2D(input [][]float64, reX int, reY int) [][]float64 {
	n := len(input)
	c := len(input[0])
	if reX == -1 {
		reX = n * c / (reY)
	} else if reY == -1 {
		reY = n * c / (reX)
	}
	var input1D []float64
	tmp := 0
	for i := range input {
		for j := range input[i] {
			input1D[tmp] = input[i][j]
			tmp++
		}
	}
	result := make2D(reX, reY)
	tmp = 0
	for i := range result {
		for j := range result[i] {
			result[i][j] = input1D[tmp]
			tmp++
		}
	}
	return result
}

func Reshape2D1D(input [][]float64) []float64 {
	n, c := Shape2D(input)
	input1D := make([]float64, n*c)
	//var input1D []float64
	tmp := 0
	for i := range input {
		for j := range input[i] {
			input1D[tmp] = input[i][j]
			tmp++
		}
	}
	return input1D
}

func Reshape1D2D(input []float64, n, c int) [][]float64 {
	input2D := make2D(n, c)
	if len(input) != n*c {
		log.Fatal("gmat.Reshape2D1D worng shape!!")
	}
	//var input1D []float64
	tmp := 0
	for i := range input2D {
		for j := range input2D[i] {
			input2D[i][j] = input[tmp]
			tmp++
		}
	}
	return input2D
}

func Reshape2D6D(input [][]float64, reN int, reC int, reH int, reW int, reX int, reY int) [][][][][][]float64 {
	var input1D []float64
	tmp := 0
	for i := range input {
		for j := range input[i] {
			input1D[tmp] = input[i][j]
			tmp++
		}
	}
	result := Make6D(reN, reC, reH, reW, reX, reY)
	tmp = 0
	for i := range result {
		for j := range result[i] {
			for k := range result[i][j] {
				for l := range result[i][j][k] {
					for m := range result[i][j][k][l] {
						for n := range result[i][j][k][l][m] {
							result[i][j][k][l][m][n] = input1D[tmp]
							tmp++
						}
					}
				}
			}
		}
	}
	return result
}

func Reshape4D(input [][][][]float64, reX int, reY int) [][]float64 {
	n := len(input)
	c := len(input[0])
	h := len(input[0][0])
	w := len(input[0][0][0])
	if reY == -1 {
		reY = n * c * h * w / reX
	} else if reX == -1 {
		reX = n * c * h * w / reY
	}
	var input1D []float64
	tmp := 0
	for i := range input {
		for j := range input[i] {
			for k := range input[i][j] {
				for l := range input[i][j][k] {
					input1D[tmp] = input[i][j][k][l]
					tmp++
				}
			}
		}
	}
	result := make2D(reX, reY)
	tmp = 0
	for i := range result {
		for j := range result[i] {
			result[i][j] = input1D[tmp]
			tmp++
		}
	}
	return result
}

func Reshape4D6D(input [][][][]float64, reN int, reC int, reH int, reW int, reX int, reY int) [][][][][][]float64 {
	//n, c, h, w := Shape4D(input)

	var input1D []float64
	tmp := 0
	for i := range input {
		for j := range input[i] {
			for k := range input[i][j] {
				for l := range input[i][j][k] {
					input1D[tmp] = input[i][j][k][l]
					tmp++
				}
			}
		}
	}
	result := Make6D(reN, reC, reH, reW, reX, reY)
	tmp = 0
	for i := range result {
		for j := range result[i] {
			for k := range result[i][j] {
				for l := range result[i][j][k] {
					for m := range result[i][j][k][l] {
						for n := range result[i][j][k][l][m] {
							result[i][j][k][l][m][n] = input1D[tmp]
							tmp++
						}
					}
				}
			}
		}
	}
	return result
}

func Reshape6D(input [][][][][][]float64, reX int, reY int) [][]float64 {
	n := len(input)
	c := len(input[0])
	h := len(input[0][0])
	w := len(input[0][0][0])
	x := len(input[0][0][0][0])
	y := len(input[0][0][0][0][0])
	if reY == -1 {
		reY = n * c * h * w * x * y / reX
	}
	var input1D []float64
	tmp := 0
	for i := range input {
		for j := range input[i] {
			for k := range input[i][j] {
				for l := range input[i][j][k] {
					for m := range input[i][j][l][l] {
						for o := range input[i][j][k][l][m] {
							input1D[tmp] = input[i][j][k][l][m][o]
							tmp++
						}
					}
				}
			}
		}
	}
	result := make2D(reX, reY)
	tmp = 0
	for i := range result {
		for j := range result[i] {
			result[i][j] = input1D[tmp]
			tmp++
		}
	}
	return result
}

func Shape2D(input [][]float64) (n int, c int) {
	n = len(input)
	c = len(input[0])
	return n, c
}

func Shape3D(input [][][]float64) (n, h, w int) {
	n = len(input)
	h = len(input[0])
	w = len(input[0][0])
	return n, h, w
}

func Shape4D(input [][][][]float64) (n int, c int, h int, w int) {
	n = len(input)
	c = len(input[0])
	h = len(input[0][0])
	w = len(input[0][0][0])
	return n, c, h, w
}

func Shape6D(input [][][][][][]float64) (n int, c int, h int, w int, x int, y int) {
	n = len(input)
	c = len(input[0])
	h = len(input[0][0])
	w = len(input[0][0][0])
	x = len(input[0][0][0][0])
	y = len(input[0][0][0][0][0])
	return n, c, h, w, x, y
}

func Pad4D(input [][][][]float64, pad [][]int) [][][][]float64 {
	padN := len(pad)
	padM := len(pad[0])
	n := len(input)
	c := len(input[0])
	h := len(input[0][0])
	w := len(input[0][0][0])
	if padN != 4 && padM == 2 {
		log.Fatal("incorrect padding dim!!")
	}
	zN := n + pad[0][0] + pad[0][1]
	zC := c + pad[1][0] + pad[1][1]
	zH := h + pad[2][0] + pad[2][1]
	zW := w + pad[3][0] + pad[3][1]
	z := make4D(zN, zC, zH, zW)
	for i := range z {
		for j := range z[i] {
			for k := range z[i][j] {
				for l := range z[i][j][k] {
					if (pad[0][0] <= i) && (pad[1][0] <= j) &&
						(pad[2][0] <= k) && (pad[3][0] <= l) &&
						(n+pad[0][1]-1 >= i) && (c+pad[1][1]-1 >= j) &&
						(h+pad[2][1]-1 >= k) && (w+pad[3][1]-1 >= l) {
						z[i][j][k][l] = input[i-pad[0][0]][j-pad[1][0]][k-pad[2][0]][l-pad[3][0]]
					}
				}
			}
		}
	}
	return z
}

func MakeInit(n int, m int, value float64) [][]float64 {
	z := make([][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([]float64, m)
	}
	for i, zArray := range z {
		for j, _ := range zArray {
			z[i][j] = value
		}
	}
	return z
}

func Add(x [][]float64, y [][]float64) [][]float64 {
	m, n := Shape2D(x)
	z := make2D(m, n)
	//fn := func(i int, j int, wg *sync.WaitGroup) {
	// 	z[i][j] = x[i][j] + y[i][j]
	// 	wg.Done()
	//}
	//wg := &sync.WaitGroup{}
	for i, xArray := range x {
		for j, _ := range xArray {
			//wg.Add(1)
			//go fn(i, j, wg)
			z[i][j] = x[i][j] + y[i][j]
		}
	}
	//wg.Wait()
	return z
}

func AddE(x [][]float64, y float64) [][]float64 {
	n := len(x)
	m := len(x[0])
	z := make([][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([]float64, m)

	}
	for i, xArray := range x {
		for j, _ := range xArray {
			z[i][j] = x[i][j] + y
		}
	}
	return z
}

func Sub(x [][]float64, y [][]float64) [][]float64 {
	n := len(x)
	m := len(x[0])
	z := make([][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([]float64, m)

	}
	for i, xArray := range x {
		for j, _ := range xArray {
			z[i][j] = x[i][j] - y[i][j]
		}
	}
	return z
}

func SubE(x [][]float64, y float64) [][]float64 {
	n := len(x)
	m := len(x[0])
	z := make([][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([]float64, m)

	}
	for i, xArray := range x {
		for j, _ := range xArray {
			z[i][j] = x[i][j] - y
		}
	}
	return z
}

func MulE(x [][]float64, y float64) [][]float64 {
	n := len(x)
	m := len(x[0])
	z := make([][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([]float64, m)
	}
	for i, xArray := range x {
		for j, _ := range xArray {
			z[i][j] = x[i][j] * y
		}
	}
	return z
}

func Mul(x [][]float64, y [][]float64) [][]float64 {
	n := len(x)
	m := len(x[0])
	z := make([][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([]float64, m)
	}
	for i, xArray := range x {
		for j, _ := range xArray {
			z[i][j] = x[i][j] * y[i][j]
		}
	}
	return z
}

func Div(x [][]float64, y [][]float64) [][]float64 {
	n := len(x)
	m := len(x[0])
	z := make([][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([]float64, m)
	}
	for i, xArray := range x {
		for j, _ := range xArray {
			z[i][j] = x[i][j] / y[i][j]
		}
	}
	return z
}

func T(x [][]float64) [][]float64 {
	n := len(x)
	m := len(x[0])
	z := make([][]float64, m)
	for i := 0; i < m; i++ {
		z[i] = make([]float64, n)
	}
	for i, zArray := range z {
		for j, _ := range zArray {
			z[i][j] = x[j][i]
		}
	}
	return z
}

func Apply(x [][]float64, fn func(float64) float64) [][]float64 {
	n := len(x)
	m := len(x[0])
	z := make([][]float64, n)
	for i := 0; i < n; i++ {
		z[i] = make([]float64, m)
	}
	for i, xArray := range x {
		for j, _ := range xArray {
			z[i][j] = fn(x[i][j])
		}
	}
	return z
}

func Dot(x, y [][]float64) [][]float64 {
	nx, mx := Shape2D(x)
	ny, my := Shape2D(y)
	if mx != ny {
		log.Fatal("Dot.mismatch matrix number")
	}
	z := make2D(nx, my)
	wg := &sync.WaitGroup{}
	ch := make(chan int, numcpu())
	fn := func(col int, z [][]float64, wg *sync.WaitGroup) {
		for i := 0; i < mx; i++ {
			for j := 0; j < my; j++ {
				z[col][j] += x[col][i] * y[i][j]
			}
		}
		ch <- 1
		defer func() {
			<-ch
			wg.Done()
		}()
	}
	for zcol := 0; zcol < nx; zcol++ {
		wg.Add(1)
		go fn(zcol, z, wg)
	}
	wg.Wait()
	return z
}

func SumRow(x [][]float64) [][]float64 {
	//sum | direction [a,b]
	//    ^           [a,b]
	m, n := Shape2D(x)
	sumArray := make2D(1, n)
	for j := 0; j < n; j++ {
		sumValue := 0.0
		for i := 0; i < m; i++ {
			sumValue += x[i][j]
		}
		sumArray[0][j] = sumValue
	}
	return sumArray
}

func SumCol(x [][]float64) [][]float64 {
	//sum -> direction [a,a]
	//				   [b,b]
	m, n := Shape2D(x)
	sumArray := make2D(m, 1)
	for j := 0; j < m; j++ {
		sumValue := 0.0
		for i := 0; i < n; i++ {
			sumValue += x[j][i]
		}
		sumArray[j][0] = sumValue
	}
	return sumArray
}

func Cast(x [][]float64, castSize int) [][]float64 {
	m, n := Shape2D(x)
	if (m != 1) && (n != 1) {
		log.Fatal("Cast.not support format")
	}
	if m == 1 {
		for i := 0; i < castSize; i++ {
			x = append(x, x[0])
		}
	}
	if n == 1 {
		for i := 0; i < castSize-1; i++ {
			for i := 0; i < m; i++ {
				x[i] = append(x[i], x[i][0])
			}
		}
	}
	return x
}

func MaxCol(x [][]float64) [][]float64 {
	//sum -> direction [a,a]
	//				   [b,b]
	n := len(x)
	m := len(x[0])
	maxArray := make([][]float64, n)
	for i := 0; i < n; i++ {
		maxArray[i] = make([]float64, m)
	}
	for j := 0; j < n; j++ {
		max := float64(0.0)
		for i := 0; i < m; i++ {
			if x[j][i] > max {
				max = x[j][i]
			}
		}
		for i := 0; i < m; i++ {
			maxArray[j][i] = max
		}
	}
	return maxArray
}

func ArgMaxCol(x [][]float64) [][]int {
	//sum -> direction [a,a]
	//				   [b,b]
	n := len(x)
	m := len(x[0])
	maxArray := make([][]int, n)
	for i := 0; i < n; i++ {
		maxArray[i] = make([]int, m)
	}
	index := 0
	for j := 0; j < n; j++ {
		max := float64(0.0)
		for i := 0; i < m; i++ {
			if x[j][i] > max {
				max = x[j][i]
				index = i
			}
		}
		for i := 0; i < m; i++ {
			maxArray[j][i] = index
		}
	}
	return maxArray
}

func RandomNorm2D(r int, c int, init float64) [][]float64 {
	z := make([][]float64, r)
	for i := 0; i < r; i++ {
		z[i] = make([]float64, c)
	}
	for i, zArray := range z {
		for j, _ := range zArray {
			z[i][j] = rand.NormFloat64() * init
		}
	}
	return z
}

func HeNorm2D(r int, c int) [][]float64 {
	z := make([][]float64, r)
	for i := 0; i < r; i++ {
		z[i] = make([]float64, c)
	}
	for i, zArray := range z {
		for j, _ := range zArray {
			z[i][j] = rand.NormFloat64() * (1 / math.Sqrt(float64(r)))
		}
	}
	return z
}

func Conv1D(input, kernel [][]float64, stride int) [][]float64 {
	bsize_i, n := Shape2D(input)
	bsize_k, k := Shape2D(kernel)
	if bsize_i != bsize_k {
		panic("not match batchsize conv1d!")
	}
	output := make2D(bsize_i, n)
	for b := 0; b < bsize_i; b++ {
		for i := 0; i < n; i++ {
			result := 0.0
			for j := 0; j < k; j++ {
				if i+j-1 >= 0 && i+j-1 < n {
					result += input[b][i+j-1] * kernel[b][j]
				}
			}
			output[b][i] = result
		}
	}
	return output
}
