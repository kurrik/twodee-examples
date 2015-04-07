#!/usr/bin/env bash

mkdir -p tmp

aseprite \
  --batch assets/originals/numbered_squares.ase \
  --save-as tmp/numbered_squares_01.png

TexturePacker \
  --data assets/textures/spritesheet.json \
  --format json-array \
  --trim-sprite-names \
  --size-constraints POT \
  --sheet assets/textures/spritesheet.png \
  tmp

rm -rf tmp
