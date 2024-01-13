package cpu6502

import (
	"fmt"
	"strings"
)

type operation struct {
	name    string
	handler func(*Cpu6502)
}

func (o *operation) exec(cpu *Cpu6502) {
	o.handler(cpu)
}

type addressingMode struct {
	name        string
	handler     func(*Cpu6502)
	size        uint16 // 0, 1 ,2
	dataFormate string
}

func (am *addressingMode) exec(cpu *Cpu6502) {
	am.handler(cpu)
}

func (am *addressingMode) dataRepresentation(data uint16) string {
	if am.size > 0 {
		return fmt.Sprintf(am.dataFormate, data)
	} else {
		return am.dataFormate
	}
}

type instruction struct {
	opcode byte
	op     operation
	am     addressingMode
	clocks uint
}

func inst(opcode byte, op operation, am addressingMode, clocks uint) instruction {
	return instruction{
		opcode,
		op,
		am,
		clocks,
	}
}

func (i *instruction) disassebly(data uint16) string {
	return strings.Join([]string{i.op.name, i.am.dataRepresentation(data)}, " ")
}

var (
	ADC = operation{"ADC", adc}
	AND = operation{"AND", and}
	ASL = operation{"ASL", asl}
	BCC = operation{"BCC", bcc}
	BCS = operation{"BCS", bcs}
	BEQ = operation{"BEQ", beq}
	BIT = operation{"BIT", bit}
	BMI = operation{"BMI", bmi}
	BNE = operation{"BNE", bne}
	BPL = operation{"BPL", bpl}
	BRK = operation{"BRK", brk}
	BVC = operation{"BVC", bvc}
	BVS = operation{"BVS", bvs}
	CLC = operation{"CLC", clc}
	CLD = operation{"CLD", cld}
	CLI = operation{"CLI", cli}
	CLV = operation{"CLV", clv}
	CMP = operation{"CMP", cmp}
	CPX = operation{"CPX", cpx}
	CPY = operation{"CPY", cpy}
	DEC = operation{"DEC", dec}
	DEX = operation{"DEX", dex}
	DEY = operation{"DEY", dey}
	EOR = operation{"EOR", eor}
	INC = operation{"INC", inc}
	INX = operation{"INX", inx}
	INY = operation{"INY", iny}
	JMP = operation{"JMP", jmp}
	JSR = operation{"JSR", jsr}
	LDA = operation{"LDA", lda}
	LDX = operation{"LDX", ldx}
	LDY = operation{"LDY", ldy}
	LSR = operation{"LSR", lsr}
	NOP = operation{"NOP", nop}
	ORA = operation{"ORA", ora}
	PHA = operation{"PHA", pha}
	PHP = operation{"PHP", php}
	PLA = operation{"PLA", pla}
	PLP = operation{"PLP", plp}
	ROL = operation{"ROL", rol}
	ROR = operation{"ROR", ror}
	RTI = operation{"RTI", rti}
	RTS = operation{"RTS", rts}
	SBC = operation{"SBC", sbc}
	SEC = operation{"SEC", sec}
	SED = operation{"SED", sed}
	SEI = operation{"SEI", sei}
	STA = operation{"STA", sta}
	STX = operation{"STX", stx}
	STY = operation{"STY", sty}
	TAX = operation{"TAX", tax}
	TAY = operation{"TAY", tay}
	TSX = operation{"TSX", tsx}
	TXA = operation{"TXA", txa}
	TXS = operation{"TXS", txs}
	TYA = operation{"TYA", tya}
	XEP = operation{"XEP", xep}
)

var (
	ACC = addressingMode{"ACC", acc, 0, "A"}
	IMM = addressingMode{"IMM", imm, 1, "#$%02X"}
	ABS = addressingMode{"ABS", abs, 2, "$%04X"}
	ZP0 = addressingMode{"ZP0", zp0, 1, "$%02X"}
	ZPX = addressingMode{"ZPX", zpx, 1, "$%02X, X"}
	ZPY = addressingMode{"ZPY", zpy, 1, "$%02X, Y"}
	ABX = addressingMode{"ABX", abx, 2, "$%04X, X"}
	ABY = addressingMode{"ABY", aby, 2, "$%04X, Y"}
	IMP = addressingMode{"IMP", imp, 0, ""}
	REL = addressingMode{"REL", rel, 1, "$%02X"}
	IDX = addressingMode{"IDX", idx, 1, "($%02X, X)"}
	IDY = addressingMode{"IDY", idy, 1, "($%02X), Y"}
	IND = addressingMode{"IND", ind, 2, "($%04X)"}
)

var instructionTable []instruction = []instruction{
	inst(0x00, BRK, IMP, 7), inst(0x01, ORA, IDX, 6), inst(0x02, XEP, IMP, 0), inst(0x03, XEP, IMP, 0), inst(0x04, XEP, IMP, 0), inst(0x05, ORA, ZP0, 3), inst(0x06, ASL, ZP0, 5), inst(0x07, XEP, IMP, 0), inst(0x08, PHP, IMP, 3), inst(0x09, ORA, IMM, 2), inst(0x0A, ASL, ACC, 2), inst(0x0B, XEP, IMP, 0), inst(0x0C, XEP, IMP, 0), inst(0x0D, ORA, ABS, 4), inst(0x0E, ASL, ABS, 6), inst(0x0F, XEP, IMP, 0),
	inst(0x10, BPL, REL, 2), inst(0x11, ORA, IDY, 5), inst(0x12, XEP, IMP, 0), inst(0x13, XEP, IMP, 0), inst(0x14, XEP, IMP, 0), inst(0x15, ORA, ZPX, 4), inst(0x16, ASL, ZPX, 6), inst(0x17, XEP, IMP, 0), inst(0x18, CLC, IMP, 2), inst(0x19, ORA, ABY, 4), inst(0x1A, XEP, IMP, 0), inst(0x1B, XEP, IMP, 0), inst(0x1C, XEP, IMP, 0), inst(0x1D, ORA, ABX, 4), inst(0x1E, ASL, ABX, 7), inst(0x1F, XEP, IMP, 0),
	inst(0x20, JSR, ABS, 6), inst(0x21, AND, IDX, 6), inst(0x22, XEP, IMP, 0), inst(0x23, XEP, IMP, 0), inst(0x24, BIT, ZP0, 3), inst(0x25, AND, ZP0, 3), inst(0x26, ROL, ZP0, 5), inst(0x27, XEP, IMP, 0), inst(0x28, PLP, IMP, 4), inst(0x29, AND, IMM, 2), inst(0x2A, ROL, ACC, 2), inst(0x2B, XEP, IMP, 0), inst(0x2C, BIT, ABS, 4), inst(0x2D, AND, ABS, 4), inst(0x2E, ROL, ABS, 6), inst(0x2F, XEP, IMP, 0),
	inst(0x30, BMI, REL, 2), inst(0x31, AND, IDY, 5), inst(0x32, XEP, IMP, 0), inst(0x33, XEP, IMP, 0), inst(0x34, XEP, IMP, 0), inst(0x35, AND, ZPX, 4), inst(0x36, ROL, ZPX, 6), inst(0x37, XEP, IMP, 0), inst(0x38, SEC, IMP, 2), inst(0x39, AND, ABY, 4), inst(0x3A, XEP, IMP, 0), inst(0x3B, XEP, IMP, 0), inst(0x3C, XEP, IMP, 0), inst(0x3D, AND, ABX, 4), inst(0x3E, ROL, ABX, 7), inst(0x3F, XEP, IMP, 0),
	inst(0x40, RTI, IMP, 6), inst(0x41, EOR, IDX, 6), inst(0x42, XEP, IMP, 0), inst(0x43, XEP, IMP, 0), inst(0x44, XEP, IMP, 0), inst(0x45, EOR, ZP0, 3), inst(0x46, LSR, ZP0, 5), inst(0x47, XEP, IMP, 0), inst(0x48, PHA, IMP, 3), inst(0x49, EOR, IMM, 2), inst(0x4A, LSR, ACC, 2), inst(0x4B, XEP, IMP, 0), inst(0x4C, JMP, ABS, 3), inst(0x4D, EOR, ABS, 4), inst(0x4E, LSR, ABS, 6), inst(0x4F, XEP, IMP, 0),
	inst(0x50, BVC, REL, 2), inst(0x51, EOR, IDY, 5), inst(0x52, XEP, IMP, 0), inst(0x53, XEP, IMP, 0), inst(0x54, XEP, IMP, 0), inst(0x55, EOR, ZPX, 4), inst(0x56, LSR, ZPX, 6), inst(0x57, XEP, IMP, 0), inst(0x58, CLI, IMP, 2), inst(0x59, EOR, ABY, 4), inst(0x5A, XEP, IMP, 0), inst(0x5B, XEP, IMP, 0), inst(0x5C, XEP, IMP, 0), inst(0x5D, EOR, ABX, 4), inst(0x5E, LSR, ABX, 7), inst(0x5F, XEP, IMP, 0),
	inst(0x60, RTS, IMP, 6), inst(0x61, ADC, IDX, 6), inst(0x62, XEP, IMP, 0), inst(0x63, XEP, IMP, 0), inst(0x64, XEP, IMP, 0), inst(0x65, ADC, ZP0, 3), inst(0x66, ROR, ZP0, 5), inst(0x67, XEP, IMP, 0), inst(0x68, PLA, IMP, 4), inst(0x69, ADC, IMM, 2), inst(0x6A, ROR, ACC, 2), inst(0x6B, XEP, IMP, 0), inst(0x6C, JMP, IND, 5), inst(0x6D, ADC, ABS, 4), inst(0x6E, ROR, ABS, 6), inst(0x6F, XEP, IMP, 0),
	inst(0x70, BVS, REL, 2), inst(0x71, ADC, IDY, 5), inst(0x72, XEP, IMP, 0), inst(0x73, XEP, IMP, 0), inst(0x74, XEP, IMP, 0), inst(0x75, ADC, ZPX, 4), inst(0x76, ROR, ZPX, 6), inst(0x77, XEP, IMP, 0), inst(0x78, SEI, IMP, 2), inst(0x79, ADC, ABY, 4), inst(0x7A, XEP, IMP, 0), inst(0x7B, XEP, IMP, 0), inst(0x7C, XEP, IMP, 0), inst(0x7D, ADC, ABX, 4), inst(0x7E, ROR, ABX, 7), inst(0x7F, XEP, IMP, 0),
	inst(0x80, XEP, IMP, 0), inst(0x81, STA, IDX, 6), inst(0x82, XEP, IMP, 0), inst(0x83, XEP, IMP, 0), inst(0x84, STY, ZP0, 3), inst(0x85, STA, ZP0, 3), inst(0x86, STX, ZP0, 3), inst(0x87, XEP, IMP, 0), inst(0x88, DEY, IMP, 2), inst(0x89, XEP, IMP, 0), inst(0x8A, TXA, IMP, 2), inst(0x8B, XEP, IMP, 0), inst(0x8C, STY, ABS, 4), inst(0x8D, STA, ABS, 4), inst(0x8E, STX, ABS, 4), inst(0x8F, XEP, IMP, 0),
	inst(0x90, BCC, REL, 2), inst(0x91, STA, IDY, 6), inst(0x92, XEP, IMP, 0), inst(0x93, XEP, IMP, 0), inst(0x94, STY, ZPX, 4), inst(0x95, STA, ZPX, 4), inst(0x96, STX, ZPY, 4), inst(0x97, XEP, IMP, 0), inst(0x98, TYA, IMP, 2), inst(0x99, STA, ABY, 5), inst(0x9A, TXS, IMP, 2), inst(0x9B, XEP, IMP, 0), inst(0x9C, XEP, IMP, 0), inst(0x9D, STA, ABX, 5), inst(0x9E, XEP, IMP, 0), inst(0x9F, XEP, IMP, 0),
	inst(0xA0, LDY, IMM, 2), inst(0xA1, LDA, IDX, 6), inst(0xA2, LDX, IMM, 2), inst(0xA3, XEP, IMP, 0), inst(0xA4, LDY, ZP0, 3), inst(0xA5, LDA, ZP0, 3), inst(0xA6, LDX, ZP0, 3), inst(0xA7, XEP, IMP, 0), inst(0xA8, TAY, IMP, 2), inst(0xA9, LDA, IMM, 2), inst(0xAA, TAX, IMP, 2), inst(0xAB, XEP, IMP, 0), inst(0xAC, LDY, ABS, 4), inst(0xAD, LDA, ABS, 4), inst(0xAE, LDX, ABS, 4), inst(0xAF, XEP, IMP, 0),
	inst(0xB0, BCS, REL, 2), inst(0xB1, LDA, IDY, 5), inst(0xB2, XEP, IMP, 0), inst(0xB3, XEP, IMP, 0), inst(0xB4, LDY, ZPX, 4), inst(0xB5, LDA, ZPX, 4), inst(0xB6, LDX, ZPY, 4), inst(0xB7, XEP, IMP, 0), inst(0xB8, CLV, IMP, 2), inst(0xB9, LDA, ABY, 4), inst(0xBA, TSX, IMP, 2), inst(0xBB, XEP, IMP, 0), inst(0xBC, LDY, ABX, 4), inst(0xBD, LDA, ABX, 4), inst(0xBE, LDX, ABY, 4), inst(0xBF, XEP, IMP, 0),
	inst(0xC0, CPY, IMM, 2), inst(0xC1, CMP, IDX, 6), inst(0xC2, XEP, IMP, 0), inst(0xC3, XEP, IMP, 0), inst(0xC4, CPY, ZP0, 3), inst(0xC5, CMP, ZP0, 3), inst(0xC6, DEC, ZP0, 5), inst(0xC7, XEP, IMP, 0), inst(0xC8, INY, IMP, 2), inst(0xC9, CMP, IMM, 2), inst(0xCA, DEX, IMP, 2), inst(0xCB, XEP, IMP, 0), inst(0xCC, CPY, ABS, 4), inst(0xCD, CMP, ABS, 4), inst(0xCE, DEC, ABS, 6), inst(0xCF, XEP, IMP, 0),
	inst(0xD0, BNE, REL, 2), inst(0xD1, CMP, IDY, 5), inst(0xD2, XEP, IMP, 0), inst(0xD3, XEP, IMP, 0), inst(0xD4, XEP, IMP, 0), inst(0xD5, CMP, ZPX, 4), inst(0xD6, DEC, ZPX, 6), inst(0xD7, XEP, IMP, 0), inst(0xD8, CLD, IMP, 2), inst(0xD9, CMP, ABY, 4), inst(0xDA, XEP, IMP, 0), inst(0xDB, XEP, IMP, 0), inst(0xDC, XEP, IMP, 0), inst(0xDD, CMP, ABX, 4), inst(0xDE, DEC, ABX, 7), inst(0xDF, XEP, IMP, 0),
	inst(0xE0, CPX, IMM, 2), inst(0xE1, SBC, IDX, 6), inst(0xE2, XEP, IMP, 0), inst(0xE3, XEP, IMP, 0), inst(0xE4, CPX, ZP0, 3), inst(0xE5, SBC, ZP0, 3), inst(0xE6, INC, ZP0, 5), inst(0xE7, XEP, IMP, 0), inst(0xE8, INX, IMP, 2), inst(0xE9, SBC, IMM, 2), inst(0xEA, NOP, IMP, 2), inst(0xEB, XEP, IMP, 0), inst(0xEC, CPX, ABS, 4), inst(0xED, SBC, ABS, 4), inst(0xEE, INC, ABS, 6), inst(0xEF, XEP, IMP, 0),
	inst(0xF0, BEQ, REL, 2), inst(0xF1, SBC, IDY, 5), inst(0xF2, XEP, IMP, 0), inst(0xF3, XEP, IMP, 0), inst(0xF4, XEP, IMP, 0), inst(0xF5, SBC, ZPX, 4), inst(0xF6, INC, ZPX, 6), inst(0xF7, XEP, IMP, 0), inst(0xF8, SED, IMP, 2), inst(0xF9, SBC, ABY, 4), inst(0xFA, XEP, IMP, 0), inst(0xFB, XEP, IMP, 0), inst(0xFC, XEP, IMP, 0), inst(0xFD, SBC, ABX, 4), inst(0xFE, INC, ABX, 7), inst(0xFF, XEP, IMP, 0),
}
