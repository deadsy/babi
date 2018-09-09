//-----------------------------------------------------------------------------
/*

Circular Buffer

*/
//-----------------------------------------------------------------------------

package cbuf

import (
	"errors"
	"sync"
)

//-----------------------------------------------------------------------------

type CircularBuffer struct {
	lock   sync.Mutex    // access locking
	buffer []interface{} // the buffer
	rd, wr int           // read/write indices
}

//-----------------------------------------------------------------------------

// Increment and wrap-around an index value.
func inc_mod(idx, size int) int {
	idx += 1
	if idx == size {
		return 0
	}
	return idx
}

//-----------------------------------------------------------------------------

// NewCircularBuffer returns a circular buffer of size elements.
func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		buffer: make([]interface{}, size),
	}
}

//-----------------------------------------------------------------------------

// Read reads a value from the circular buffer, or returns "empty" as an error.
func (cb *CircularBuffer) Read() (interface{}, error) {
	cb.lock.Lock()
	defer cb.lock.Unlock()
	if cb.rd != cb.wr {
		val := cb.buffer[cb.rd]
		cb.rd = inc_mod(cb.rd, len(cb.buffer))
		return val, nil
	}
	return nil, errors.New("empty")
}

// Write writes a value to the circular buffer, or returns "full" as an error.
func (cb *CircularBuffer) Write(val interface{}) error {
	cb.lock.Lock()
	defer cb.lock.Unlock()
	wr_inc := inc_mod(cb.wr, len(cb.buffer))
	if wr_inc == cb.rd {
		return errors.New("full")
	}
	cb.buffer[cb.wr] = val
	cb.wr = wr_inc
	return nil
}

// Empty returns true if the circular buffer is empty.
func (cb *CircularBuffer) Empty() bool {
	cb.lock.Lock()
	defer cb.lock.Unlock()
	return cb.rd == cb.wr
}

//-----------------------------------------------------------------------------
