package models

type Movable interface {
	Element
	SetX(int)
	SetY(int)
	CalculateMove(matrix [][]Element) (int, int)
}
