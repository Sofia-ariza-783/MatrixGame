package models

type Phone struct {
	name string
	posX int
	posY int
}

func NewPhone(initialPosX int, initialPosY int, id string) *Phone {
	return &Phone{
		name: id,
		posX: initialPosX,
		posY: initialPosY,
	}
}

func (p Phone) GetY() int {
	return p.posY
}

func (p Phone) GetX() int {
	return p.posX
}

func (n Phone) GetName() string {
	return n.name
}
