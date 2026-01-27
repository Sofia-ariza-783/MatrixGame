package models

type Movable interface {
	Element
	CalculateMove(matrix [][]Element) (int, int)
}
