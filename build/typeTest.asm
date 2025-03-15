global _start


getAge_I.I.:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	push QWORD [rbp + 16]
	push QWORD [rbp + 24]
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop QWORD [rbp + 32]
	add rsp, 0
	pop rbp
	ret

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


_start:
	mov rax, 0
	push rax

	mov rax, 16
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	push 0
	mov rax, 16
	push rax
	push QWORD [rsp + 24]
	call getAge_I.I.
	add rsp, 16
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	mov rax, 0
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	mov rax, 65
	push rax
	pop rax
	mov QWORD [rsp + 0], rax


	mov rax, 60
	mov rdi, 0
	syscall
