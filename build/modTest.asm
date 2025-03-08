global _start


_start:
	mov rax, 0
	push rax

	mov rax, 11
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	mov rax, 3
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	push QWORD [rsp + 16]
	push QWORD [rsp + 16]
	pop rbx
	pop rax
	mov rdx, 0
	div rbx
	push rdx
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	push QWORD [rsp + 24]
	push QWORD [rsp + 16]
	pop rbx
	pop rax
	mov rdx, 0
	div rbx
	push rdx
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
