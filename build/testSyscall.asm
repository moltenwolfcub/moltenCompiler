global _start


_start:
	mov rax, 60
	push rax
	mov rax, 3
	push rax
	mov rax, 2
	push rax
	pop rbx
	pop rax
	add rax, rbx
	push rax
	pop rdi
	pop rax
	syscall

	mov rax, 60
	mov rdi, 0
	syscall
