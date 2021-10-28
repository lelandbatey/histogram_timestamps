#!/usr/bin/env bash

# File exists to run a command which is hard to express in a Makefile in a way
# that's easy to understand. You can run this file from anywhere and it'll
# always try to output the same thing. However, this script makes some
# assumptions which are:
# 1. This script is located inside of the `jsbuild/` directory which contains
#    all the files necessary (main.js, package.json, etc) for node to know how
#    to build 'bundle.js'.
# 2. The output directory where 'bundle.js' shall be placed is the parent
#    directory to wherever this bash script is stored on disk.

# Find the directory of this bash script, from https://stackoverflow.com/a/246128
SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
ORIGINAL_DIR="${PWD}"

cd "${SCRIPT_DIR}"

$(npm bin)/browserify main.js -o ../bundle.js

echo "bundle.js has been successfully written to the directory '${SCRIPT_DIR}/../'"
