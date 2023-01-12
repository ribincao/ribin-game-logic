package handler

import (
	"github.com/ribincao/ribin-game-logic/logic"
	errs "github.com/ribincao/ribin-game-server/error"
	"github.com/ribincao/ribin-game-server/logger"
	"github.com/ribincao/ribin-game-server/manager"
	"github.com/ribincao/ribin-game-server/network"
	"github.com/ribincao/ribin-protocol/base"
	"go.uber.org/zap"
)

func CreateRoom(enterRoomReq *base.ReqBody, conn *network.WrapConnection) (*base.RoomInfo, *errs.Error) {
	roomId, playerId := enterRoomReq.RoomId, enterRoomReq.PlayerId
	room := logic.NewNormalRoom(roomId)
	manager.RoomMng.AddRoom(room)
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
	var err *errs.Error
	switch req.RoomMessageReq.MsgType {
	case base.MsgType_E_MSGTYPE_CHAT:
		err = roomChat(player.Id, room)
	}
	if err != nil {
		logger.Error("HandleMessageError",
			zap.Any("MsgType", req.RoomMessageReq.MsgType),
			zap.String("PlayerId", player.Id),
			zap.String("RoomId", room.Id),
			zap.Error(err))
		return err
	}
	return nil
}

func roomChat(playerId string, room *logic.NormalRoom) *errs.Error {
	data := &base.BstBody{
		FromPlayerId: playerId,
	}
	room.Broadcast(base.Server2ClientBstType_E_PUSH_ROOM_MESSAGE, data, "")
	return nil
}
