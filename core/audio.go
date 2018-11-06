//-----------------------------------------------------------------------------
/*

Audio Interface Objects

*/
//-----------------------------------------------------------------------------

package core

import (
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Audio contains state for the audio stream.
type Audio struct {
}

// NewAudio returns an audio interface.
func NewAudio() (*Audio, error) {
	log.Info.Printf("")
	return &Audio{}, nil
}

// Close closes an audio stream.
func (a *Audio) Close() {
	log.Info.Printf("")
}

// Write writes an audio stream.
func (a *Audio) Write(l, r *Buf) {
}

//-----------------------------------------------------------------------------
