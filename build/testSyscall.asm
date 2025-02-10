global _start


_exit_1:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 60
	push rax
	push QWORD [rbp + 16]
	pop rdi
	pop rax
	syscall

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


_start:

	mov rax, 3
	push rax
	mov rax, 8
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call _exit_1
	add rsp, 8
	add rsp, 0

	mov rax, 60
	mov rdi, 0
	syscall
