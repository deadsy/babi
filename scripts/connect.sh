#!/bin/bash

CONNECT=/usr/local/bin/jack_connect

$CONNECT system:midi_capture_1 babi:midi_in_0
$CONNECT system:playback_3 babi:audio_out_0
$CONNECT system:playback_4 babi:audio_out_1
