package common

var PacketFactory = map[uint16]func() Packet{
	OpPing: func() Packet {
		return &PingPacket{}
	},
	OpPong: func() Packet {
		return &PongPacket{}
	},
	OpRegisterReq: func() Packet {
		return &RegisterReq{}
	},
	OpRegisterResult: func() Packet {
		return &RegisterResult{}
	},
}
