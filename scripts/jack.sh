#!/bin/bash

JACK=/usr/local/bin/jackd

DEV="hw:CARD=sndrpihifiberry,DEV=0"

case "$1" in
  start)
    $JACK -r -d alsa -d $DEV -S -r 48000 -P &
  ;;
  stop)
    killall -s SIGHUP jackd
  ;;
  restart)
    $0 stop
    $0 start
  ;;
  *)
    echo "Usage: $0 (start|stop|restart)"
    exit 1
esac

exit $?
