package seq

/*
   #cgo LDFLAGS: -lasound
   #include <alsa/asoundlib.h>
*/
import "C"
import (
	"unsafe"

	. "ransan.tk/seq/alsa"
)

type Seq C.snd_seq_t

type Addr C.snd_seq_addr_t

const (
	OPEN_OUTPUT int = C.SND_SEQ_OPEN_OUTPUT
	OPEN_INPUT      = C.SND_SEQ_OPEN_INPUT
	OPEN_DUPLEX     = C.SND_SEQ_OPEN_DUPLEX
)

const NONBLOCK int = C.SND_SEQ_NONBLOCK

const (
	CAP_READ       int = C.SND_SEQ_PORT_CAP_READ
	CAP_WRITE          = C.SND_SEQ_PORT_CAP_WRITE
	CAP_SYNC_READ      = C.SND_SEQ_PORT_CAP_SYNC_READ
	CAP_SYNC_WRITE     = C.SND_SEQ_PORT_CAP_SYNC_WRITE
	CAP_DUPLEX         = C.SND_SEQ_PORT_CAP_DUPLEX
	CAP_SUBS_READ      = C.SND_SEQ_PORT_CAP_SUBS_READ
	CAP_SUBS_WRITE     = C.SND_SEQ_PORT_CAP_SUBS_WRITE
	CAP_NO_EXPORT      = C.SND_SEQ_PORT_CAP_NO_EXPORT
)

const (
	TYPE_SPECIFIC     int = C.SND_SEQ_PORT_TYPE_SPECIFIC
	TYPE_MIDI_GENERIC     = C.SND_SEQ_PORT_TYPE_MIDI_GENERIC
	TYPE_MIDI_GM          = C.SND_SEQ_PORT_TYPE_MIDI_GM
	TYPE_MIDI_GM2         = C.SND_SEQ_PORT_TYPE_MIDI_GM2
	TYPE_MIDI_GS          = C.SND_SEQ_PORT_TYPE_MIDI_GS
	TYPE_MIDI_XG          = C.SND_SEQ_PORT_TYPE_MIDI_XG
	TYPE_MIDI_MT32        = C.SND_SEQ_PORT_TYPE_MIDI_MT32
	TYPE_HARDWARE         = C.SND_SEQ_PORT_TYPE_HARDWARE
	TYPE_SOFTWARE         = C.SND_SEQ_PORT_TYPE_SOFTWARE
	TYPE_SYNTHESIZER      = C.SND_SEQ_PORT_TYPE_SYNTHESIZER
	TYPE_PORT             = C.SND_SEQ_PORT_TYPE_PORT
	TYPE_APPLICATION      = C.SND_SEQ_PORT_TYPE_APPLICATION
)

func Open(name string, streams, mode int) (seq *Seq, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cerr := C.snd_seq_open(
		(**C.snd_seq_t)(unsafe.Pointer(&seq)),
		cname, C.int(streams), C.int(mode),
	)
	if cerr < 0 {
		err = AlsaError(cerr)
	}

	return
}

func (seq *Seq) Close() error {
	cerr := C.snd_seq_close(seq.toC())
	if cerr < 0 {
		return AlsaError(cerr)
	}

	return nil
}

func (seq *Seq) CreateSimplePort(name string, caps, pType uint) (port int, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cerr := C.snd_seq_create_simple_port(
		seq.toC(), cname, C.uint(caps), C.uint(pType),
	)
	if cerr < 0 {
		err = AlsaError(cerr)
	} else {
		port = int(cerr)
	}

	return
}

func (seq *Seq) DeleteSimplePort(port int) error {
	cerr := C.snd_seq_delete_simple_port(seq.toC(), C.int(port))
	if cerr < 0 {
		return AlsaError(cerr)
	}

	return nil
}

func (seq *Seq) SetClientName(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cerr := C.snd_seq_set_client_name(seq.toC(), cname)
	if cerr < 0 {
		return AlsaError(cerr)
	}

	return nil
}

func (seq *Seq) ParseAddress(s string) (addr Addr, err error) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	cerr := C.snd_seq_parse_address(seq.toC(), addr.toC(), cstr)
	if cerr < 0 {
		err = AlsaError(cerr)
	}

	return
}

func (seq *Seq) ConnectFrom(port int, addr Addr) error {
	cerr := C.snd_seq_connect_from(
		seq.toC(), C.int(port),
		C.int(addr.client), C.int(addr.port),
	)
	if cerr < 0 {
		return AlsaError(cerr)
	}

	return nil
}

func (seq *Seq) DisconnectFrom(port int, addr Addr) error {
	cerr := C.snd_seq_disconnect_from(
		seq.toC(), C.int(port),
		C.int(addr.client), C.int(addr.port),
	)
	if cerr < 0 {
		return AlsaError(cerr)
	}

	return nil
}

func (seq *Seq) toC() *C.snd_seq_t {
	return (*C.snd_seq_t)(seq)
}

func (addr *Addr) toC() *C.snd_seq_addr_t {
	return (*C.snd_seq_addr_t)(addr)
}
