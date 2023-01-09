package logic

type NormalPlayer struct {
	Id   string
	Name string
}

func NewNormalPlayer(playerId string, name string) *NormalPlayer {
	return &NormalPlayer{
		Id:   playerId,
		Name: name,
	}
}

func (p *NormalPlayer) GetId() string {
	return p.Id
}

func (p *NormalPlayer) GetName() string {
	return p.Name
}
