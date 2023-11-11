global _start
_start:
	mov rbx, 0
	push rbx
	push QWORD [rsp + 0]
	mov rax, 60
	pop rdi
	syscall
	mov rax, 60
	mov rdi, 0
	syscall
