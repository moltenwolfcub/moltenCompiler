#! /bin/bash

nasm -felf64 build/test.asm && ld build/test.o -o build/test && ./build/test; echo $?
