package constant

import "time"

const (
	CONFIG_PATH = "./conf.yaml"
	ROOM_SERVER = "room"

	HEATH_CHECK_DURATION = time.Second * 3
	BAD_NETWORK_TIME     = 30
	FRAME_SEND_TIME      = 60
	MaxFrameCnt          = 80
)

type PLAYER_NETWORK_STATE uint8

const (
	PLAYER_STATE_ONLINE PLAYER_NETWORK_STATE = iota
	PLAYER_STATE_BACKGROUND
	PLAYER_STATE_OFFLINE
)
