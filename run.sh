#! /bin/bash

nasm -felf64 build/$1 -o build/out.o && ld build/out.o -o build/out && ./build/out; echo $?
