package cpu6502

func acc(cpu *Cpu6502) {
}

func imm(cpu *Cpu6502) {
	cpu.amAdr = cpu.pc
	cpu.incrementPc()
}

func abs(cpu *Cpu6502) {
	adrL := cpu.readPc()
	adrH := cpu.readPc()
	cpu.amAdr = uint16(adrH)<<8 | uint16(adrL)
}

func zp0(cpu *Cpu6502) {
	cpu.amAdr = uint16(cpu.readPc())
}

func zpx(cpu *Cpu6502) {
	adr := uint16(cpu.readPc()) + uint16(cpu.x)
	adr = adr & 0xFF
	cpu.amAdr = adr
}

func zpy(cpu *Cpu6502) {
	adr := uint16(cpu.readPc()) + uint16(cpu.y)
	adr = adr & 0xFF
	cpu.amAdr = adr
}

func abx(cpu *Cpu6502) {
	adrL := cpu.readPc()
	adrH := cpu.readPc()
	baseAdr := uint16(adrH)<<8 | uint16(adrL)
	cpu.amAdr = baseAdr + uint16(cpu.x)
	if byte(baseAdr>>8) != adrH {
		cpu.clockCounter += 1
	}
}

func aby(cpu *Cpu6502) {
	adrL := cpu.readPc()
	adrH := cpu.readPc()
	baseAdr := uint16(adrH)<<8 | uint16(adrL)
	cpu.amAdr = baseAdr + uint16(cpu.y)
	if byte(baseAdr>>8) != adrH {
		cpu.clockCounter += 1
	}
}

func imp(cpu *Cpu6502) {
	// operands are contains directly in the instruction
}

func rel(cpu *Cpu6502) {
	offset := uint16(cpu.readPc())
	if offset&0x80 != 0 {
		offset = 0xFF00 | offset
	}
	cpu.amAdr = cpu.pc + uint16(offset)
}

func idx(cpu *Cpu6502) {
	zAdr := cpu.readPc()
	zAdr = zAdr + cpu.x
	adrL := cpu.bus.CpuRead(uint16(zAdr))
	adrH := cpu.bus.CpuRead(uint16(zAdr + 1))
	cpu.amAdr = uint16(adrH)<<8 | uint16(adrL)
}

func idy(cpu *Cpu6502) {
	zAdr := cpu.readPc()
	adrL := cpu.bus.CpuRead(uint16(zAdr))
	adrH := cpu.bus.CpuRead(uint16(zAdr + 1))
	baseAdr := uint16(adrH)<<8 | uint16(adrL)
	cpu.amAdr = baseAdr + uint16(cpu.y)
	if byte(cpu.amAdr>>8) != adrH {
		cpu.clockCounter += 1
	}
}

func ind(cpu *Cpu6502) { // JMP[IND] only
	iAdrL := cpu.readPc()
	iAdrH := cpu.readPc()
	iAdr := uint16(iAdrH)<<8 | uint16(iAdrL)
	adrL := cpu.bus.CpuRead(iAdr)
	adrH := cpu.bus.CpuRead(iAdr + 1)
	cpu.amAdr = uint16(adrH)<<8 | uint16(adrL)
}
