global _start


exit_1:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 60
	push rax
	push QWORD [rbp + 16]
	pop rdi
	pop rax
	syscall

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


num_0:
	;=====FUNCTION SETUP=====
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
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


leave_2:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	push QWORD [rbp + 16]
	push QWORD [rbp + 24]
	pop rbx
	pop rax
	add rax, rbx
	push rax
	call exit_1
	add rsp, 8
	add rsp, 0

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


leave_3:
	;=====FUNCTION SETUP=====
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
	call exit_1
	add rsp, 8
	add rsp, 0

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


empty_0:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


ret_0:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	add rsp, 0
	pop rbp
	ret

	add rsp, 0
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


pmRet_1:
	;=====FUNCTION SETUP=====
	push rbp
	mov rbp, rsp

	;=====FUNCTION BODY=====
	mov rax, 0
	push rax

	;---start_scope---
	mov rax, 0
	push rax

	;---start_scope---
	mov rax, 0
	push rax

	add rsp, 8
	;---end_scope---

	add rsp, 8
	;---end_scope---

	add rsp, 8
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


test_0:
	;=====FUNCTION SETUP=====
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
	;=====FUNCTION CLEANUP=====
	pop rbp
	ret


_start:








	push 0
	push 0
	push 0
	call test_0
	add rsp, 0
	add rsp, 24

	call empty_0
	add rsp, 0
	add rsp, 0

	call ret_0
	add rsp, 0
	add rsp, 0

	mov rax, 5
	push rax
	call pmRet_1
	add rsp, 8
	add rsp, 0

	mov rax, 60
	mov rdi, 0
	syscall
