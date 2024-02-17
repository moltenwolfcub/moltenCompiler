global _start


num:
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 22
	push rax
	pop QWORD [rbp + 16]
	add rsp, 0
	pop rbp
	ret

	add rsp, 0
	pop rbp
	ret


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
	push 0
	call num
	add rsp, 0
	call leave
	add rsp, 16
	add rsp, 0

	mov rax, 60
	mov rdi, 0
	syscall
