#!/bin/bash

JACK=/usr/local/bin/jackd

PLATFORM=$(uname -m)

case $PLATFORM in
	x86_64)
		DEV="hw:CARD=PCH,DEV=0"
	;;
	armv7l)
		DEV="hw:CARD=sndrpihifiberry,DEV=0"
	;;
	*)
		echo "unknown platform"
		exit 1
	;;
esac

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
