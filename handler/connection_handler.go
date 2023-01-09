package handler

import (
	"github.com/ribincao/ribin-game-server/memory"
	"github.com/ribincao/ribin-game-server/network"
)

func OnClose(conn *network.WrapConnection) {
	roomId := memory.GetRoomIdByPlayerId(conn.PlayerId)
	if roomId == "" {
		return
	}
	// TODO: room kill player
	// TODO: Broadcast
}
