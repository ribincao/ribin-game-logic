package handler

import (
	"github.com/ribincao/ribin-game-logic/logic"
	errs "github.com/ribincao/ribin-game-server/error"
	"github.com/ribincao/ribin-game-server/network"
	"github.com/ribincao/ribin-protocol/base"
)

func CreateRoom(enterRoomReq *base.ReqBody, conn *network.WrapConnection) (*base.RoomInfo, *errs.Error) {
	roomId, playerId := enterRoomReq.RoomId, enterRoomReq.PlayerId
	room := logic.NewNormalRoom(roomId)
	return JoinRoom(room, playerId, conn)
}
func JoinRoom(room *logic.NormalRoom, playerId string, conn *network.WrapConnection) (*base.RoomInfo, *errs.Error) {
	player := logic.NewNormalPlayer(playerId, "TEST")
	player.SetRoomConn(conn)
	room.AddPlayer(player)
	data := &base.BstBody{}
	room.Broadcast(base.Server2ClientBstType_E_PUSH_ROOM_ENTER, data, "")
	return room.RoomInfo, nil
}

func HandleMessage(room *logic.NormalRoom, player *logic.NormalPlayer, req *base.ReqBody) *errs.Error {
	switch req.RoomMessageReq.MsgType {
	case base.MsgType_E_MSGTYPE_CHAT:
		data := &base.BstBody{}
		room.Broadcast(base.Server2ClientBstType_E_PUSH_ROOM_MESSAGE, data, "")
	}
	return nil
}
