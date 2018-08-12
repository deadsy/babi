//-----------------------------------------------------------------------------
/*

Babi

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------

type Babi struct {
	patch        Patch
	audio        Audio
	out_l, out_r Buf
}

func NewBabi(audio Audio) *Babi {
	return &Babi{
		audio: audio,
	}
}

func (b *Babi) AddPatch(patch Patch) {
	b.patch = patch
}

func (b *Babi) Run() {
	for {
		b.out_l.Zero()
		b.out_r.Zero()
		b.patch.Process()
		b.audio.Write(&b.out_l, &b.out_r)
	}
}

func (b *Babi) OutLR(l, r *Buf) {
	b.out_l.Add(l)
	b.out_r.Add(r)
}

//-----------------------------------------------------------------------------
