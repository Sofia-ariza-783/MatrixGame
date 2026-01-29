package models

import "fmt"

type Agent struct {
	name string
	posX int
	posY int
}

func NewAgent(initialPosX int, initialPosY int, id string) *Agent {
	return &Agent{
		name: id,
		posX: initialPosX,
		posY: initialPosY,
	}
}

func (a Agent) GetY() int {
	return a.posY
}

func (a Agent) GetX() int {
	return a.posX
}

func (a *Agent) SetY(y int) {
	a.posY = y
}

func (a *Agent) SetX(x int) {
	a.posX = x
}

func (a Agent) GetName() string {
	return a.name
}

func (a *Agent) findNeo(matrix [][]Element) (int, int, bool) {
	rows := len(matrix)
	cols := len(matrix[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	queue := [][3]int{{a.posX, a.posY, 0}}
	visited[a.posX][a.posY] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		x, y, dist := current[0], current[1], current[2]

		if element := matrix[x][y]; element != nil {
			if _, isNeo := element.(*Neo); isNeo {
				return x, y, true
			}
		}

		for _, dir := range directions {
			newX, newY := x+dir[0], y+dir[1]
			if newX >= 0 && newX < rows && newY >= 0 && newY < cols && !visited[newX][newY] && a.isSafe(newX, newY, matrix) {
				visited[newX][newY] = true
				queue = append(queue, [3]int{newX, newY, dist + 1})
			}
		}
	}
	return -1, -1, false
}

func (a *Agent) isSafe(x, y int, matrix [][]Element) bool {
	if x < 0 || x >= len(matrix) || y < 0 || y >= len(matrix[0]) {
		return false
	}
	if element := matrix[x][y]; element != nil {
		if _, isAgent := element.(*Agent); isAgent {
			return false
		}
		if _, isPhone := element.(*Phone); isPhone {
			return false
		}
	}
	return true
}

func (a *Agent) CalculateMove(matrix [][]Element) (int, int) {
	neoX, neoY, found := a.findNeo(matrix)
	if !found {
		fmt.Printf("El agente %s decidio quedarse en el mismo lugar", a.name)
		return a.posX, a.posY
	}

	dx := neoX - a.posX
	dy := neoY - a.posY

	newX, newY := a.posX, a.posY
	if dx != 0 {
		if dx > 0 && a.isSafe(a.posX+1, a.posY, matrix) {
			newX = a.posX + 1
		} else if dx < 0 && a.isSafe(a.posX-1, a.posY, matrix) {
			newX = a.posX - 1
		}
	} else if dy != 0 {
		if dy > 0 && a.isSafe(a.posX, a.posY+1, matrix) {
			newY = a.posY + 1
		} else if dy < 0 && a.isSafe(a.posX, a.posY-1, matrix) {
			newY = a.posY - 1
		}
	}

	fmt.Printf("El agente %s decidio moverse a las coordenadas: %d %d\n", a.name, newX+1, newY+1)
	return newX, newY
}
