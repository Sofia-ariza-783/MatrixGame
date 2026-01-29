package models

import "fmt"

type Neo struct {
	name string
	posX int
	posY int
}

func NewNeo(initialPosX int, initialPosY int) *Neo {
	return &Neo{
		name: " N ",
		posX: initialPosX,
		posY: initialPosY,
	}
}

func (n *Neo) GetY() int {
	return n.posY
}

func (n *Neo) GetX() int {
	return n.posX
}

func (n *Neo) SetY(y int) {
	n.posY = y
}

func (n *Neo) SetX(x int) {
	n.posX = x
}

func (n *Neo) GetName() string {
	return n.name
}

func (n *Neo) IsTrapped(matrix [][]Element) bool {
	element := matrix[n.posX][n.posY]
	if _, isAgent := element.(*Agent); isAgent {
		return true
	}
	return false
}

func (n *Neo) isSafe(x, y int, matrix [][]Element) bool {
	if x < 0 || x >= len(matrix) || y < 0 || y >= len(matrix[0]) {
		return false
	}
	element := matrix[x][y]
	if _, isAgent := element.(*Agent); isAgent {
		return false
	}

	return true
}

func (n *Neo) findNearestPhone(matrix [][]Element) (int, int, bool) {
	rows := len(matrix)
	cols := len(matrix[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	queue := [][3]int{{n.posX, n.posY, 0}}
	visited[n.posX][n.posY] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		x, y, dist := current[0], current[1], current[2]

		if element := matrix[x][y]; element != nil {
			if _, isPhone := element.(*Phone); isPhone {
				return x, y, true
			}
		}

		for _, dir := range directions {
			newX, newY := x+dir[0], y+dir[1]
			if newX >= 0 && newX < rows && newY >= 0 && newY < cols && !visited[newX][newY] && n.isSafe(newX, newY, matrix) {
				visited[newX][newY] = true
				queue = append(queue, [3]int{newX, newY, dist + 1})
			}
		}
	}
	return -1, -1, false
}

func (n *Neo) CalculateMove(matrix [][]Element) (int, int) {
	phoneX, phoneY, found := n.findNearestPhone(matrix)
	if !found {
		fmt.Println("Neo decidio quedarse en el mismo lugar")
		return n.posX, n.posY
	}

	dx := phoneX - n.posX
	dy := phoneY - n.posY

	newX, newY := n.posX, n.posY

	if dx != 0 {
		if dx > 0 && n.isSafe(n.posX+1, n.posY, matrix) {
			newX = n.posX + 1
		} else if dx < 0 && n.isSafe(n.posX-1, n.posY, matrix) {
			newX = n.posX - 1
		}
	} else if dy != 0 {
		if dy > 0 && n.isSafe(n.posX, n.posY+1, matrix) {
			newY = n.posY + 1
		} else if dy < 0 && n.isSafe(n.posX, n.posY-1, matrix) {
			newY = n.posY - 1
		}
	}

	fmt.Printf("Neo decidio moverse n las coordenadas: %d %d\n", newX+1, newY+1)
	return newX, newY
}
