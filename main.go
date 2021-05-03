package main

/*
   #cgo LDFLAGS: -lasound
   #include <alsa/asoundlib.h>
*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

func seq_init() (seq *C.snd_seq_t) {
	cstr := C.CString("default")
	defer C.free(unsafe.Pointer(cstr))

	cerr := C.snd_seq_open(&seq, cstr, C.SND_SEQ_OPEN_DUPLEX, 0)
	if cerr < 0 {
		log.Fatal("Couldn't open sequencer")
	}

	name := C.CString("test")
	defer C.free(unsafe.Pointer(name))

	cerr = C.snd_seq_set_client_name(seq, name)
	if cerr < 0 {
		log.Fatal("Couldn't set client name")
	}

	return
}

func seq_deinit(seq *C.snd_seq_t) {
	C.snd_seq_close(seq)
}

func print_ports(seq *C.snd_seq_t) {
	var cinfo *C.snd_seq_client_info_t
	var pinfo *C.snd_seq_port_info_t

	cerr := C.snd_seq_client_info_malloc(&cinfo)
	if cerr < 0 {
		log.Fatal("Couldn't allocate memory for client info")
	}
	defer C.snd_seq_client_info_free(cinfo)

	cerr = C.snd_seq_port_info_malloc(&pinfo)
	if cerr < 0 {
		log.Fatal("Couldn't allocate memory for port info")
	}
	defer C.snd_seq_port_info_free(pinfo)

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

func main() {
	seq := seq_init()
	defer seq_deinit(seq)

	print_ports(seq)
}
