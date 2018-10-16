#!/usr/bin/env python3
#------------------------------------------------------------------------------
"""

Generate a loopless FFT function (in golang) for a given input length.

"""
#------------------------------------------------------------------------------

_package = "core"
_fp = 64
_fft_n = 32

_fft_hn = _fft_n >> 1
_fft_hmask = _fft_hn - 1

_fft_nbits = _fft_n.bit_length() - 1
_complex_type = ("complex128", "complex64")[_fp == 32]

_buf_type = "FFT%dBuf" % _fft_n
_w_type = "FFT%dTwiddle" % _fft_n

#------------------------------------------------------------------------------

def reverse_bits(x, n):
  """reverse n bits of x"""
  rev = 0
  for i in range(n):
    rev = (rev << 1) + (x & 1)
    x >>= 1
  return rev

#------------------------------------------------------------------------------

def main():

  # declare package
  print("package %s" % _package)

  # declare types
  print("type %s [%d]%s" % (_buf_type, _fft_n, _complex_type))
  print("type %s [%d]%s" % (_w_type, _fft_hn, _complex_type))

  # declare function
  print("func FFT%d(in *%s, w *%s) *%s {" % (_fft_n, _buf_type, _w_type, _buf_type))
  print("var out %s" % _buf_type)
  print("var tmp %s" % _complex_type)

  # reverse the input buffer
  for i in range(_fft_n):
    print("out[%d] = in[%d]" % (reverse_bits(i, _fft_nbits), i))

  one_mask = 1
  hi_mask = -1
  lo_mask = 0
  shift = _fft_nbits - 1

  # butterfly stages
  for s in range(_fft_nbits):
    print("// stage %d" % s)
    for i in range(_fft_hn):
      j = (i&hi_mask)<<1 | (i & lo_mask)
      k = j | one_mask
      w = (i << shift) & _fft_hmask
      print("tmp = out[%d] * w[%d]" % (k, w))
      print("out[%d] = out[%d] - tmp" % (k, j))
      print("out[%d] += tmp" % j)
    shift -= 1
    one_mask <<= 1
    hi_mask <<= 1
    lo_mask = (lo_mask << 1) | 1

  print("return &out")
  print("}")

main()


"""
  // run the butterflies
  oneMask := 1
  hiMask := -1
  loMask := 0
  shift := uint(fk.stages - 1)

  for s := 0; s < fk.stages; s++ {
    for i := 0; i < fk.hn; i++ {
      j := (i&hiMask)<<1 | (i & loMask)
      k := j | oneMask
      tmp := out[k] * fk.w[(i<<shift)&fk.hmask]
      out[k] = out[j] - tmp
      out[j] += tmp
    }
    shift--
    oneMask <<= 1
    hiMask <<= 1
    loMask = (loMask << 1) | 1
  }

"""



#------------------------------------------------------------------------------
