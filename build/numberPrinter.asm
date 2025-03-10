global _start


len_1:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 0
	push rax

	mov rax, 0
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

label1_startWhile:
	push QWORD [rbp + 16]
	pop rax
	test rax, rax
	jz label2_endWhile
	;---start_scope---
	push QWORD [rsp + 0]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	push QWORD [rbp + 16]
	mov rax, 10
	push rax
	pop rbx
	pop rax
	mov rdx, 0
	div rbx
	push rax
	pop rax
	mov QWORD [rbp + 16], rax

	add rsp, 0
	;---end_scope---
	jmp label1_startWhile
label2_endWhile:

	push QWORD [rsp + 0]
	pop QWORD [rbp + 24]
	add rsp, 8
	pop rbp
	ret

	add rsp, 8
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


exp_2:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 0
	push rax

	mov rax, 1
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

label3_startWhile:
	push QWORD [rbp + 24]
	pop rax
	test rax, rax
	jz label4_endWhile
	;---start_scope---
	push QWORD [rsp + 0]
	push QWORD [rbp + 16]
	pop rbx
	pop rax
	mul rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	push QWORD [rbp + 24]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	mov QWORD [rbp + 24], rax

	add rsp, 0
	;---end_scope---
	jmp label3_startWhile
label4_endWhile:

	push QWORD [rsp + 0]
	pop QWORD [rbp + 32]
	add rsp, 8
	pop rbp
	ret

	add rsp, 8
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


getDigit_2:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	push QWORD [rbp + 16]
	push 0
	push QWORD [rbp + 24]
	mov rax, 10
	push rax
	call exp_2
	add rsp, 16
	pop rbx
	pop rax
	mov rdx, 0
	div rbx
	push rax
	mov rax, 10
	push rax
	pop rbx
	pop rax
	mov rdx, 0
	div rbx
	push rdx
	pop QWORD [rbp + 32]
	add rsp, 0
	pop rbp
	ret

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


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


printDigit_1:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 48
	push rax
	push QWORD [rbp + 16]
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call print_1
	add rsp, 8
	add rsp, 0

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


printNumber_1:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 0
	push rax

	push 0
	push QWORD [rbp + 16]
	call len_1
	add rsp, 8
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	push QWORD [rsp + 8]
	pop rax
	mov QWORD [rsp + 0], rax

label5_startWhile:
	push QWORD [rsp + 0]
	pop rax
	test rax, rax
	jz label6_endWhile
	;---start_scope---
	push QWORD [rsp + 0]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	push 0
	push QWORD [rsp + 16]
	push QWORD [rbp + 16]
	call getDigit_2
	add rsp, 16
	pop rax
	mov QWORD [rsp + 0], rax

	push QWORD [rsp + 0]
	call printDigit_1
	add rsp, 8
	add rsp, 0

	add rsp, 8
	;---end_scope---
	jmp label5_startWhile
label6_endWhile:

	add rsp, 16
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


_start:






	mov rax, 314159
	push rax
	call printNumber_1
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
