global _start

set:
	mov rdi, [rsp+8] ;READ FIRST ARG
	add rax, rdi	 	 ;RETURN VALUE
	ret

_start:
	mov rdi, 0

	mov rbx, 50
	push rbx		;SET ARG (reverseOrd)
	call set		;CALL
	add rsp, 8		;RESET STACK
					;RAX contains return

	mov rdi, rax

	mov rax, 60
	syscall
