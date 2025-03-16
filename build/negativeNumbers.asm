global _start


_start:
	mov rax, 0
	push rax

	mov rax, 5
	push rax
	pop rax
	neg rax
	push rax
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
