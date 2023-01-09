package handler

import (
	"github.com/ribincao/ribin-game-logic/logic"
	"github.com/ribincao/ribin-game-server/manager"
	"github.com/ribincao/ribin-game-server/network"
	"github.com/ribincao/ribin-protocol/base"
)

func OnClose(conn *network.WrapConnection) {
	roomId := manager.GetRoomIdByPlayerId(conn.PlayerId)
	if roomId == "" {
		return
	}
	room := manager.GetRoom[*logic.NormalRoom](roomId)
	if room == nil {
		return
	}
	room.RemovePlayer(conn.PlayerId)
	room.Broadcast(base.Server2ClientBstType_E_PUSH_ROOM_MESSAGE, nil, "") // TODO: Broadcast
}
