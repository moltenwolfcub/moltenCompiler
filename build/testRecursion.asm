global _start


exit_1:
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


factorial_1:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	push QWORD [rbp + 16]
	pop rax
	test rax, rax
	jz label1_else
	;---start_scope---
	push QWORD [rbp + 16]
	push 0
	push QWORD [rbp + 16]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	call factorial_1
	add rsp, 8
	pop rbx
	pop rax
	mul rbx
	push rax
	pop QWORD [rbp + 24]
	add rsp, 0
	pop rbp
	ret

	add rsp, 0
	;---end_scope---
label1_else:

	mov rax, 1
	push rax
	pop QWORD [rbp + 24]
	add rsp, 0
	pop rbp
	ret

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


_start:


	push 0
	mov rax, 5
	push rax
	call factorial_1
	add rsp, 8
	call exit_1
	add rsp, 8
	add rsp, 0

	mov rax, 60
	mov rdi, 0
	syscall
