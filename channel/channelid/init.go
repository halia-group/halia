package channelid

import (
	"encoding/hex"
	"math/rand"
	"net"
	"sync"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		loadMachineId()
	})
}

func loadMachineId() {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, item := range interfaces {
			if item.Flags&net.FlagLoopback == 0 && item.Flags&net.FlagUp != 0 {
				machineId = item.HardwareAddr
				return
			}
		}
	}
	// use a random bytes as machineId
	loadDefaultMachineId()
	return
}

func loadDefaultMachineId() {
	machineId = make([]byte, 8)
	if _, err := rand.Read(machineId); err == nil {
		logger.Warnf(
			"Failed to find a usable hardware address from the network interfaces; using random bytes: %v",
			hex.EncodeToString(machineId),
		)
	}
}
