global _start

gen_code:

	mov rax, 5
	mov [rsp + 16], rax
	mov rax, 10
	mov [rsp + 8], rax
	ret

_start:

	push 0
	push 0
	call gen_code
	pop rdi
	pop rbx

	mov rax, 60
	syscall
