#!/bin/zsh

export SENZING_PATH=${HOME}/senzing
export DYLD_LIBRARY_PATH=${SENZING_PATH}/er/lib:/opt/homebrew/lib
export LD_LIBRARY_PATH=${DYLD_LIBRARY_PATH}

"$@"
