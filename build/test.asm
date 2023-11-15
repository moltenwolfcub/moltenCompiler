global _start
_start:
	mov rax, 0
	push rax

	mov rax, 14
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	push QWORD [rsp + 8]
	mov rax, 4
	push rax
	mov rax, 6
	push rax
	pop rbx
	pop rax
	mul rbx
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	mov rax, 3
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax
	mov rax, 6
	push rax
	pop rax
	mov QWORD [rsp + 0], rax
	mov rax, 0
	push rax
	mov rax, 7
	push rax
	pop rax
	mov QWORD [rsp + 0], rax
	mov rax, 2
	push rax
	push QWORD [rsp + 8]
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 8], rax
	add rsp, 8
	push QWORD [rsp + 8]
	push QWORD [rsp + 8]
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 8], rax
	add rsp, 8

	mov rax, 0
	push rax

	push QWORD [rsp + 8]
	mov rax, 10
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	mov rax, 5
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

label1_startWhile:
	push QWORD [rsp + 0]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	test rax, rax
	jz label2_endWhile
	push QWORD [rsp + 0]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax
	add rsp, 0
	jmp label1_startWhile
label2_endWhile:

	push QWORD [rsp + 0]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	test rax, rax
	jz label3_if
	mov rax, 3
	push rax
	mov rax, 60
	pop rdi
	syscall
	add rsp, 0
label3_if:

	push QWORD [rsp + 8]
	mov rax, 4
	push rax
	pop rbx
	pop rax
	div rbx
	push rax
	mov rax, 2
	push rax
	mov rax, 4
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop rbx
	pop rax
	mul rbx
	push rax
	mov rax, 60
	pop rdi
	syscall

	mov rax, 60
	mov rdi, 0
	syscall
