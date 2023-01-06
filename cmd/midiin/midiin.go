package main

import (
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	in := midi.GetInPorts()
	fmt.Println(in.String())

	dev, err := midi.FindInPort("iRig Keys IO 49")
	if err != nil {
		log.Fatal(err)
	}
	analyze(dev)
}

func analyze(in drivers.In) {
	var channel uint8
	var key uint8
	var velocity uint8

	stop := make(chan bool, 0)
	_, err := midi.ListenTo(in, func(msg midi.Message, absms int32) {
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
