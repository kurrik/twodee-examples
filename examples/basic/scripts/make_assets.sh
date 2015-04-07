#!/usr/bin/env bash

mkdir -p tmp

aseprite \
  --batch basic/assets/originals/numbered_squares.ase \
  --save-as tmp/numbered_squares_01.png

TexturePacker \
  --data basic/assets/textures/spritesheet.json \
  --format json-array \
  --sheet basic/assets/textures/spritesheet.png \
  tmp

rm -rf tmp
