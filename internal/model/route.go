package model

type Route struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Stations []Station `json:"stations"`
}

func (r *Route) GetID() int {
	return r.Id
}

func (r *Route) GetName() string {
	return r.Name
}

func (r *Route) SetID(id int) {
	r.Id = id
}

func (r *Route) SetName(name string) {
	r.Name = name
}

func (r *Route) DBTable() string {
	return "route"
}
