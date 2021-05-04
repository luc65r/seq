package main

/*
   #cgo LDFLAGS: -lasound
   #include <alsa/asoundlib.h>
*/
import "C"
import (
	"flag"
	"fmt"
	"log"
	"unsafe"
)

const CLIENT_NAME = "test seq"

var listPorts = flag.Bool("list", false, "")
var portName = flag.String("port", "", "")

func newSeq() (seq *C.snd_seq_t) {
	cstr := C.CString("default")
	defer C.free(unsafe.Pointer(cstr))

	cerr := C.snd_seq_open(&seq, cstr, C.SND_SEQ_OPEN_DUPLEX, 0)
	if cerr < 0 {
		log.Fatalln("Couldn't open sequencer")
	}

	name := C.CString(CLIENT_NAME)
	defer C.free(unsafe.Pointer(name))

	cerr = C.snd_seq_set_client_name(seq, name)
	if cerr < 0 {
		log.Fatalln("Couldn't set client name")
	}

	return
}

func (seq *C.snd_seq_t) free() {
	C.snd_seq_close(seq)
}

func (seq *C.snd_seq_t) printPorts() {
	var cinfo *C.snd_seq_client_info_t
	var pinfo *C.snd_seq_port_info_t

	cerr := C.snd_seq_client_info_malloc(&cinfo)
	if cerr < 0 {
		log.Fatalln("Couldn't allocate memory for client info")
	}
	defer C.snd_seq_client_info_free(cinfo)

	cerr = C.snd_seq_port_info_malloc(&pinfo)
	if cerr < 0 {
		log.Fatalln("Couldn't allocate memory for port info")
	}
	defer C.snd_seq_port_info_free(pinfo)

	fmt.Println(" Port   Client name                       Port name")
	C.snd_seq_client_info_set_client(cinfo, -1)
	for C.snd_seq_query_next_client(seq, cinfo) >= 0 {
		client := C.snd_seq_client_info_get_client(cinfo)
		C.snd_seq_port_info_set_client(pinfo, client)

		C.snd_seq_port_info_set_port(pinfo, -1)
		for C.snd_seq_query_next_port(seq, pinfo) >= 0 {
			fmt.Printf("%3v:%-3v %-33.32v %v\n",
				C.snd_seq_port_info_get_client(pinfo),
				C.snd_seq_port_info_get_port(pinfo),
				C.GoString(C.snd_seq_client_info_get_name(cinfo)),
				C.GoString(C.snd_seq_port_info_get_name(pinfo)),
			)
		}
	}
}

func (seq *C.snd_seq_t) initPort() (port C.int) {
	name := C.CString(CLIENT_NAME)
	defer C.free(unsafe.Pointer(name))

	port = C.snd_seq_create_simple_port(seq, name,
		C.SND_SEQ_PORT_CAP_WRITE|C.SND_SEQ_PORT_CAP_SUBS_WRITE,
		C.SND_SEQ_PORT_TYPE_MIDI_GENERIC|C.SND_SEQ_PORT_TYPE_APPLICATION,
	)
	if port < 0 {
		log.Fatalln("Couldn't create port")
	}

	return
}

func (seq *C.snd_seq_t) deinitPort(port C.int) {
	cerr := C.snd_seq_delete_simple_port(seq, port)
	if cerr < 0 {
		log.Fatalln("Couldn't delete port")
	}
}

func (seq *C.snd_seq_t) parsePort(s string) (addr C.snd_seq_addr_t) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	cerr := C.snd_seq_parse_address(seq, &addr, cstr)
	if cerr < 0 {
		log.Fatalf("Invalid port %v", s)
	}

	return
}

func (seq *C.snd_seq_t) connectPort(addr C.snd_seq_addr_t) {
	cerr := C.snd_seq_connect_from(seq, 0, C.int(addr.client), C.int(addr.port))
	if cerr < 0 {
		log.Fatalf("Couldn't connect from port %v:%v\n",
			addr.client, addr.port)
	}
}

func (seq *C.snd_seq_t) disconnectPort(addr C.snd_seq_addr_t) {
	cerr := C.snd_seq_disconnect_from(seq, 0,
		C.int(addr.client), C.int(addr.port))
	if cerr < 0 {
		log.Fatalf("Couldn't disconnect from port %v:%v\n",
			addr.client, addr.port)
	}
}

func (seq *C.snd_seq_t) recieveEvents(c chan<- C.snd_seq_event_t) {
	var event *C.snd_seq_event_t
	for C.snd_seq_event_input(seq, &event) >= 0 {
		c <- *event
	}
	close(c)
}

func main() {
	flag.Parse()

	seq := newSeq()
	defer seq.free()

	if *listPorts {
		seq.printPorts()
	}

	port := seq.initPort()
	defer seq.deinitPort(port)

	portConnect := seq.parsePort(*portName)

	seq.connectPort(portConnect)
	defer seq.disconnectPort(portConnect)

	c := make(chan C.snd_seq_event_t)
	go seq.recieveEvents(c)

	for e := range c {
		fmt.Println(e)
	}
}
