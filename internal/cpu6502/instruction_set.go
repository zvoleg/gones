package cpu6502

const lslAcc uint8 = 0x4A
const aslAcc uint8 = 0x0A
const rolAcc uint8 = 0x2A
const rorAcc uint8 = 0x6A

func add(a, b uint8) (res uint8, carry bool, overflow bool) {
	sum := uint16(a) + uint16(b)

	res = uint8(sum)
	carry = sum&0x100 != 0
	overflow = (a&0x80 == b&0x80) && res&0x80 != a&0x80
	return res, carry, overflow
}

func branchProcess(cpu *Cpu6502) {
	cpu.clockCounter += 1
	if cpu.amAdr&0xFF00 != cpu.pc&0xFF00 { // comparing the high order of bits
		cpu.clockCounter += 1
	}
	cpu.pc = cpu.amAdr
}

func adc(cpu *Cpu6502) { // decimal mod is not implemented
	sum, carrySum, overflSum := add(cpu.a, cpu.amOpr)
	res, carryCarry, overflCarry := add(sum, cpu.getFlag(C))
	cpu.a = res
	cpu.setFlag(C, carrySum || carryCarry)
	cpu.setFlag(V, overflSum || overflCarry)
	cpu.setFlag(N, res&0x80 != 0)
	cpu.setFlag(Z, res == 0)
}

func and(cpu *Cpu6502) {
	res := cpu.a & cpu.amOpr
	cpu.a = res
	cpu.setFlag(Z, res == 0)
	cpu.setFlag(N, res&0x80 != 0)
}

func asl(cpu *Cpu6502) {
	popedBit := cpu.amOpr >> 7
	res := cpu.amOpr << 1
	if cpu.opcode == aslAcc {
		cpu.a = res
	} else {
		cpu.bus.CpuWrite(cpu.amAdr, res)
	}
	cpu.setFlag(C, popedBit == 1)
	cpu.setFlag(N, res&0x80 == 1)
	cpu.setFlag(Z, res == 0)
}

func bcc(cpu *Cpu6502) {
	if cpu.getFlag(C) == 0 {
		branchProcess(cpu)
	}
}

func bcs(cpu *Cpu6502) {
	if cpu.getFlag(C) == 1 {
		branchProcess(cpu)
	}
}

func beq(cpu *Cpu6502) {
	if cpu.getFlag(Z) == 1 {
		branchProcess(cpu)
	}
}

func bit(cpu *Cpu6502) {
	res := cpu.a & cpu.amOpr
	cpu.setFlag(N, res&0x80 != 0)
	cpu.setFlag(V, res&0x40 != 0)
	cpu.setFlag(Z, res == 0)
}

func bmi(cpu *Cpu6502) {
	if cpu.getFlag(N) == 1 {
		branchProcess(cpu)
	}
}

func bne(cpu *Cpu6502) {
	if cpu.getFlag(Z) == 0 {
		branchProcess(cpu)
	}
}

func bpl(cpu *Cpu6502) {
	if cpu.getFlag(N) == 0 {
		branchProcess(cpu)
	}
}

func brk(cpu *Cpu6502) {
	pcH := uint8(cpu.pc >> 8)
	cpu.push(pcH)
	pcL := uint8(cpu.pc)
	cpu.push(pcL)
	cpu.setFlag(I, true)
	cpu.setFlag(B, true)
	cpu.push(cpu.status)

	pcL = cpu.bus.CpuRead(irqVector)
	pcH = cpu.bus.CpuRead(irqVector + 1)
	cpu.pc = uint16(pcH)<<8 | uint16(pcL)
}

func bvc(cpu *Cpu6502) {
	if cpu.getFlag(V) == 0 {
		branchProcess(cpu)
	}
}

func bvs(cpu *Cpu6502) {
	if cpu.getFlag(V) == 1 {
		branchProcess(cpu)
	}
}

func clc(cpu *Cpu6502) {
	cpu.setFlag(C, false)
}

func cld(cpu *Cpu6502) {
	cpu.setFlag(D, false)
}

func cli(cpu *Cpu6502) {
	cpu.setFlag(I, false)
}

func clv(cpu *Cpu6502) {
	cpu.setFlag(V, false)
}

func cmp(cpu *Cpu6502) {
	res := cpu.a - cpu.amOpr
	cpu.setFlag(Z, res == 0)
	cpu.setFlag(N, res&0x80 != 0)
	cpu.setFlag(C, cpu.amOpr <= cpu.a)
}

func cpx(cpu *Cpu6502) {
	res := cpu.x - cpu.amOpr
	cpu.setFlag(Z, res == 0)
	cpu.setFlag(N, res&0x80 != 0)
	cpu.setFlag(C, cpu.amOpr <= cpu.x)
}

func cpy(cpu *Cpu6502) {
	res := cpu.y - cpu.amOpr
	cpu.setFlag(Z, res == 0)
	cpu.setFlag(N, res&0x80 != 0)
	cpu.setFlag(C, cpu.amOpr <= cpu.y)
}

func dec(cpu *Cpu6502) {
	res := cpu.amOpr - 1
	cpu.bus.CpuWrite(cpu.amAdr, res)
	cpu.setFlag(N, res&0x80 != 0)
	cpu.setFlag(Z, res == 0)
}

func dex(cpu *Cpu6502) {
	cpu.x -= 1
	cpu.setFlag(N, cpu.x&0x80 != 0)
	cpu.setFlag(Z, cpu.x == 0)
}

func dey(cpu *Cpu6502) {
	cpu.y -= 1
	cpu.setFlag(N, cpu.y&0x80 != 0)
	cpu.setFlag(Z, cpu.y == 0)
}

func eor(cpu *Cpu6502) {
	res := cpu.a ^ cpu.amOpr
	cpu.a = res
	cpu.setFlag(Z, res == 0)
	cpu.setFlag(N, res&0x80 != 0)
}

func inc(cpu *Cpu6502) {
	res := cpu.amOpr + 1
	cpu.bus.CpuWrite(cpu.amAdr, res)
	cpu.setFlag(N, res&0x80 != 0)
	cpu.setFlag(Z, res == 0)
}

func inx(cpu *Cpu6502) {
	cpu.x += 1
	cpu.setFlag(N, cpu.x&0x80 != 0)
	cpu.setFlag(Z, cpu.x == 0)
}

func iny(cpu *Cpu6502) {
	cpu.y += 1
	cpu.setFlag(N, cpu.y&0x80 != 0)
	cpu.setFlag(Z, cpu.y == 0)
}

func jmp(cpu *Cpu6502) {
	cpu.pc = cpu.amAdr
}

func jsr(cpu *Cpu6502) {
	pcH := uint8(cpu.pc >> 8)
	cpu.push(pcH)
	pcL := uint8(cpu.pc)
	cpu.push(pcL)
	cpu.pc = cpu.amAdr
}

func lda(cpu *Cpu6502) {
	cpu.a = cpu.amOpr
	cpu.setFlag(N, cpu.a&0x80 != 0)
	cpu.setFlag(Z, cpu.a == 0)
}

func ldx(cpu *Cpu6502) {
	cpu.x = cpu.amOpr
	cpu.setFlag(Z, cpu.amOpr == 0)
	cpu.setFlag(N, cpu.amOpr&0x80 != 0)
}

func ldy(cpu *Cpu6502) {
	cpu.y = cpu.amOpr
	cpu.setFlag(Z, cpu.amOpr == 0)
	cpu.setFlag(N, cpu.amOpr&0x80 != 0)
}

func lsr(cpu *Cpu6502) {
	popedBit := cpu.amOpr & 1
	res := cpu.amOpr >> 1
	if cpu.opcode == lslAcc {
		cpu.a = res
	} else {
		cpu.bus.CpuWrite(cpu.amAdr, res)
	}
	cpu.setFlag(C, popedBit == 1)
	cpu.setFlag(N, false)
	cpu.setFlag(Z, res == 0)
}

func nop(cpu *Cpu6502) {

}

func ora(cpu *Cpu6502) {
	res := cpu.a | cpu.amOpr
	cpu.a = res
	cpu.setFlag(Z, res == 0)
	cpu.setFlag(N, res&0x80 != 0)
}

func pha(cpu *Cpu6502) {
	cpu.push(cpu.a)
}

func php(cpu *Cpu6502) {
	cpu.push(cpu.status)
}

func pla(cpu *Cpu6502) {
	cpu.a = cpu.pop()
	cpu.setFlag(N, cpu.a&0x80 != 0)
	cpu.setFlag(Z, cpu.a == 0)
}

func plp(cpu *Cpu6502) {
	cpu.status = cpu.pop()
}

func rol(cpu *Cpu6502) {
	popedBit := cpu.amOpr >> 7
	res := cpu.amOpr << 1
	res = res | cpu.getFlag(C)
	if cpu.opcode == rolAcc {
		cpu.a = res
	} else {
		cpu.bus.CpuWrite(cpu.amAdr, res)
	}
	cpu.setFlag(C, popedBit == 1)
	cpu.setFlag(N, res&0x80 == 1)
	cpu.setFlag(Z, res == 0)
}

func ror(cpu *Cpu6502) {
	popedBit := cpu.amOpr & 1
	res := cpu.amOpr >> 1
	res = res | (cpu.getFlag(C) << 7)
	if cpu.opcode == rorAcc {
		cpu.a = res
	} else {
		cpu.bus.CpuWrite(cpu.amAdr, res)
	}
	cpu.setFlag(C, popedBit == 1)
	cpu.setFlag(N, res&0x80 == 1)
	cpu.setFlag(Z, res == 0)
}

func rti(cpu *Cpu6502) {
	cpu.status = cpu.pop()
	cpu.setFlag(B, false)
	pcL := cpu.pop()
	pcH := cpu.pop()
	cpu.pc = uint16(pcH)<<8 | uint16(pcL)
}

func rts(cpu *Cpu6502) {
	pcL := cpu.pop()
	pcH := cpu.pop()
	cpu.pc = uint16(pcH)<<8 | uint16(pcL)
}

func sbc(cpu *Cpu6502) {
	sub := ^cpu.amOpr + cpu.getFlag(C)
	res, _, overflow := add(cpu.a, sub)
	cpu.a = res
	cpu.setFlag(C, int8(res) >= 0)
	cpu.setFlag(V, overflow)
	cpu.setFlag(N, res&0x80 != 0)
	cpu.setFlag(Z, res == 0)
}

func sec(cpu *Cpu6502) {
	cpu.setFlag(C, true)
}

func sed(cpu *Cpu6502) {
	cpu.setFlag(D, true)
}

func sei(cpu *Cpu6502) {
	cpu.setFlag(I, true)
}

func sta(cpu *Cpu6502) {
	cpu.bus.CpuWrite(cpu.amAdr, cpu.a)
}

func stx(cpu *Cpu6502) {
	cpu.bus.CpuWrite(cpu.amAdr, cpu.x)
}

func sty(cpu *Cpu6502) {
	cpu.bus.CpuWrite(cpu.amAdr, cpu.y)
}

func tax(cpu *Cpu6502) {
	cpu.x = cpu.a
	cpu.setFlag(N, cpu.x&0x80 != 0)
	cpu.setFlag(Z, cpu.x == 0)
}

func tay(cpu *Cpu6502) {
	cpu.y = cpu.a
	cpu.setFlag(N, cpu.y&0x80 != 0)
	cpu.setFlag(Z, cpu.y == 0)
}

func tsx(cpu *Cpu6502) {
	cpu.x = cpu.s
	cpu.setFlag(N, cpu.x&0x80 != 0)
	cpu.setFlag(Z, cpu.x == 0)
}

func txa(cpu *Cpu6502) {
	cpu.a = cpu.x
	cpu.setFlag(N, cpu.a&0x80 != 0)
	cpu.setFlag(Z, cpu.a == 0)
}

func txs(cpu *Cpu6502) {
	cpu.s = cpu.x
}

func tya(cpu *Cpu6502) {
	cpu.a = cpu.y
	cpu.setFlag(N, cpu.a&0x80 != 0)
	cpu.setFlag(Z, cpu.a == 0)
}

func xep(cpu *Cpu6502) {

}
