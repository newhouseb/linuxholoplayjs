package main

import (
	"encoding/binary"
	"fmt"
	// The below is a fork of "github.com/karalabe/hid" that supports SendFeatureReport
	"github.com/gorilla/websocket"
	"github.com/newhouseb/hid"
	"net/http"
)

func readPage(device *hid.Device, index uint16, length int) ([]byte, error) {
	input := make([]byte, 68)
	binary.BigEndian.PutUint16(input[2:], index)

	output := make([]byte, 68)
	read := 0
	var err error

	// The Microchip controller has a state machine that we're somehow confusing,
	// so we want to keep requesting the same report until we get ther right one back
	// (This usually doesn't happen super often)
	for true {
		device.SendFeatureReport(input)

		read, err = device.Read(output)
		if err != nil {
			return []byte{}, err
		}
		page := binary.BigEndian.Uint16(output[2:4])
		if page == index {
			break
		}
		fmt.Println("Clearing out some garbage", page, string(output))
	}

	// Each page is 64 bytes, but read returns 64 bytes including the page
	// index, so we actually need to do a couple calls to get the whole thing
	if read < length+4 {
		_, err := device.Read(output[64:])
		if err != nil {
			return []byte{}, err
		}
	}
	return output, nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		if err != nil {
			fmt.Println("Error upgrading", err)
			return
		}

		for _, info := range hid.Enumerate(0x04d8, 0xef7e) {
			if info.Product == "HoloPlay" {
				fmt.Println("Info", info, info.Product)
				device, err := info.Open()
				defer device.Close()
				if err != nil {
					fmt.Println("Failed to open device (try running as root with sudo or changing your udev rules)", err)
					return
				}

				// Read how long of a config to read
				packet, err := readPage(device, 0, 64)
				length := int32(binary.BigEndian.Uint32(packet[4:8]))
				config := []byte{}

				// Read out each page
				page := 0
				for length > 0 {
					packet, err = readPage(device, uint16(page), int(length))
					if length > 64 {
						config = append(config, packet[4:68]...)
					} else {
						config = append(config, packet[4:length+8]...)
					}
					length -= 64
					page += 1
				}

				// Chop off the initial length
				config = config[4:]

				// Serve it to holoplay.js
				fmt.Println("Serving config", string(config))
				if err = conn.WriteMessage(websocket.TextMessage, config); err != nil {
					return
				}
			}
		}

	})

	http.ListenAndServe(":11222", nil)
}
