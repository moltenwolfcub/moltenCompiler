global _start


leave:
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	push QWORD [rbp + 16]
	push QWORD [rbp + 24]
	pop rbx
	pop rax
	add rax, rbx
	push rax
	mov rax, 60
	pop rdi
	syscall

	add rsp, 16
	pop rbp
	ret


_start:

	mov rax, 4
	push rax
	mov rax, 32
	push rax
	call leave
	add rsp, 16
	add rsp, 0

	mov rax, 60
	mov rdi, 0
	syscall
