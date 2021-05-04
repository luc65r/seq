package main

import (
	"flag"
	"fmt"
	"log"

	. "ransan.tk/seq/alsa/seq"
)

const CLIENT_NAME = "test seq"

var listPorts = flag.Bool("list", false, "")
var portName = flag.String("port", "", "")

func printPorts(seq *Seq) {
	cinfo, err := ClientInfoMalloc()
	if err != nil {
		log.Fatal(err)
	}
	defer cinfo.Free()

	pinfo, err := PortInfoMalloc()
	if err != nil {
		log.Fatal(err)
	}
	defer pinfo.Free()

	fmt.Println(" Port   Client name                       Port name")
	cinfo.SetClient(-1)
	for {
		_, err := seq.QueryNextClient(cinfo)
		if err != nil {
			break
		}
		client, err := cinfo.Client()
		if err != nil {
			log.Fatal(err)
		}
		pinfo.SetClient(client)

		pinfo.SetPort(-1)
		for {
			port, err := seq.QueryNextPort(pinfo)
			if err != nil {
				break
			}
			fmt.Printf("%3v:%-3v %-33.32v %v\n",
				client, port,
				cinfo.Name(), pinfo.Name(),
			)
		}
	}
}

func recieveEvents(seq *Seq, c chan<- Event) {
	for {
		ev, err := seq.EventInput()
		if err != nil {
			break
		}
		c <- *ev
	}
	close(c)
}

func dumpEvent(ev *Event) {
	d := ev.Data()
	switch v := d.(type) {
	case Note:
		fmt.Printf("Note: channel %v, note %v, velocity %v\n",
			v.Channel(), v.Note(), v.Velocity())

	default:
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	seq, err := Open("default", OPEN_DUPLEX, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer seq.Close()

	if *listPorts {
		printPorts(seq)
	}

	port, err := seq.CreateSimplePort(
		CLIENT_NAME,
		CAP_WRITE|CAP_SUBS_WRITE,
		TYPE_MIDI_GENERIC|TYPE_APPLICATION,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer seq.DeleteSimplePort(port)

	if *portName == "" {
		return
	}

	portConnect, err := seq.ParseAddress(*portName)
	if err != nil {
		log.Fatal(err)
	}

	err = seq.ConnectFrom(0, portConnect)
	if err != nil {
		log.Fatal(err)
	}
	defer seq.DisconnectFrom(0, portConnect)

	c := make(chan Event)
	go recieveEvents(seq, c)

	for ev := range c {
		dumpEvent(&ev)
	}
}
