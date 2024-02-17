global _start

gen_code:
	push rbp
	mov rbp, rsp

	mov rax, 5
	mov [rbp + 16], rax
	mov rax, 10
	mov [rsp + 24], rax

	pop rbp
	ret

_start:

	push 0
	push 0
	;push parameters
	call gen_code
	;clean parameters
	;clean or access returns
	pop rdi
	pop rbx

	mov rax, 60
	syscall
