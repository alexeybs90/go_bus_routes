package model

type Station struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (r *Station) GetID() int {
	return r.Id
}

func (r *Station) GetName() string {
	return r.Name
}

func (r *Station) SetID(id int) {
	r.Id = id
}

func (r *Station) SetName(name string) {
	r.Name = name
}

func (r *Station) DBTable() string {
	return "station"
}
