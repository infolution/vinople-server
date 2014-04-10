#!/bin/sh

case $1 in
	dev)
		GO=/usr/local/go/bin/go
		;;
	prod)
		GO=/home/admin/build/golang/go/bin/go
		;;
	test)
		GO=/home/admin/build/golang/go/bin/go
		;;
	*)
		echo "Usage: $0 {dev|prod|test}"
		exit 1
esac

$GO build -o bin/vinople-server src/*
