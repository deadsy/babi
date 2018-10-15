#!/usr/bin/env python3
#------------------------------------------------------------------------------
"""

Generate a loopless FFT function (in golang) for a given input length.

"""
#------------------------------------------------------------------------------

_package = "core"
_fp_length = 32
_fft_length = 256
_complex_type = ("complex128", "complex64")[_fp_length == 32]

#------------------------------------------------------------------------------

def main():
  print("package %s" % _package)
  print("func FFT%d(in []%s) []%s {" % (_fft_length, _complex_type, _complex_type))
  print("out := make([]%s, %d)" % (_complex_type, _fft_length))
  print("return out")
  print("}")

main()

#------------------------------------------------------------------------------
