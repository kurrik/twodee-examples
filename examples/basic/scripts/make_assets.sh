#!/usr/bin/env bash

mkdir -p tmp

aseprite \
  --batch assets/originals/numbered_squares.ase \
  --save-as tmp/numbered_squares_01.png

aseprite \
  --batch assets/originals/numbered_squares_tall.ase \
  --save-as tmp/numbered_squares_tall_01.png

aseprite \
  --batch assets/originals/numbered_squares_wide.ase \
  --save-as tmp/numbered_squares_wide_01.png

TexturePacker \
  --data assets/textures/spritesheet.json \
  --format json-array \
  --trim-sprite-names \
  --size-constraints POT \
  --disable-rotation \
  --sheet assets/textures/spritesheet.png \
  tmp

rm -rf tmp
