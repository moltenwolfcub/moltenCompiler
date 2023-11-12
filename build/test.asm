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
	mov rax, 9
	push rax
	pop rbx
	pop rax
	mul rbx
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

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

	push QWORD [rsp + 0]
	mov rax, 4
	push rax
	pop rbx
	pop rax
	div rbx
	push rax
	mov rax, 2
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
