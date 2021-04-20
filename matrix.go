package main

import (
	"fmt"
	"math"
)

type Matrix struct {
	value      [][]float64
	answer     []float64
	cols, rows int
}

func (m *Matrix) printMatrix(text string) {

	fmt.Println(text)
	fmt.Println("-------")
	for i := 0; i < m.cols; i++ {
		str := ""
		for j := 0; j < m.rows; j++ {
			str += fmt.Sprintf("%5.2f  ", m.value[i][j])
		}
		str += "|  " + fmt.Sprintf("%5.2f  ", m.answer[i])
		fmt.Println(str)
	}
}
func (m *Matrix) maxColValue(col int) (float64, int) {
	maxIndex := col
	max := m.value[maxIndex][col]
	for j := col; j < m.rows; j++ {
		if max < math.Abs(m.value[j][col]) {
			max = math.Abs(m.value[j][col])
			maxIndex = j
		}
	}
	if math.Abs(max) < 0.0000000000001 {
		panic("Матрица вырожденна...")
	}
	return max, maxIndex
}
func (m *Matrix) rowsSwap(rowFirst, rowLast int) {
	//Меняем местами строки в матрице
	newArr := make([]float64, len(m.value[rowLast]), len(m.value[rowLast]))
	for i := 0; i < len(m.value[rowLast]); i++ {
		newArr[i] = m.value[rowLast][i]
	}
	m.value[rowLast] = m.value[rowFirst]
	m.value[rowFirst] = newArr

	//Меняем местами строки свободных членов
	newValue := m.answer[rowFirst]
	m.answer[rowFirst] = m.answer[rowLast]
	m.answer[rowLast] = newValue
}

func (m *Matrix) toSingleMatrix(col, row int) {
	deleter := m.value[row][col]

	for i := 0; i < m.cols; i++ {
		m.value[row][i] = m.value[row][i] / deleter
	}

	m.answer[row] = m.answer[row] / deleter
	//m.printMatrix("Матрица с единицей на главной диагонали: ")

	for i := 0; i < m.rows; i++ {
		if i != row {
			m.rowsDifferent(row, i, col)
		}
	}
}

func (m *Matrix) rowsDifferent(firstRow, secondRow, col int) {
	kef := m.value[secondRow][col]
	for i := 0; i < m.cols; i++ {
		m.value[secondRow][i] = m.value[secondRow][i] - (m.value[firstRow][i] * kef)
	}
	m.answer[secondRow] = m.answer[secondRow] - (m.answer[firstRow] * kef)
}

func NewMatrix(cols, rows int) *Matrix {
	return &Matrix{
		cols:   cols,
		rows:   rows,
		value:  make([][]float64, 0, 0),
		answer: make([]float64, 0, 0),
	}
}

func Execute(a [][]float64, b []float64) []float64 {
	m := &Matrix{
		value:  a,
		answer: b,
		cols:   len(a),
		rows:   len(a),
	}

	for j := 0; j < m.cols; j++ {
		_, maxRow := m.maxColValue(j)
		m.rowsSwap(j, maxRow)
		m.toSingleMatrix(j, j)
	}

	return m.answer
}
