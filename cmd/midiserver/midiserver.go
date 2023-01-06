package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	var outPort int

	flag.IntVar(&outPort, "out", -1, "output device")
	flag.Parse()

	out := midi.GetOutPorts()
	outDevices := strings.Split(strings.Trim(out.String(), "\n"), "\n")
	if len(outDevices) == 0 {
		log.Fatal("no MIDI output device found")
	}

	if len(outDevices) != 1 && outPort == -1 {
		fmt.Println(out.String())
		log.Fatal("multiple MIDI output devices found, use `-out` to select")
	}

	if outPort == -1 {
		outPort = 0
	}

	outDev, err := midi.OutPort(outPort)
	if err != nil {
		log.Fatal(err)
	}

	pc, err := net.ListenPacket("udp", "0.0.0.0:1053")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	fmt.Println("+ redirecting playback to", outDev)
	playback(pc, outDev)
}

func playback(pc net.PacketConn, out drivers.Out) {
	var channel uint8
	var key uint8
	var velocity uint8

	for {
		buf := make([]byte, 1024)
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}

		var msg midi.Message = buf[:n]

		switch msg.Type() {
		case midi.NoteOnMsg:
			msg.GetNoteOn(&channel, &key, &velocity)
			if velocity == 0 {
				fmt.Println("NOTE OFF", key)
			} else {
				fmt.Println("NOTE ON", key)
			}
			out.Send(buf[:n])

		case midi.NoteOffMsg:
			fmt.Println("NOTE OFF")
			out.Send(buf[:n])

		default:
			fmt.Println(msg)
		}
	}
	/*
		stop := make(chan bool, 0)
		_, err := midi.ListenTo(in, func(msg midi.Message, absms int32) {
			out.Send(msg.Bytes())
			switch msg.Type() {
			case midi.NoteOnMsg:
				msg.GetNoteOn(&channel, &key, &velocity)
				if velocity == 0 {
					fmt.Println("NOTE OFF", key)
				} else {
					fmt.Println("NOTE ON", key)
				}

			case midi.NoteOffMsg:
				fmt.Println("NOTE OFF")
			default:
				fmt.Println(msg)
			}
		})
		if err != nil {

			log.Fatal(err)
		}
		<-stop
	*/
	//	}
}
