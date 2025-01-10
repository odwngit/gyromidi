#!/bin/bash
cd ssl
FILE=./gyromidi.crt
if [ ! -f "$FILE" ]
then
    echo Requesting permissions to make ssl_generator.sh executable...
    sudo chmod 755 ssl_generator.sh
    ./ssl_generator.sh
fi
cd ..

go run .