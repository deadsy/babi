#!/bin/bash

DEV="hw:CARD=sndrpihifiberry,DEV=0"
#DEV="plughw:CARD=sndrpihifiberry,DEV=0"

alsabat -P$DEV -f S16_LE -r 48000 -c 2
