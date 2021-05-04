package alsa

/*
   #cgo LDFLAGS: -lasound
   #include <alsa/asoundlib.h>
*/
import "C"

type AlsaError C.int

func (e AlsaError) Error() string {
	return C.GoString(C.snd_strerror(C.int(e)))
}
