global _start
_start:
	mov rax, 0
	push rax

	mov rax, 15
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 10
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

	mov rax, 3
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	push QWORD [rsp + 0]
	pop rax
	mov QWORD [rsp + 8], rax

	push QWORD [rsp + 8]
	mov rax, 60
	pop rdi
	syscall

	mov rax, 60
	mov rdi, 0
	syscall
