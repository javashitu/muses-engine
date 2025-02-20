#!/bin/bash

function transcode360p(){
    ffmpeg -i $1 -c:v libx264 -b:v 256k -s 480x360 -r 25 -c:a copy $2
}

transcode360p $1 $2