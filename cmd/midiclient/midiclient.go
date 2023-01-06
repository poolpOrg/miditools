package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	var inPort int

	flag.IntVar(&inPort, "in", -1, "input device")
	flag.Parse()

	in := midi.GetInPorts()
	inDevices := strings.Split(in.String(), "\n")
	if len(inDevices) == 0 {
		log.Fatal("no MIDI input device found")
	}

	if len(inDevices) != 1 && inPort == -1 {
		fmt.Println(in.String())
		log.Fatal("multiple MIDI input devices found, use `-in` to select")
	}

	if inPort == -1 {
		inPort = 0
	}

	var inDev drivers.In
	var err error
	inDev, err = midi.InPort(inPort)
	if err != nil {
		log.Fatal(err)
	}

	udpServer, err := net.ResolveUDPAddr("udp", ":1053")
	if err != nil {
		println("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	//close the connection
	defer conn.Close()

	playback(conn, inDev)
}

func playback(conn *net.UDPConn, in drivers.In) {
	stop := make(chan bool, 0)
	_, err := midi.ListenTo(in, func(msg midi.Message, absms int32) {
		_, err := conn.Write(msg.Bytes())
		if err != nil {
			println("Write data failed:", err.Error())
			return
		}
	})
	if err != nil {

		log.Fatal(err)
	}
	<-stop
}
