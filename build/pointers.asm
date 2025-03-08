global _start


_start:
	mov rax, 0
	push rax

	mov rax, 5
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	lea rax, QWORD [rsp + 8]
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	mov rax, [rsp + 8]
	mov rax, [rax]
	push rax
	mov rax, 5
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	push QWORD [rsp + 0]
	pop rax
	mov QWORD [rsp + 16], rax

	mov rax, 60
	push rax
	mov rax, [rsp + 16]
	mov rax, [rax]
	push rax
	pop rdi
	pop rax
	syscall

	mov rax, 60
	mov rdi, 0
	syscall
