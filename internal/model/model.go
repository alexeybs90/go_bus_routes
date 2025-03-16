package model

type Model interface {
	GetID() int
	GetName() string
	SetID(id int)
	SetName(name string)
	DBTable() string
}
