package logic

import (
	"sync"
	"time"

	"github.com/ribincao/ribin-game-logic/constant"
	"github.com/ribincao/ribin-game-server/logger"
	"github.com/ribincao/ribin-game-server/network"
	"go.uber.org/zap"
)

type NormalPlayer struct {
	Id             string
	Name           string
	State          constant.PLAYER_NETWORK_STATE
	LastActiveTime time.Time
	RoomConn       *network.WrapConnection
	sync.RWMutex
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

func (p *NormalPlayer) GetRoomConn() *network.WrapConnection {
	return p.RoomConn
}

func (p *NormalPlayer) SetRoomConn(conn *network.WrapConnection) {
	p.Lock()
	if conn == p.RoomConn {
		p.Unlock()
		return
	}
	var oldConn *network.WrapConnection
	if conn != p.RoomConn && p.RoomConn != nil && !p.RoomConn.IsClosed.Load() {
		logger.Info("SetRoomConn",
			zap.Any("OldConnection", p.RoomConn.Connection.RemoteAddr()),
			zap.Any("NewConnection", conn.Connection.RemoteAddr()))
		oldConn = p.RoomConn
	}
	conn.PlayerId = p.Id
	p.RoomConn = conn
	p.Unlock()
	if oldConn != nil {
		oldConn.Close()
	}
	p.RoomConn.OnConnect()
}
