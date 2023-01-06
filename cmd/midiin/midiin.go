package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	var inPort int
	var outPort int

	flag.IntVar(&inPort, "in", -1, "input device")
	flag.IntVar(&outPort, "out", -1, "output device")
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

	out := midi.GetOutPorts()
	outDevices := strings.Split(out.String(), "\n")
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

	var inDev drivers.In
	var err error
	inDev, err = midi.InPort(inPort)
	if err != nil {
		log.Fatal(err)
	}

	var outDev drivers.Out
	outDev, err = midi.OutPort(outPort)
	if err != nil {
		log.Fatal(err)
	}

	analyze(inDev, outDev)
}

func analyze(in drivers.In, out drivers.Out) {
	var channel uint8
	var key uint8
	var velocity uint8

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
}
