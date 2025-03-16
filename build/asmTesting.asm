global _start


_start:

	mov rax, -1
	push rax

	mov rax, 5
	neg rax
	push rax

	pop rax
	pop rbx
	add rbx, rax

	mov rdi, rbx

	mov rax, 60
	syscall
