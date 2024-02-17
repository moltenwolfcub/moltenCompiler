global _start


num_0:
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


leave_2:
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


leave_3:
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	push QWORD [rbp + 16]
	push QWORD [rbp + 24]
	pop rbx
	pop rax
	add rax, rbx
	push rax
	push QWORD [rbp + 32]
	pop rbx
	pop rax
	add rax, rbx
	push rax
	mov rax, 60
	pop rdi
	syscall

	add rsp, 24
	pop rbp
	ret


test_0:
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 6
	push rax
	pop QWORD [rbp + 16]
	mov rax, 7
	push rax
	pop QWORD [rbp + 24]
	mov rax, 8
	push rax
	pop QWORD [rbp + 32]
	add rsp, 0
	pop rbp
	ret

	add rsp, 0
	pop rbp
	ret


_start:




	push 0
	push 0
	push 0
	call test_0
	add rsp, 0
	add rsp, 24

	mov rax, 4
	push rax
	push 0
	call num_0
	add rsp, 0
	call leave_2
	add rsp, 16
	add rsp, 0

	mov rax, 5
	push rax
	mov rax, 4
	push rax
	mov rax, 3
	push rax
	call leave_3
	add rsp, 24
	add rsp, 0

	mov rax, 60
	mov rdi, 0
	syscall
