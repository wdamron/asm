package i64

// optab is a table of Op+Reg+Reg combinations matched against reprenstation data
// (op code, prefix requirements, etc). Some of the table is listed directly, and
// some of the more repetitive sections of the table are generated by init.
var optab = map[opKey]opVal{
	opKey{CMPQ, Reg, Reg}:     opVal{c1: 0x3b, rex: true},
	opKey{CMPQ, Reg, Ind}:     opVal{c1: 0x3b, rex: true},
	opKey{PUSHQ, Imm32, None}: opVal{c1: 0x68, mod: modNone},
	opKey{PUSHQ, Imm8, None}:  opVal{c1: 0x6a, mod: modNone},
	opKey{PUSHQ, Reg, None}:   opVal{c1: 0x50, addReg: true, mod: modNone},
	opKey{POPQ, None, Reg}:    opVal{c1: 0x58, addReg: true, mod: modNone},

	opKey{JE, None, Rel8}:   opVal{c1: 0x74, mod: modNone},
	opKey{JE, None, Rel16}:  opVal{c1: 0x0f, c2: 0x84, mod: modNone},
	opKey{JE, None, Rel32}:  opVal{c1: 0x0f, c2: 0x84, mod: modNone},
	opKey{JNE, None, Rel8}:  opVal{c1: 0x75, mod: modNone},
	opKey{JNE, None, Rel16}: opVal{c1: 0x0f, c2: 0x85, mod: modNone},
	opKey{JNE, None, Rel32}: opVal{c1: 0x0f, c2: 0x85, mod: modNone},
	opKey{JHI, None, Rel8}:  opVal{c1: 0x77, mod: modNone},
	opKey{JHI, None, Rel16}: opVal{c1: 0x0f, c2: 0x87, mod: modNone},
	opKey{JHI, None, Rel32}: opVal{c1: 0x0f, c2: 0x87, mod: modNone},

	opKey{MOVB, Reg, Reg}: opVal{c1: 0x8a},
	opKey{MOVB, Ind, Reg}: opVal{c1: 0x8a},
	opKey{MOVB, Reg, Ind}: opVal{c1: 0x88},
	opKey{MOVL, Reg, Reg}: opVal{c1: 0x8b},
	opKey{MOVL, Ind, Reg}: opVal{c1: 0x8b},
	opKey{MOVL, Reg, Ind}: opVal{c1: 0x89},
	opKey{MOVQ, Reg, Reg}: opVal{c1: 0x8b, rex: true},
	opKey{MOVQ, Ind, Reg}: opVal{c1: 0x8b, rex: true},
	opKey{MOVQ, Reg, Ind}: opVal{c1: 0x89, rex: true},

	opKey{LEAL, Reg, Ind}: opVal{c1: 0x8d},
	opKey{LEAQ, Reg, Ind}: opVal{c1: 0x8d, rex: true},

	opKey{RET, None, None}: opVal{c1: 0xc3, mod: modNone},
	opKey{MOVQ, Imm8, Reg}: opVal{c1: 0xc6, rex: true, mod: mod0},
	opKey{MOVQ, Imm8, Ind}: opVal{c1: 0xc6, rex: true, mod: mod0},

	opKey{MOVSS, Ind, Xmm}: opVal{c0: 0xf3, c1: 0x0f, c2: 0x10},
	opKey{MOVSS, Xmm, Ind}: opVal{c0: 0xf3, c1: 0x0f, c2: 0x11},
	opKey{ADDSS, Xmm, Xmm}: opVal{c0: 0xf3, c1: 0x0f, c2: 0x58},
	opKey{MULSS, Xmm, Xmm}: opVal{c0: 0xf3, c1: 0x0f, c2: 0x59},
	opKey{SUBSS, Xmm, Xmm}: opVal{c0: 0xf3, c1: 0x0f, c2: 0x5c},
	opKey{MINSS, Xmm, Xmm}: opVal{c0: 0xf3, c1: 0x0f, c2: 0x5d},
	opKey{DIVSS, Xmm, Xmm}: opVal{c0: 0xf3, c1: 0x0f, c2: 0x5e},
	opKey{MAXSS, Xmm, Xmm}: opVal{c0: 0xf3, c1: 0x0f, c2: 0x5f},

	opKey{MOVSD, Ind, Xmm}: opVal{c0: 0xf2, c1: 0x0f, c2: 0x10},
	opKey{MOVSD, Xmm, Ind}: opVal{c0: 0xf2, c1: 0x0f, c2: 0x11},
	opKey{ADDSD, Xmm, Xmm}: opVal{c0: 0xf2, c1: 0x0f, c2: 0x58},
	opKey{MULSD, Xmm, Xmm}: opVal{c0: 0xf2, c1: 0x0f, c2: 0x59},
	opKey{SUBSD, Xmm, Xmm}: opVal{c0: 0xf2, c1: 0x0f, c2: 0x5c},
	opKey{MINSD, Xmm, Xmm}: opVal{c0: 0xf2, c1: 0x0f, c2: 0x5d},
	opKey{DIVSD, Xmm, Xmm}: opVal{c0: 0xf2, c1: 0x0f, c2: 0x5e},
	opKey{MAXSD, Xmm, Xmm}: opVal{c0: 0xf2, c1: 0x0f, c2: 0x5f},
}

func init() {
	regs := []AddrType{Reg, Ind, Xmm, Imm8, Imm16, Imm32, Imm64, Rel8, Rel16, Rel32}
	expand := func(r AddrType) (addrs []AddrType) {
		if r == None {
			return []AddrType{None}
		}
		for _, reg := range regs {
			if r&reg != 0 {
				addrs = append(addrs, reg)
			}
		}
		return addrs
	}
	add := func(op Op, r1, r2 AddrType, c opVal) {
		regs2 := expand(r2)
		for _, reg1 := range expand(r1) {
			for _, reg2 := range regs2 {
				optab[opKey{op, reg1, reg2}] = c
			}
		}
	}
	for i := ADD; i <= CMP; i++ {
		m := mod0 + modBits(i-ADD)
		add(i, Imm8, Reg|Ind, opVal{c1: 0x80, mod: m})
	}
	for i := ADDL; i <= CMPL; i++ {
		m := mod0 + modBits(i-ADDL)
		add(i, Imm16|Imm32, Reg|Ind, opVal{c1: 0x81, mod: m})
		add(i, Imm8, Reg|Ind, opVal{c1: 0x83, mod: m})
		opOff := uint8(i-ADDQ) * 8
		add(i, Reg|Ind, Reg, opVal{c1: opOff + 0x01})
		add(i, Reg, Ind, opVal{c1: opOff + 0x03})
	}
	for i := ADDQ; i <= CMPQ; i++ {
		m := mod0 + modBits(i-ADDQ)
		add(i, Imm16|Imm32, Reg|Ind, opVal{c1: 0x81, rex: true, mod: m})
		add(i, Imm8, Reg|Ind, opVal{c1: 0x83, rex: true, mod: m})
		opOff := uint8(i-ADDQ) * 8
		add(i, Reg|Ind, Reg, opVal{c1: opOff + 0x01, rex: true})
		add(i, Reg, Ind, opVal{c1: opOff + 0x03, rex: true})
	}
	add(MOVL, Imm16|Imm32|Imm64, Reg, opVal{c1: 0xb8, addReg: true, mod: modNone})
	add(MOVQ, Imm16|Imm32, Reg|Ind, opVal{c1: 0xc7, rex: true, mod: mod0})
	add(MOVQ, Imm64, Reg, opVal{c1: 0xb8, addReg: true, rex: true, mod: modNone})
	add(IMULL, Reg, Reg|Ind, opVal{c1: 0x0f, c2: 0xaf})
	add(IMULQ, Reg, Reg|Ind, opVal{c1: 0x0f, c2: 0xaf, rex: true})
	add(IDIVL, None, Reg|Ind, opVal{c1: 0xf7, mod: mod7})
	add(IDIVQ, None, Reg|Ind, opVal{c1: 0xf7, rex: true, mod: mod7})
	add(CALL, None, Rel16|Rel32, opVal{c1: 0xe8, mod: modNone})
	add(JMP, None, Rel16|Rel32, opVal{c1: 0xe9, mod: modNone})
}

type opKey struct {
	Op   Op
	From AddrType
	To   AddrType
}

type opVal struct {
	c0     uint8 // 1-byte prefix, either 66, f3, or f2.
	c1     uint8 // 1-byte op code
	c2     uint8 // 2nd byte of op code. Only used if code==0x0f.
	rex    bool  // REX prefix is present, W is set.
	addReg bool  // add the register number to the op code.
	mod    modBits
}

// modBits describes what the ModRM.mod bits are used for.
type modBits int

const (
	// modDefault means ModRM.mod designates a register.
	// (If this instruction has a ModRM byte.)
	modDefault modBits = iota

	// modNone means no ModRM byte for this instruction.
	modNone

	// The following are ModRM.mod used as an opcode
	// extension. E.g. opcode 0x83 uses mod==0 to
	// mean ADDL, mod==1 means ORL.
	mod0
	mod1
	mod2
	mod3
	mod4
	mod5
	mod6
	mod7
	mod8
)
