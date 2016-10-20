#!/bin/sh

# This script is meant to be used to launch the application inside of its Docker container.

# Dump udocs env to stdout for debugging
/usr/local/bin/udocs env

# Run the app in headless mode, seeding any Git repos that were specified in the env vars
/usr/local/bin/udocs serve --headless
