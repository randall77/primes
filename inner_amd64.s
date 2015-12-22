// takes sieve slice, 8 indexes, 8 masks, p.

TEXT ·inner(SB), 4, $0-48
	MOVQ	a+0(FP), BP	// sieve pointer
	MOVQ	a+8(FP), R14	// sieve length

	MOVQ	a+24(FP), CX	// ptr to indexes
	MOVQ	(CX), R8	// 8 indexes
	MOVQ	8(CX), R9
	MOVQ	16(CX), R10
	MOVQ	24(CX), AX
	MOVQ	32(CX), R12
	MOVQ	40(CX), SI
	MOVQ	48(CX), R13
	MOVQ	56(CX), DI

	MOVW	a+32(FP), R11	// 8 byte masks (loaded 2 at a time)
	MOVW	a+34(FP), BX
	MOVW	a+36(FP), CX
	MOVW	a+38(FP), DX
	
	MOVQ	a+40(FP), R15	// p
	
	CMPQ	DI, R14
	JAE	tail

loop:
	ANDB	R11, (BP)(R8*1)
	RORQ	$8, R11
	ANDB	R11, (BP)(R9*1)
	ROLQ	$8, R11
	ANDB	BL, (BP)(R10*1)
	ANDB	BH, (BP)(AX*1)
	ANDB	CL, (BP)(R12*1)
	ANDB	CH, (BP)(SI*1)
	ANDB	DL, (BP)(R13*1)
	ANDB	DH, (BP)(DI*1)
	ADDQ	R15, BP
	SUBQ	R15, R14
	CMPQ	DI, R14
	JB	loop

tail:
	CMPQ	R8, R14
	JAE	exit
	ANDB	R11, (BP)(R8*1)

	CMPQ	R9, R14
	JAE	exit
	SHRQ	$8, R11
	ANDB	R11, (BP)(R9*1)

	CMPQ	R10, R14
	JAE	exit
	ANDB	BL, (BP)(R10*1)

	CMPQ	AX, R14
	JAE	exit
	SHRQ	$8, BX
	ANDB	BL, (BP)(AX*1)

	CMPQ	R12, R14
	JAE	exit
	ANDB	CL, (BP)(R12*1)

	CMPQ	SI, R14
	JAE	exit
	ANDB	CH, (BP)(SI*1)

	CMPQ	R13, R14
	JAE	exit
	ANDB	DL, (BP)(R13*1)
exit:
	RET
