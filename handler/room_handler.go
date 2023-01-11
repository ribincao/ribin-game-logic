package handler

import (
	"context"
	"time"

	"github.com/ribincao/ribin-game-logic/logic"
	errs "github.com/ribincao/ribin-game-server/error"
	"github.com/ribincao/ribin-game-server/logger"
	"github.com/ribincao/ribin-game-server/manager"
	"github.com/ribincao/ribin-game-server/network"
	"github.com/ribincao/ribin-protocol/base"
	"go.uber.org/zap"
)

func HandleServerMessage(ctx context.Context, conn *network.WrapConnection, req *base.Client2ServerReq) (*base.Server2ClientRsp, error) {
	var (
		err *errs.Error
		rsp = &base.Server2ClientRsp{
			Seq: req.Seq,
		}
		rspBody = &base.RspBody{}
	)

	logger.Debug("HandleServerMessage IN", zap.Any("Req", req))
	switch req.Cmd {
	case base.Client2ServerReqCmd_E_CMD_HEART_BEAT:
		rspBody, err = handleHeartBeat(ctx, conn, req.Body, req.Seq)
	case base.Client2ServerReqCmd_E_CMD_ROOM_ENTER:
		rspBody, err = handleEnterRoom(ctx, conn, req.Body, req.Seq)
	case base.Client2ServerReqCmd_E_CMD_ROOM_LEAVE:
		rspBody, err = handleLeaveRoom(ctx, conn, req.Body, req.Seq)
	case base.Client2ServerReqCmd_E_CMD_ROOM_MESSAGE:
		rspBody, err = handleRoomMessage(ctx, req.Body, req.Seq)
	}

	if err != nil {
		rsp.Code = err.Code
		rsp.Msg = err.Message
	}
	rsp.Body = rspBody
	logger.Debug("HandleServerMessage OUT", zap.Any("Rsp", req))
	return rsp, err
}

func CheckReqParam(req *base.ReqBody) (*logic.NormalRoom, *logic.NormalPlayer, *errs.Error) {
	playerId := req.GetPlayerId()
	if playerId == "" {
		return nil, nil, errs.PlayerIdParamError
	}
	roomId := manager.GetRoomIdByPlayerId(playerId)
	if roomId == "" {
		return nil, nil, errs.RoomIdParamError
	}

	room := manager.GetRoom[*logic.NormalRoom](roomId)
	if room == nil {
		return nil, nil, errs.RoomUnexistError
	}

	player := room.GetPlayer(playerId)
	if player == nil {
		return room, nil, errs.PlayerNotInRoomError
	}
	return room, player, nil
}

// EnterRoom
func handleEnterRoom(ctx context.Context, conn *network.WrapConnection, enterRoomReq *base.ReqBody, seq string) (*base.RspBody, *errs.Error) {
	logger.Info("HandleEnterRoom IN", zap.Any("EnterRoomReq", enterRoomReq), zap.String("Seq", seq))
	var (
		err          *errs.Error
		enterRoomRsp = &base.RspBody{
			EnterRoomRsp: &base.EnterRoomRsp{},
		}
	)
	defer func() {
		logger.Info("HandleEnterRoom OUT", zap.Any("EnterRoomRsp", enterRoomRsp), zap.String("Seq", seq))
	}()

	manager.AddRoomToPlayerMap(enterRoomReq.RoomId, enterRoomReq.PlayerId)

	room, _, err := CheckReqParam(enterRoomReq)
	if err == errs.RoomUnexistError {
		roomInfo, err := CreateRoom(enterRoomReq, conn)
		enterRoomRsp.EnterRoomRsp.RoomInfo = roomInfo
		return enterRoomRsp, err
	}
	if err == errs.PlayerNotInRoomError {
		roomInfo, err := JoinRoom(room, enterRoomReq.PlayerId, conn)
		enterRoomRsp.EnterRoomRsp.RoomInfo = roomInfo
		return enterRoomRsp, err
	}
	return enterRoomRsp, err
}

// LeaveRoom
func handleLeaveRoom(ctx context.Context, conn *network.WrapConnection, leaveRoomReq *base.ReqBody, seq string) (*base.RspBody, *errs.Error) {
	logger.Info("HandleLeaveRoom IN", zap.Any("LeaveRoomReq", leaveRoomReq), zap.String("Seq", seq))
	var (
		err          *errs.Error
		leaveRoomRsp = &base.RspBody{
			LeaveRoomRsp: &base.LeaveRoomRsp{},
		}
	)
	defer func() {
		logger.Info("HandleLeaveRoom OUT", zap.Any("LeaveRoomRsp", leaveRoomRsp), zap.String("Seq", seq))
	}()
	room, player, err := CheckReqParam(leaveRoomReq)
	if err != nil {
		return leaveRoomRsp, err
	}
	room.RemovePlayer(player.GetId())
	data := &base.BstBody{}
	room.Broadcast(base.Server2ClientBstType_E_PUSH_ROOM_LEAVE, data, seq)
	return leaveRoomRsp, err
}

// HeartBeat
func handleHeartBeat(ctx context.Context, conn *network.WrapConnection, heartBeatReq *base.ReqBody, seq string) (*base.RspBody, *errs.Error) {
	var (
		err          *errs.Error
		heartBeatRsp = &base.RspBody{}
	)
	_, player, err := CheckReqParam(heartBeatReq)
	if err != nil {
		return heartBeatRsp, err
	}
	conn.UpdateLastActiveTime(time.Now().UnixMilli())
	player.LastActiveTime = time.Now()
	playerId := player.GetId()
	if playerId == "" {
		conn.PlayerId = playerId
	}
	player.SetRoomConn(conn)
	return heartBeatRsp, err
}

// Message
func handleRoomMessage(ctx context.Context, roomMessageReq *base.ReqBody, seq string) (*base.RspBody, *errs.Error) {
	logger.Info("HandleRoomMessage IN", zap.Any("RoomMessageReq", roomMessageReq), zap.String("Seq", seq))
	var (
		err            *errs.Error
		roomMessageRsp = &base.RspBody{
			RoomMessageRsp: &base.RoomMessageRsp{},
		}
	)
	defer func() {
		logger.Info("HandleRoomMessage OUT", zap.Any("RoomMessageRsp", roomMessageRsp), zap.String("Seq", seq))
	}()
	room, player, err := CheckReqParam(roomMessageReq)
	if err != nil {
		return roomMessageRsp, err
	}

	err = HandleMessage(room, player, roomMessageReq)

	return roomMessageRsp, err
}
