global _start


num:
	mov rax, 22
	push rax
	ret
	add rsp, 0
	ret


ex:
	mov rax, 69
	push rax
	mov rax, 60
	pop rdi
	syscall
	add rsp, 0
	ret


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

	mov rax, 15
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

label1_startWhile:
	mov rax, 1
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
	push QWORD [rsp + 0]
	mov rax, 10
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	test rax, rax
	jz label3_else
	jmp label1_startWhile
	add rsp, 0
label3_else:
	jmp label2_endWhile
	add rsp, 0
	jmp label1_startWhile
label2_endWhile:

	mov rax, 0
	push rax

	mov rax, 25
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

label4_startWhile:
	push QWORD [rsp + 0]
	pop rax
	test rax, rax
	jz label5_endWhile
	push QWORD [rsp + 0]
	mov rax, 7
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax
	mov rax, 0
	push rax
	mov rax, 2
	push rax
	pop rax
	mov QWORD [rsp + 0], rax
label6_startWhile:
	mov rax, 1
	push rax
	pop rax
	test rax, rax
	jz label7_endWhile
	push QWORD [rsp + 8]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 8], rax
	push QWORD [rsp + 0]
	mov rax, 1
	push rax
	pop rbx
	pop rax
	sub rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax
	push QWORD [rsp + 0]
	pop rax
	test rax, rax
	jz label8_else
	jmp label6_startWhile
	add rsp, 0
label8_else:
	jmp label7_endWhile
	add rsp, 0
	jmp label6_startWhile
label7_endWhile:
	add rsp, 8
	jmp label4_startWhile
label5_endWhile:

	call num
	pop rax



	mov rax, 60
	mov rdi, 0
	syscall
