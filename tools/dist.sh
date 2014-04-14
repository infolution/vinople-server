#!/bin/sh

if [ -z "$1" ]
 then
    echo "Usage: $0 path/to/install/dir"
    exit 1
fi

cp -R bin $1
