//-----------------------------------------------------------------------------
/*

Audio

*/
//-----------------------------------------------------------------------------

package core

import "fmt"

//-----------------------------------------------------------------------------

var lbuf, rbuf Buf

func AudioOutLR(l, r *Buf) {
	lbuf.Add(l)
	rbuf.Add(r)
}

func AudioOutL(l *Buf) {
	lbuf.Add(l)
}

func AudioOutR(r *Buf) {
	rbuf.Add(r)
}

func AudioClear() {
	lbuf.Zero()
	rbuf.Zero()
}

func AudioDump() {
	fmt.Printf("lbuf %s\n", &lbuf)
	fmt.Printf("rbuf %s\n", &rbuf)
}

//-----------------------------------------------------------------------------
