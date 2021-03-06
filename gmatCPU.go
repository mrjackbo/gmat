// +build !gpu
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

package gmat

import (
	"github.com/kuroko1t/gmat/cpu"
	"log"
)

type Tensor cpu.Tensor

func Make(shape []int) Tensor {
	tensor := Tensor{Shape: shape}
	if len(shape) == 2 {
		tensor = Tensor(cpu.Make([]int{shape[0], shape[1]}))
	} else if len(shape) == 4 {
		tensor = Tensor(cpu.Make([]int{shape[0], shape[1], shape[2], shape[3]}))
	} else if len(shape) == 6 {
		tensor = Tensor(cpu.Make([]int{shape[0], shape[1], shape[2], shape[3], shape[4], shape[5]}))
	}
	return tensor
}

func Make2DInitArray(x [][]float64) Tensor {
	z := Tensor{CPU: x}
	z.CPU = x
	z.Shape = []int{len(x), len(x[0])}
	return z
}

func Trans2D(input Tensor, n int, c int) Tensor {
	z := Tensor{}
	z.CPU = cpu.Trans2D(input.CPU, n, c)
	z.Shape = []int{n, c}
	return z
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
	result := Make([]int{reN, reC, reH, reW, reX, reY}).CPU6D
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

func Reshape4D(input Tensor, reX int, reY int) Tensor {
	z := Tensor{}
	z.CPU = cpu.Reshape4D(input.CPU4D, reX, reY)
	return z
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
	result := Make([]int{reN, reC, reH, reW, reX, reY}).CPU6D
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

func Reshape2D1D(x Tensor) []float64 {
	y := cpu.Reshape2D1D(x.CPU)
	return y
}

func Reshape1D2D(x []float64, n, c int) Tensor {
	z := Tensor{}
	z.Shape = []int{n, c}
	z.CPU = cpu.Reshape1D2D(x, n, c)
	return z
}

//func Reshape6D(input [][][][][][]float64, reX int, reY int) [][]float64 {
// 	n := len(input)
// 	c := len(input[0])
// 	h := len(input[0][0])
// 	w := len(input[0][0][0])
// 	x := len(input[0][0][0][0])
// 	y := len(input[0][0][0][0][0])
// 	if reY == -1 {
// 		reY = n * c * h * w * x * y / reX
// 	}
// 	var input1D []float64
// 	tmp := 0
// 	for i := range input {
// 		for j := range input[i] {
// 			for k := range input[i][j] {
// 				for l := range input[i][j][k] {
// 					for m := range input[i][j][l][l] {
// 						for o := range input[i][j][k][l][m] {
// 							input1D[tmp] = input[i][j][k][l][m][o]
// 							tmp++
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	result := Make2D(reX, reY)
// 	tmp = 0
// 	for i := range result {
// 		for j := range result[i] {
// 			result[i][j] = input1D[tmp]
// 			tmp++
// 		}
// 	}
// 	return result
//}

func Shape2D(x Tensor) (n int, c int) {
	n, c = cpu.Shape2D(x.CPU)
	return n, c
}

func Shape4D(input Tensor) (n int, c int, h int, w int) {
	n, c, h, w = cpu.Shape4D(input.CPU4D)
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
	z := Make([]int{zN, zC, zH, zW}).CPU4D
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

func MakeInit(n int, m int, value float64) Tensor {
	z := Tensor{}
	z.CPU = cpu.MakeInit(n, m, value)
	return z
}

func Add(x, y Tensor) Tensor {
	z := Tensor{}
	z.CPU = cpu.Add(x.CPU, y.CPU)
	z.Shape = x.Shape
	return z
}

func AddE(x Tensor, y float64) Tensor {
	z := Tensor{}
	z.CPU = cpu.AddE(x.CPU, y)
	return z
}

func Sub(x, y Tensor) Tensor {
	z := x
	z.CPU = cpu.Sub(x.CPU, y.CPU)
	return z
}

func SubE(x Tensor, y float64) Tensor {
	z := Tensor{}
	z.CPU = cpu.SubE(x.CPU, y)
	return z
}

func MulE(x Tensor, y float64) Tensor {
	z := Tensor{}
	mule := cpu.MulE(x.CPU, y)
	z.CPU = mule
	return z
}

func Mul(x, y Tensor) Tensor {
	z := Tensor{}
	z.CPU = cpu.Mul(x.CPU, y.CPU)
	return z
}

func Div(x, y Tensor) Tensor {
	z := Tensor{}
	z.CPU = cpu.Div(x.CPU, y.CPU)
	return z
}

func T(x Tensor) Tensor {
	z := Tensor{}
	z.CPU = cpu.T(x.CPU)
	return z
}

func Apply(x Tensor, fn func(float64) float64) Tensor {
	z := Tensor{}
	z.CPU = cpu.Apply(x.CPU, fn)
	return z
}

func Dot(x, y Tensor) Tensor {
	z := Tensor{}
	z.CPU = cpu.Dot(x.CPU, y.CPU)
	return z
}

func SumRow(x Tensor) Tensor {
	//sum | direction [a,b]
	//    ^           [a,b]
	z := Tensor{}
	z.CPU = cpu.SumRow(x.CPU)
	return z
}

func SumCol(x Tensor) Tensor {
	//sum -> direction [a,a]
	//				   [b,b]
	z := Tensor{}
	z.CPU = cpu.SumCol(x.CPU)
	return z
}

func Cast(x Tensor, castSize int) Tensor {
	z := Tensor{}
	z.CPU = cpu.Cast(x.CPU, castSize)
	return z
}

func MaxCol(x Tensor) Tensor {
	//sum -> direction [a,a]
	//				   [b,b]
	z := Tensor{}
	z.CPU = cpu.MaxCol(x.CPU)
	return z
}

func ArgMaxCol(x Tensor) [][]int {
	//sum -> direction [a,a]
	//				   [b,b]
	//z := Tensor{}
	maxArray := cpu.ArgMaxCol(x.CPU)
	return maxArray
}

func RandomNorm2D(r int, c int, init float64) Tensor {
	z := Tensor{}
	z.CPU = cpu.RandomNorm2D(r, c, init)
	return z
}

func HeNorm2D(r int, c int) Tensor {
	z := Tensor{}
	z.CPU = cpu.HeNorm2D(r, c)
	return z
}

func Conv1D(x, filter Tensor, stride int) Tensor {
	z := Tensor{}
	z.CPU = cpu.Conv1D(x.CPU, filter.CPU, stride)
	return z
}
