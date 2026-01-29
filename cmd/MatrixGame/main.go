package main

import (
	"Matrix/internal/models"
	"fmt"
	"strings"
	"sync"
)

var (
	neoPosX, neoPosY, nPhones, nAgents, n, m int
	Matrix                                   [][]models.Element
	Agents                                   []models.Movable
	Neo                                      *models.Neo
	isOnGame                                 bool
	roundBarrier                             *sync.WaitGroup
)

type MoveResult struct {
	OldX, OldY int
	NewX, NewY int
	Movable    models.Movable
}

func main() {
	getMatrixSize()
	getNeo()
	getPhones()
	getAgents()

	drawGameBoard()
	isOnGame = true
	play()
}

func getMatrixSize() {
	fmt.Print("► Ingresa el tamaño de la matriz (nxm): ")
	fmt.Scan(&n, &m)

	Matrix = make([][]models.Element, n)
	for i := range Matrix {
		Matrix[i] = make([]models.Element, m)
	}
}

func getNeo() {
	fmt.Print("► Ingresa la posicion de Neo: ")
	fmt.Scan(&neoPosX, &neoPosY)

	if !verifyPosition(neoPosX-1, neoPosY-1) {
		fmt.Println("Posición inválida para Neo")
		return
	}

	Neo = models.NewNeo(neoPosX-1, neoPosY-1)
	Matrix[neoPosX-1][neoPosY-1] = Neo
}

func getPhones() {
	fmt.Print("► Ingresa cuantos telefonos hay: ")
	fmt.Scan(&nPhones)

	for i := 0; i < nPhones; i++ {
		var phonePosX, phonePosY int
		fmt.Printf("► Ingresa la posición del teléfono %d: ", i+1)
		fmt.Scan(&phonePosX, &phonePosY)

		if !verifyPosition(phonePosX-1, phonePosY-1) {
			fmt.Println("Posición inválida para teléfono")
			i--
			continue
		}

		if Matrix[phonePosX-1][phonePosY-1] != nil {
			fmt.Println("Error: ya existe un elemento en esta posición")
			i--
			continue
		}

		Phone := models.NewPhone(phonePosX-1, phonePosY-1, " ☎ ")
		Matrix[phonePosX-1][phonePosY-1] = Phone
	}
}

func getAgents() {
	fmt.Print("► Ingresa cuantos agentes hay: ")
	fmt.Scan(&nAgents)
	Agents = make([]models.Movable, 0, nAgents)

	for i := 0; i < nAgents; i++ {
		var agentPosX, agentPosY int
		fmt.Printf("► Ingresa la posición del agente %d: ", i+1)
		fmt.Scan(&agentPosX, &agentPosY)

		if !verifyPosition(agentPosX-1, agentPosY-1) {
			fmt.Println("Posición inválida para agente")
			i--
			continue
		}

		if Matrix[agentPosX-1][agentPosY-1] != nil {
			fmt.Println("Error: ya existe un elemento en esta posición")
			i--
			continue
		}

		agent := models.NewAgent(agentPosX-1, agentPosY-1, " A ")
		Matrix[agentPosX-1][agentPosY-1] = agent
		Agents = append(Agents, agent)
	}
}

func verifyPosition(x, y int) bool {
	return x >= 0 && x < n && y >= 0 && y < m
}

func play() {
	roundNumber := 1
	fmt.Println("\n\n╭┈┈┈┈┈┈┈┈┈┈┈┈┈┈ JUEGO INICIADO ┈┈┈┈┈┈┈┈┈┈┈┈┈┈╮")
	for isOnGame {
		fmt.Printf("\n╰┈➤ ❝ [Ronda %d] ❞\n\n", roundNumber)
		roundNumber++

		moves := executeRound()
		applyMoves(moves)

		drawGameBoard()
		fmt.Print("Enter para continuar")
		fmt.Scanln()
	}
	fmt.Printf("\n\n╰┈┈┈┈┈┈┈┈┈┈┈┈┈┈ JUEGO TERMINADO ┈┈┈┈┈┈┈┈┈┈┈┈┈┈╯")
}

func executeRound() []MoveResult {
	roundBarrier = &sync.WaitGroup{}
	totalMovers := len(Agents) + 1

	roundBarrier.Add(totalMovers)

	moveResults := make(chan MoveResult, totalMovers)

	for i := 0; i < len(Agents); i++ {
		go func(agent models.Movable, idx int) {
			defer roundBarrier.Done()
			newX, newY := agent.CalculateMove(Matrix)

			moveResults <- MoveResult{
				OldX:    agent.GetX(),
				OldY:    agent.GetY(),
				NewX:    newX,
				NewY:    newY,
				Movable: agent,
			}

		}(Agents[i], i)
	}

	go func() {
		defer roundBarrier.Done()
		newX, newY := Neo.CalculateMove(Matrix)

		moveResults <- MoveResult{
			OldX:    Neo.GetX(),
			OldY:    Neo.GetY(),
			NewX:    newX,
			NewY:    newY,
			Movable: Neo,
		}
	}()

	roundBarrier.Wait()
	close(moveResults)

	var moves []MoveResult
	for move := range moveResults {
		moves = append(moves, move)
	}

	return moves
}

func applyMoves(moves []MoveResult) {

	var neoNewX, neoNewY int
	agentPositions := make(map[string][2]int)

	for _, move := range moves {
		if _, isNeo := move.Movable.(*models.Neo); isNeo {
			neoNewX, neoNewY = move.NewX, move.NewY
			move.Movable.SetX(neoNewX)
			move.Movable.SetY(neoNewY)
		} else if _, isAgent := move.Movable.(*models.Agent); isAgent {
			agentPositions[move.Movable.GetName()] = [2]int{move.NewX, move.NewY}
			move.Movable.SetX(move.NewX)
			move.Movable.SetY(move.NewY)
		}
	}

	if neoNewX >= 0 && neoNewX < n && neoNewY >= 0 && neoNewY < m {
		if element := Matrix[neoNewX][neoNewY]; element != nil {
			if _, isPhone := element.(*models.Phone); isPhone {
				fmt.Printf("\n------- Neo Escapo de la Matrix -------")
				isOnGame = false
				return
			}
		}
	}

	capturedNeo := false
	for _, pos := range agentPositions {
		if pos[0] == neoNewX && pos[1] == neoNewY {
			capturedNeo = true
			break
		}
	}

	for _, move := range moves {
		if move.OldX >= 0 && move.OldX < n && move.OldY >= 0 && move.OldY < m {
			Matrix[move.OldX][move.OldY] = nil
		}
	}

	for _, move := range moves {
		if move.NewX >= 0 && move.NewX < n && move.NewY >= 0 && move.NewY < m {
			Matrix[move.NewX][move.NewY] = move.Movable
		}
	}

	if capturedNeo || Neo.IsTrapped(Matrix) {
		fmt.Println("------- Neo fue atrapado -------")
		isOnGame = false
	}
}

func drawGameBoard() {
	fmt.Printf("\n╭" + strings.Repeat(("────"), len(Matrix)) + "╮\n")
	for _, row := range Matrix {
		for _, cell := range row {
			if cell == nil {
				fmt.Print(" ▢ ")
			} else {
				fmt.Print(cell.GetName())
			}
		}
		fmt.Println("")
	}
	fmt.Printf("╰" + strings.Repeat(("────"), len(Matrix)) + "╯\n")
}
