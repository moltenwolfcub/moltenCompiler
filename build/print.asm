global _start


print_1:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 1
	push rax
	mov rax, 1
	push rax
	lea rax, QWORD [rbp + 16]
	push rax
	mov rax, 1
	push rax
	pop rdx
	pop rsi
	pop rdi
	pop rax
	syscall

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


_start:

	mov rax, 64
	push rax
	mov rax, 8
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 96
	push rax
	mov rax, 5
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 96
	push rax
	mov rax, 12
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 96
	push rax
	mov rax, 12
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 96
	push rax
	mov rax, 15
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 32
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 64
	push rax
	mov rax, 23
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 96
	push rax
	mov rax, 15
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 96
	push rax
	mov rax, 18
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 96
	push rax
	mov rax, 12
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 96
	push rax
	mov rax, 4
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 33
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 10
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	mov rax, 60
	mov rdi, 0
	syscall
