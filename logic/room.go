package logic

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/ribincao/ribin-game-server/codec"
	"github.com/ribincao/ribin-game-server/network"
	"github.com/ribincao/ribin-protocol/base"
	"google.golang.org/protobuf/proto"
)

type NormalRoom struct {
	Id        string
	Type      string
	RoomInfo  *base.RoomInfo
	playerMap sync.Map
}

func NewNormalRoom(roomId string) *NormalRoom {
	return &NormalRoom{
		Id: roomId,
	}
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

func (r *NormalRoom) Broadcast(cmd base.Server2ClientBstType, data *base.RspBody, seq string) {
	msg := &base.Server2ClientBst{
		Type: cmd,
		Body: data,
		Seq:  seq,
	}
	reqbuf, err := proto.Marshal(msg)
	if err != nil {
		return
	}

	frame, err := codec.DefaultCodec.Encode(reqbuf, codec.Broadcast)
	if err != nil {
		return
	}

	var conns []*network.WrapConnection
	players := r.GetAllPlayers()
	for _, player := range players {
		c := player.GetRoomConn()
		if c == nil || c.Connection == nil || c.IsClosed.Load() {
			continue
		}
		conns = append(conns, c)
	}

	for _, c := range conns {
		c.Write(websocket.BinaryMessage, frame)
	}
}
