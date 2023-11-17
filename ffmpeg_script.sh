#!/bin/bash

ffmpeg -listen 1 -i rtmp://localhost:1935/test \
  -f flv \
  -c copy \
  -map 0 \
  -f segment -segment_list playlist.m3u8 \
  -segment_list_flags +live -segment_time 2 \
  out%03d.ts
