#!/bin/bash

CONNECT=/usr/local/bin/jack_connect

$CONNECT system:midi_capture_1 simple:midi_in_0
$CONNECT system:playback_3 simple:audio_out_0
$CONNECT system:playback_4 simple:audio_out_1

#$CONNECT system:playback_3 jack_simple_client:output1
#$CONNECT system:playback_4 jack_simple_client:output2
