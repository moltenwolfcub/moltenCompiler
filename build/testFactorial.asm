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


_start:

	mov rax, 0
	push rax

	mov rax, 7
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	mov rax, 1
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

label1_startWhile:
	mov rax, 1
	push rax
	pop rax
	test rax, rax
	jz label2_endWhile
	;---start_scope---
	push QWORD [rsp + 0]
	push QWORD [rsp + 16]
	pop rbx
	pop rax
	mul rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	push QWORD [rsp + 8]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 8], rax

	push QWORD [rsp + 8]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	test rax, rax
	jz label3_else
	;---start_scope---
	jmp label1_startWhile

	add rsp, 0
	;---end_scope---
label3_else:

	jmp label2_endWhile

	add rsp, 0
	;---end_scope---
	jmp label1_startWhile
label2_endWhile:

	push QWORD [rsp + 0]
	call exit_1
	add rsp, 8
	add rsp, 0

	mov rax, 60
	mov rdi, 0
	syscall
