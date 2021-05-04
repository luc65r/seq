package seq

/*
   #cgo LDFLAGS: -lasound
   #include <alsa/asoundlib.h>
*/
import "C"
import (
	. "ransan.tk/seq/alsa"
	"unsafe"
)

type (
	ClientInfo C.snd_seq_client_info_t
	PortInfo   C.snd_seq_port_info_t
)

func (seq *Seq) QueryNextClient(cinfo *ClientInfo) (client int, err error) {
	cerr := C.snd_seq_query_next_client(seq.toC(), cinfo.toC())
	if cerr < 0 {
		err = AlsaError(cerr)
	} else {
		client = int(cerr)
	}

	return
}

func ClientInfoMalloc() (cinfo *ClientInfo, err error) {
	cerr := C.snd_seq_client_info_malloc(
		(**C.snd_seq_client_info_t)(unsafe.Pointer(&cinfo)),
	)
	if cerr < 0 {
		err = AlsaError(cerr)
	}

	return
}

func (cinfo *ClientInfo) Free() {
	C.snd_seq_client_info_free(cinfo.toC())
}

func (cinfo *ClientInfo) SetClient(c int) {
	C.snd_seq_client_info_set_client(cinfo.toC(), C.int(c))
}

func (cinfo *ClientInfo) Client() (client int, err error) {
	cerr := C.snd_seq_client_info_get_client(cinfo.toC())
	if cerr < 0 {
		err = AlsaError(cerr)
	} else {
		client = int(cerr)
	}

	return
}

func (cinfo *ClientInfo) Name() string {
	return C.GoString(C.snd_seq_client_info_get_name(cinfo.toC()))
}

func (seq *Seq) QueryNextPort(pinfo *PortInfo) (port int, err error) {
	cerr := C.snd_seq_query_next_port(seq.toC(), pinfo.toC())
	if cerr < 0 {
		err = AlsaError(cerr)
	} else {
		port = int(cerr)
	}

	return
}

func PortInfoMalloc() (pinfo *PortInfo, err error) {
	cerr := C.snd_seq_port_info_malloc(
		(**C.snd_seq_port_info_t)(unsafe.Pointer(&pinfo)),
	)
	if cerr < 0 {
		err = AlsaError(cerr)
	}

	return
}

func (pinfo *PortInfo) Free() {
	C.snd_seq_port_info_free(pinfo.toC())
}

func (pinfo *PortInfo) SetPort(c int) {
	C.snd_seq_port_info_set_port(pinfo.toC(), C.int(c))
}

func (pinfo *PortInfo) Port() (port int, err error) {
	cerr := C.snd_seq_port_info_get_port(pinfo.toC())
	if cerr < 0 {
		err = AlsaError(cerr)
	} else {
		port = int(cerr)
	}

	return
}

func (pinfo *PortInfo) SetClient(c int) {
	C.snd_seq_port_info_set_client(pinfo.toC(), C.int(c))
}

func (pinfo *PortInfo) Client() (client int, err error) {
	cerr := C.snd_seq_port_info_get_client(pinfo.toC())
	if cerr < 0 {
		err = AlsaError(cerr)
	} else {
		client = int(cerr)
	}

	return
}

func (pinfo *PortInfo) Name() string {
	return C.GoString(C.snd_seq_port_info_get_name(pinfo.toC()))
}

func (cinfo *ClientInfo) toC() *C.snd_seq_client_info_t {
	return (*C.snd_seq_client_info_t)(cinfo)
}

func (pinfo *PortInfo) toC() *C.snd_seq_port_info_t {
	return (*C.snd_seq_port_info_t)(pinfo)
}
