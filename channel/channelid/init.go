/*
 *
 *  * MIT License
 *  *
 *  * Copyright (c) [2021] [xialeistudio]
 *  *
 *  * Permission is hereby granted, free of charge, to any person obtaining a copy
 *  * of this software and associated documentation files (the "Software"), to deal
 *  * in the Software without restriction, including without limitation the rights
 *  * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  * copies of the Software, and to permit persons to whom the Software is
 *  * furnished to do so, subject to the following conditions:
 *  *
 *  * The above copyright notice and this permission notice shall be included in all
 *  * copies or substantial portions of the Software.
 *  *
 *  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  * SOFTWARE.
 *
 */

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
