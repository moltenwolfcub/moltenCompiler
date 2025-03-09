global _start


test_1:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 0
	push rax

	push QWORD [rbp + 16]
	pop rax
	mov QWORD [rsp + 0], rax

	push QWORD [rsp + 0]
	pop QWORD [rbp + 24]
	add rsp, 8
	pop rbp
	ret

	add rsp, 8
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


_start:

	mov rax, 0
	push rax

	mov rax, 15
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	push 0
	mov rax, 5
	push rax
	call test_1
	add rsp, 8
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 60
	push rax
	push QWORD [rsp + 8]
	pop rdi
	pop rax
	syscall

	mov rax, 60
	mov rdi, 0
	syscall
