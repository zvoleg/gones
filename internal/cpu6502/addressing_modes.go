package cpu6502

func acc(cpu *Cpu6502) {
	cpu.amOpr = cpu.a
}

func imm(cpu *Cpu6502) {
	data := cpu.readPc()
	cpu.amOpr = data
}

func abs(cpu *Cpu6502) {
	adrL := cpu.readPc()
	adrH := cpu.readPc()
	cpu.amAdr = uint16(adrH)<<8 | uint16(adrL)
	cpu.amOpr = cpu.bus.CpuRead(cpu.amAdr)
}

func zp0(cpu *Cpu6502) {
	cpu.amAdr = uint16(cpu.readPc())
	cpu.amOpr = cpu.bus.CpuRead(cpu.amAdr)
}

func zpx(cpu *Cpu6502) {
	adr := uint16(cpu.readPc()) + uint16(cpu.x)
	adr = adr & 0xFF
	cpu.amAdr = adr
	cpu.amOpr = cpu.bus.CpuRead(cpu.amAdr)
}

func zpy(cpu *Cpu6502) {
	adr := uint16(cpu.readPc()) + uint16(cpu.y)
	adr = adr & 0xFF
	cpu.amAdr = adr
	cpu.amOpr = cpu.bus.CpuRead(cpu.amAdr)
}

func abx(cpu *Cpu6502) {
	adrL := cpu.readPc()
	adrH := cpu.readPc()
	cpu.amAdr = (uint16(adrH)<<8 | uint16(adrL)) + uint16(cpu.x)
	cpu.amOpr = cpu.bus.CpuRead(cpu.amAdr)
}

func aby(cpu *Cpu6502) {
	adrL := cpu.readPc()
	adrH := cpu.readPc()
	cpu.amAdr = (uint16(adrH)<<8 | uint16(adrL)) + uint16(cpu.y)
	cpu.amOpr = cpu.bus.CpuRead(cpu.amAdr)
}

func imp(cpu *Cpu6502) {
	// operands are contains directly in the instruction
}

func rel(cpu *Cpu6502) {
	offset := cpu.readPc()
	cpu.amAdr = cpu.pc + uint16(offset)
}

func idx(cpu *Cpu6502) {
	zAdr := cpu.readPc()
	zAdr = zAdr + cpu.x
	adrL := cpu.bus.CpuRead(uint16(zAdr))
	adrH := cpu.bus.CpuRead(uint16(zAdr + 1))
	cpu.amAdr = uint16(adrH)<<8 | uint16(adrL)
	cpu.amOpr = cpu.bus.CpuRead(cpu.amAdr)
}

func idy(cpu *Cpu6502) {
	zAdr := cpu.readPc()
	zAdr = zAdr + cpu.y
	adrL := cpu.bus.CpuRead(uint16(zAdr))
	adrH := cpu.bus.CpuRead(uint16(zAdr + 1))
	cpu.amAdr = uint16(adrH)<<8 | uint16(adrL)
	cpu.amOpr = cpu.bus.CpuRead(cpu.amAdr)
}

func ind(cpu *Cpu6502) { // JMP[IND] only
	iAdrL := cpu.readPc()
	iAdrH := cpu.readPc()
	iAdr := uint16(iAdrH)<<8 | uint16(iAdrL)
	adrL := cpu.bus.CpuRead(iAdr)
	adrH := cpu.bus.CpuRead(iAdr + 1)
	cpu.amAdr = uint16(adrH)<<8 | uint16(adrL)
}
