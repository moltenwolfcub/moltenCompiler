global _start


num:
	;=====FUNCTION BODY=====
	mov rax, 22
	push rax
	ret

	add rsp, 0
	ret


_start:

	mov rax, 0
	push rax

	call num
	add rsp, 0
	pop rax
	mov QWORD [rsp + 0], rax

	push QWORD [rsp + 0]
	mov rax, 60
	pop rdi
	syscall

	mov rax, 60
	mov rdi, 0
	syscall
