package handler

import (
	"github.com/ribincao/ribin-game-server/manager"
	"github.com/ribincao/ribin-game-server/network"
	"github.com/ribincao/ribin-game-server/types"
)

func OnClose(conn *network.WrapConnection) {
	roomId := manager.GetRoomIdByPlayerId(conn.PlayerId)
	if roomId == "" {
		return
	}
	// TODO: room kill player
	room := manager.GetRoom[types.Room](roomId)
	if room == nil {
		return
	}
	// TODO: Broadcast
}
