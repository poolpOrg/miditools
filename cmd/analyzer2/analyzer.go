package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/poolpOrg/go-harmony/intervals"
	"github.com/poolpOrg/go-harmony/notes"
	"gitlab.com/gomidi/midi/v2"

	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func analyze(absms int32, pitches []uint8) {
	fmt.Printf("@%d\t", absms)
	switch len(pitches) {
	case 0:
		fmt.Println("rest")
	case 1:
		midiNote := midi.Note(pitches[0])
		n, err := notes.Parse(fmt.Sprintf("%s%d", midiNote.Name(), midiNote.Octave()))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("note:", n.OctaveName())

	case 2:
		fmt.Println("interval:", midi.Note(pitches[0]), midi.Note(pitches[1]))
		for _, interval := range intervals.Intervals() {
			if uint8(interval.Semitones()) == pitches[1]-pitches[0] {
				fmt.Println("\t", interval.Name())
			}
		}
	default:
		fmt.Println("chord:", pitches)
	}
}

func main() {
	var opt_input string

	flag.StringVar(&opt_input, "input", "", "MIDI input device")
	flag.Parse()

	var inPort drivers.In

	inputPortID, err := strconv.Atoi(opt_input)
	if err != nil {
		inPort, err = midi.FindInPort(opt_input)
		if err == nil {
			if !inPort.IsOpen() {
				err = inPort.Open()
			}
		}
	} else {
		inPort, err = midi.InPort(inputPortID)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(inPort)

	activeNotes := make(map[uint8]bool)
	fmt.Println(activeNotes)

	_, err = midi.ListenTo(inPort, func(msg midi.Message, absms int32) {
		var channel uint8
		var key uint8
		var velocity uint8

		if msg.GetNoteOn(&channel, &key, &velocity) {
			if velocity == 0 {
				delete(activeNotes, key)
			} else {
				if _, exists := activeNotes[key]; !exists {
					activeNotes[key] = true
				}
			}
		} else if msg.GetNoteOff(&channel, &key, &velocity) {
			delete(activeNotes, key)
		}

		pitches := make([]uint8, 0)
		for pitch, _ := range activeNotes {
			pitches = append(pitches, pitch)
		}
		sort.Slice(pitches, func(i int, j int) bool { return pitches[i] < pitches[j] })
		analyze(absms, pitches)
	})
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan bool)
}
