#!/bin/bash

# ffmpeg -loglevel verbose -listen 1 -i rtmp://localhost:1935/test -f flv -c copy \
#   -c:v libx264 -x264opts keyint=120:no-scenecut -s 1920x1080 -r 60 -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_list ./1080_60/playlist_1080_60.m3u8  -segment_list_flags +live -segment_time 2 ./1080_60/out_1080_60_%03d.ts \
#   -c:v libx264 -x264opts keyint=120:no-scenecut -s 1280x720 -r 60 -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_list ./720_60/playlist_720_60.m3u8  -segment_list_flags +live -segment_time 2 ./720_60/out_720_60_%03d.ts \
#   -c:v libx264 -x264opts keyint=60:no-scenecut -s 1280x720 -r 30  -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_list ./720_30/playlist_720_30.m3u8  -segment_list_flags +live -segment_time 2 ./720_30/out_720_30_%03d.ts \
#   -c:v libx264 -x264opts keyint=60:no-scenecut -s 852x480 -r 30  -profile:v main -preset veryfast -c:a aac -sws_flags bilinear  -f segment -segment_list ./480_30/playlist_480_30.m3u8  -segment_list_flags +live -segment_time 2 ./480_30/out_480_30_%03d.ts 

# ffmpeg -loglevel verbose -listen 1 -i rtmp://localhost:1935/test -f flv -c copy \
#   -c:v libx264 -x264opts keyint=120:no-scenecut -s 1920x1080 -r 60 -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_time 2 ./1080_60/out_1080_60_%03d.ts \
#   -c:v libx264 -x264opts keyint=120:no-scenecut -s 1280x720 -r 60 -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_time 2 ./720_60/out_720_60_%03d.ts \
#   -c:v libx264 -x264opts keyint=60:no-scenecut -s 1280x720 -r 30  -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_time 2 ./720_30/out_720_30_%03d.ts \
#   -c:v libx264 -x264opts keyint=60:no-scenecut -s 852x480 -r 30  -profile:v main -preset veryfast -c:a aac -sws_flags bilinear  -f segment -segment_time 2 ./1/202312021430/852x480_30/out_480_30_%03d.ts 

# ffmpeg -loglevel verbose -listen 1 -i rtmp://localhost:1935/test \ 
#   -c:v libx264 -x264opts keyint=30:no-scenecut -s 1920x1080 -r 30 -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_time 2 ./temp/1/202312021430/1920x1080_30/out_%03d.ts \
#   -c:v libx264 -x264opts keyint=60:no-scenecut -s 1280x720 -r 30 -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_time 2 ./temp/1/202312021430/1280x720_30/out_%03d.ts \
#   -c:v libx264 -x264opts keyint=60:no-scenecut -s 852x480 -r 30 -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_time 2 ./temp/1/202312021430/852x480_30/out_%03d.ts 

ffmpeg -loglevel verbose -listen 1 -i rtmp://localhost:1935/test -f flv -c copy \
  -c:v libx264 -x264opts keyint=60:no-scenecut -s 852x480 -r 30  -profile:v main -preset veryfast -c:a aac -sws_flags bilinear  -f segment -segment_time 2 "./temp/1/202312021430/852x480_30/%06d.ts"


ffmpeg -loglevel verbose -listen 1 -i rtmp://localhost:1935/test -f flv -c copy \
  -c:v libx264 -x264opts keyint=60:no-scenecut -s 852x480 -r 30  -profile:v main -preset veryfast -c:a aac -sws_flags bilinear  -f segment -segment_time 2 "./out_%06d.ts"