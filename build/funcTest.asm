global _start

set:
	mov rdi, [rsp+8] ;READ FIRST ARG
	add rax, 5	 	 ;RETURN VALUE
	ret

_start:
	mov rdi, 0

	push $5			;SET ARG (reverseOrd)
	call set		;CALL
	add rsp, 8		;RESET STACK
					;RAX contains return

	mov rbx, 50
	push rbx
	call set
	add rsp, 8

	mov rax, 60
	syscall
