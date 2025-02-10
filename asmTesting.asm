global _start

; rax | rdi, rsi, rdx, r10, r8, r9

_start:
	; ;open file
	; mov rax, 2
	; mov rdi, pathname
	; mov rsi, 00000000
	; syscall
	; push rax

	; ;read file to buffer
	; mov rax, 0
	; pop rdi
	; mov rsi, buffer
	; mov rdx, 1024
	; syscall


	; mov rax, 74
	; mov [buffer], rax

	; push character to stack
	push 72
	push 69
	push 76
	push 76
	push 79
	push 10

	; setup syscall
	mov rax, 1
	mov rdi, 1
	mov rdx, 1

	; load first value
	lea rsi, [rsp+40]
	syscall
	
	lea rsi, [rsp+32]
	syscall

	lea rsi, [rsp+24]
	syscall

	lea rsi, [rsp+16]
	syscall

	lea rsi, [rsp+8]
	syscall

	lea rsi, [rsp]
	syscall

	; exit
	mov rax, 60
	mov rdi, 0
	syscall

section .data
	pathname DD "/home/oliver/Desktop/dev/go/compilers/molten/build/read.txt"
; section .bss
; 	buffer: resb 1024
