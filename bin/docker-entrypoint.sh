#!/bin/sh

# This script is meant to be used to launch the application inside
# its docker container.

# source the environment variables
. /usr/local/udocs.env

# Run the app
/usr/local/bin/udocs serve&
while true; do
  udocs publish
  sleep 10
done
