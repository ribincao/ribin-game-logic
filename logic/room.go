package logic

import (
	"sync"
)

type NormalRoom struct {
	Id        string
	Type      string
	playerMap sync.Map
}

func (r *NormalRoom) GetId() string {
	return r.Id
}

func (r *NormalRoom) GetPlayer(playerId string) *NormalPlayer {
	player, ok := r.playerMap.Load(playerId)
	if !ok {
		return nil
	}
	return player.(*NormalPlayer)
}

func (r *NormalRoom) GetAllPlayers() []*NormalPlayer {
	var players []*NormalPlayer
	r.playerMap.Range(func(key interface{}, value interface{}) bool {
		players = append(players, value.(*NormalPlayer))
		return true
	})
	return players
}

func (r *NormalRoom) AddPlayer(player *NormalPlayer) {
	r.playerMap.Store(player.GetId(), player)
}

func (r *NormalRoom) RemovePlayer(playerId string) {
	r.playerMap.Delete(playerId)
}
