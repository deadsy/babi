//-----------------------------------------------------------------------------
/*

Jack Client Object

*/
//-----------------------------------------------------------------------------

package core

import (
	"errors"

	"github.com/deadsy/babi/jack"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Jack contains state for the jack client.
type Jack struct {
	client *jack.Client
}

// NewJack returns a jack client object.
func NewJack(m Module) (*Jack, error) {
	log.Info.Printf("jack version %s", jack.GetVersionString())

	j := &Jack{}

	client, status := jack.ClientOpen("jack_test", jack.NoStartServer)
	if status != 0 {
		return nil, errors.New(status.String())
	}

	j.client = client

	return j, nil
}

// Close closes the jack client.
func (j *Jack) Close() {
	log.Info.Printf("")
	if j.client != nil {
		j.client.Close()
	}
}

// WriteAudio writes data to an audio stream.
func (j *Jack) WriteAudio(audio []Buf) {
}

//-----------------------------------------------------------------------------
