#!/bin/sh

if ! command -v glide >/dev/null 2>&1; then
  echo 'no exists glide'
  exit 1
fi

glide update -v

