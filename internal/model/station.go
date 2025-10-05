package model

type Station struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	StopTime []string `json:"stop_time"`
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

func (r *Station) GetStopTime() []string {
	return r.StopTime
}

func (r *Station) SetStopTime(stopTime []string) {
	r.StopTime = stopTime
}

func (r *Station) DBTable() string {
	return "station"
}
