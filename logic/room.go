package logic

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ribincao/ribin-game-logic/constant"
	"github.com/ribincao/ribin-game-server/codec"
	"github.com/ribincao/ribin-game-server/network"
	"github.com/ribincao/ribin-game-server/timer"
	"github.com/ribincao/ribin-game-server/utils"
	"github.com/ribincao/ribin-protocol/base"
	"google.golang.org/protobuf/proto"
)

type NormalRoom struct {
	Id        string
	Type      string
	IsExist   bool
	RoomInfo  *base.RoomInfo
	TimeWheel *timer.TimeWheel
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

func (r *NormalRoom) Run() {
	r.IsExist = false
	utils.GoWithRecover(func() {
		r.HealthCheck()
	})
	r.TimeWheel.Start()
}

func (r *NormalRoom) HealthCheck() {
	ticker := time.NewTicker(constant.HEATH_CHECK_DURATION)
	defer ticker.Stop()

	for {
		if !r.IsExist {
			return
		}
		<-ticker.C

		var deletePlayerIds []string
		for _, player := range r.GetAllPlayers() {
			playerId := player.GetId()
			lastActiveTime := player.LastActiveTime
			if time.Since(lastActiveTime).Seconds() > constant.BAD_NETWORK_TIME {
				deletePlayerIds = append(deletePlayerIds, playerId)
			}
			if player.State == constant.PLAYER_STATE_OFFLINE {
				deletePlayerIds = append(deletePlayerIds, playerId)
			}
		}
		for _, playerId := range deletePlayerIds {
			r.RemovePlayer(playerId)
		}

		if len(r.GetAllPlayers()) == 0 {
			r.Destroy()
		}

	}

}

func (r *NormalRoom) Destroy() {

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

func (r *NormalRoom) Broadcast(cmd base.Server2ClientBstType, data *base.BstBody, seq string) {
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
