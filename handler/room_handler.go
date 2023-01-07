package handler

import (
	"context"

	errs "github.com/ribincao/ribin-game-server/error"
	"github.com/ribincao/ribin-game-server/logger"
	"github.com/ribincao/ribin-game-server/network"
	"github.com/ribincao/ribin-protocol/base"
	"go.uber.org/zap"
)

func HandleRoomMessage(ctx context.Context, conn *network.WrapConnection, req *base.Test) (*base.Test, error) {
	var (
		err *errs.Error
		rsp = &base.Test{}
	)

	logger.Info("HandleRoomMessage IN", zap.Any("Req", req))

	logger.Info("HandleRoomMessage OUT", zap.Any("Rsp", req))
	return rsp, err
}
