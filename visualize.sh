#!/bin/bash

# Usage: ./visualize.sh [package] [output]

godepgraph -nostdlib -novendor $1 | dot -Tpng -o $2
