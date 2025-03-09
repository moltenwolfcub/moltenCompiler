global _start

_start:
	mov rax, 0
	push rax

	mov rax, 5
	push rax
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 0
	push rax

;==============
; important bit for making &variable work
	lea rax, QWORD [rsp + 8]
	push rax
;===========
	pop rax
	mov QWORD [rsp + 0], rax

	mov rax, 60
	push rax
	push QWORD [rsp + 8]
	pop rdi
	pop rax
	syscall

	mov rax, 60
	mov rdi, 0
	syscall

; 5
; ptr->65
; 60
; ptr->65
;
;
;

; _start:
; 	mov rax, 65
; 	push rax

; 	lea rax, [rsp + 0]
; 	push rax

; 	mov rax, 1
; 	push rax

; 	mov rax, 1
; 	push rax

; 	mov rax, QWORD [rsp + 16]
; 	push rax

; 	mov rax, 1
; 	push rax

; 	pop rdx
; 	pop rsi
; 	pop rdi
; 	pop rax

; 	syscall


; 	mov rax, 60
; 	mov rdi, 0
; 	syscall

; 65
; address of -> 65
; 1
; 1
; address of -> 65
; 1
;
;
;
