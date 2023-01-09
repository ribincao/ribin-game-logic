package logic

type NormalPlayer struct {
	Id   string
	Name string
}

func (p *NormalPlayer) GetId() string {
	return p.Id
}

func (p *NormalPlayer) GetName() string {
	return p.Name
}
