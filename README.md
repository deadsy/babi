# babi
*It's Babi Guling time Dave.*

[![Go Report Card](https://goreportcard.com/badge/github.com/deadsy/babi)](https://goreportcard.com/report/github.com/deadsy/babi)
[![GoDoc](https://godoc.org/github.com/deadsy/babi/core?status.svg)](https://godoc.org/github.com/deadsy/babi/core)

## What is It?
Babi is a yet another modular software synthesizer. This one is written in Golang.

## What does it run on?

Anything that can run a golang program and supports a JACK server.
I normally run it on Linux PC/RaspberryPi.

## How do I run it?
* Start JACK. E.g. "./scripts/jack.sh start"
* Start an example program. E.g. "./cmd/poly/poly"
* Connect the client to JACK. E.g. "./scripts/connect.sh"
* Start jammin' on your MIDI input device.

## Specifications
* 32-bit floats for DSP operations
* 48000 samples/sec (compile time constant)
* 128 samples/buffer (compile time constant)
* Connects to the world as a JACK client.

## Resources
* https://golang.org/
* http://jackaudio.org/

## What's a Module?
A module is a software analog to the modules you might find in a hardware modular synth (only cheaper and less tactile).

* Modules are processing elements with input and output ports.
* Modules can have sub-modules.
* Modules can dynamically create and destroy modules while they are running.

## Patch Module
A patch is a module suitable for use as the top-level module of the synthesizer.

It normally has:

* Audio inputs (audio streams coming from outside sources) 
* Audio outputs (audio streams going to outside destinations)
* MIDI inputs (MIDI events coming in to control the synthesizer)
* MIDI outputs (generated MIDI events sent out)

These module ports will be mapped to JACK ports which are then connected to a JACK server.

## Voice Module
A voice is a module which outputs audio for a single note.

It has standard port interfaces:
* Gate
* Frequency
* Audio Output

A voice will typically be a submodule to a polyphonic module which will then allocate and run multiple voices concurrently.
