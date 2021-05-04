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

type (
	Event C.snd_seq_event_t
	Note  C.snd_seq_ev_note_t
)

const (
	SYSTEM            int = C.SND_SEQ_EVENT_SYSTEM
	RESULT                = C.SND_SEQ_EVENT_RESULT
	NOTE                  = C.SND_SEQ_EVENT_NOTE
	NOTEON                = C.SND_SEQ_EVENT_NOTEON
	NOTEOFF               = C.SND_SEQ_EVENT_NOTEOFF
	KEYPRESS              = C.SND_SEQ_EVENT_KEYPRESS
	CONTROLLER            = C.SND_SEQ_EVENT_CONTROLLER
	PGMCHANGE             = C.SND_SEQ_EVENT_PGMCHANGE
	CHANPRESS             = C.SND_SEQ_EVENT_CHANPRESS
	PITCHBEND             = C.SND_SEQ_EVENT_PITCHBEND
	CONTROL14             = C.SND_SEQ_EVENT_CONTROL14
	NONREGPARAM           = C.SND_SEQ_EVENT_NONREGPARAM
	REGPARAM              = C.SND_SEQ_EVENT_REGPARAM
	SONGPOS               = C.SND_SEQ_EVENT_SONGPOS
	SONGSEL               = C.SND_SEQ_EVENT_SONGSEL
	QFRAME                = C.SND_SEQ_EVENT_QFRAME
	TIMESIGN              = C.SND_SEQ_EVENT_TIMESIGN
	KEYSIGN               = C.SND_SEQ_EVENT_KEYSIGN
	START                 = C.SND_SEQ_EVENT_START
	CONTINUE              = C.SND_SEQ_EVENT_CONTINUE
	STOP                  = C.SND_SEQ_EVENT_STOP
	SETPOS_TICK           = C.SND_SEQ_EVENT_SETPOS_TICK
	SETPOS_TIME           = C.SND_SEQ_EVENT_SETPOS_TIME
	TEMPO                 = C.SND_SEQ_EVENT_TEMPO
	CLOCK                 = C.SND_SEQ_EVENT_CLOCK
	TICK                  = C.SND_SEQ_EVENT_TICK
	QUEUE_SKEW            = C.SND_SEQ_EVENT_QUEUE_SKEW
	SYNC_POS              = C.SND_SEQ_EVENT_SYNC_POS
	TUNE_REQUEST          = C.SND_SEQ_EVENT_TUNE_REQUEST
	RESET                 = C.SND_SEQ_EVENT_RESET
	SENSING               = C.SND_SEQ_EVENT_SENSING
	ECHO                  = C.SND_SEQ_EVENT_ECHO
	OSS                   = C.SND_SEQ_EVENT_OSS
	CLIENT_START          = C.SND_SEQ_EVENT_CLIENT_START
	CLIENT_EXIT           = C.SND_SEQ_EVENT_CLIENT_EXIT
	CLIENT_CHANGE         = C.SND_SEQ_EVENT_CLIENT_CHANGE
	PORT_START            = C.SND_SEQ_EVENT_PORT_START
	PORT_EXIT             = C.SND_SEQ_EVENT_PORT_EXIT
	PORT_CHANGE           = C.SND_SEQ_EVENT_PORT_CHANGE
	PORT_SUBSCRIBED       = C.SND_SEQ_EVENT_PORT_SUBSCRIBED
	PORT_UNSUBSCRIBED     = C.SND_SEQ_EVENT_PORT_UNSUBSCRIBED
	USR0                  = C.SND_SEQ_EVENT_USR0
	USR1                  = C.SND_SEQ_EVENT_USR1
	USR2                  = C.SND_SEQ_EVENT_USR2
	USR3                  = C.SND_SEQ_EVENT_USR3
	USR4                  = C.SND_SEQ_EVENT_USR4
	USR5                  = C.SND_SEQ_EVENT_USR5
	USR6                  = C.SND_SEQ_EVENT_USR6
	USR7                  = C.SND_SEQ_EVENT_USR7
	USR8                  = C.SND_SEQ_EVENT_USR8
	USR9                  = C.SND_SEQ_EVENT_USR9
	SYSEX                 = C.SND_SEQ_EVENT_SYSEX
	BOUNCE                = C.SND_SEQ_EVENT_BOUNCE
	USR_VAR0              = C.SND_SEQ_EVENT_USR_VAR0
	USR_VAR1              = C.SND_SEQ_EVENT_USR_VAR1
	USR_VAR2              = C.SND_SEQ_EVENT_USR_VAR2
	USR_VAR3              = C.SND_SEQ_EVENT_USR_VAR3
	USR_VAR4              = C.SND_SEQ_EVENT_USR_VAR4
	NONE                  = C.SND_SEQ_EVENT_NONE
)

func (seq *Seq) EventInput() (ev *Event, err error) {
	cerr := C.snd_seq_event_input(
		seq.toC(),
		(**C.snd_seq_event_t)(unsafe.Pointer(&ev)),
	)
	if cerr < 0 {
		err = AlsaError(cerr)
	}

	return
}

func (ev *Event) Data() interface{} {
	switch ev._type {
	case NOTE, NOTEON, NOTEOFF, KEYPRESS:
		return *(*Note)(unsafe.Pointer(&ev.data))

	default:
		return nil
	}
}

func (n Note) Channel() uint8 {
	return uint8(n.channel)
}

func (n Note) Note() uint8 {
	return uint8(n.note)
}

func (n Note) Velocity() uint8 {
	return uint8(n.velocity)
}

func (ev *Event) toC() *C.snd_seq_event_t {
	return (*C.snd_seq_event_t)(ev)
}
