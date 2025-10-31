#!/bin/bash
#
# This file is part of the PyRDP project.
# Copyright (C) 2022-2023 GoSecure Inc.
# Licensed under the GPLv3 or later.
#
# We extracted a set of important tests that were run as part of a GitHub
# workflow before. Having them all here makes them easy to run from a
# development environment. The GitHub worfklows can still run them.
#
# NOTE: Running these locally requires the test/files/test_files.zip to be
#       extracted in test/files/.

# Any non-zero exit code becomes an error now
set -e

# Sets how to launch commands. GitHub workflows sets the CI environment variable
if [[ -z "${CI}" ]]; then
	SOURCE="local"
else
	SOURCE="ci"
fi

declare -A ENTRYPOINT

ENTRYPOINT["player","local"]="pyrdp-player"
ENTRYPOINT["player","ci"]="coverage run --append -m pyrdp.bin.player"
ENTRYPOINT["convert","local"]="pyrdp-convert"
ENTRYPOINT["convert","ci"]="coverage run --append -m pyrdp.bin.convert"

export QT_QPA_PLATFORM=offscreen

echo ===================================================
echo pyrdp-player read a replay in headless mode test
${ENTRYPOINT["player",${SOURCE}]} --headless test/files/test_session.replay
echo

echo ===================================================
echo pyrdp-convert to MP4
${ENTRYPOINT["convert",${SOURCE}]} test/files/test_convert.pyrdp -f mp4
echo

echo ===================================================
echo Verify the MP4 file
file test_convert.mp4 | grep "MP4 Base Media"
rm test_convert.mp4
echo

echo ===================================================
echo pyrdp-convert replay to JSON
${ENTRYPOINT["convert",${SOURCE}]} test/files/test_convert.pyrdp -f json
echo

echo ===================================================
echo Verify the replay to JSON file
./test/validate_json.sh test_convert.json
rm test_convert.json
echo

echo ===================================================
echo pyrdp-convert PCAP to JSON
${ENTRYPOINT["convert",${SOURCE}]} test/files/test_session.pcap -f json
echo

echo ===================================================
echo Verify the PCAP to JSON file
./test/validate_json.sh "20200319000716_192.168.38.1_20989-192.168.38.1_3389.json"
rm "20200319000716_192.168.38.1_20989-192.168.38.1_3389.json"
echo

echo ===================================================
echo pyrdp-convert PCAP to replay
${ENTRYPOINT["convert",${SOURCE}]} test/files/test_session.pcap -f replay
echo

echo ===================================================
echo Verify that the replay file exists
file -E "20200319000716_192.168.38.1_20989-192.168.38.1_3389.pyrdp"
rm "20200319000716_192.168.38.1_20989-192.168.38.1_3389.pyrdp"

echo ===================================================
echo pyrdp-convert.py regression issue 428
${ENTRYPOINT["convert",${SOURCE}]} test/files/test_convert_428.pyrdp -f mp4
echo

echo ===================================================
echo Verify the MP4 file
file test_convert_428.mp4 | grep "MP4 Base Media"
rm test_convert_428.mp4
echo
