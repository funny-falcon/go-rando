#!/bin/bash
set -x
go build head.go
go build rando.go
gcc -O3 chacha20.c -lsodium -o chacha20
gcc -O3 rando.c -lsodium -o rando_c

N=$((1024*1024*1024))
time ./head $N /dev/urandom
time ./chacha20 | ./head $N
time ./rando | ./head $N
time ./rando_c | ./head $N
