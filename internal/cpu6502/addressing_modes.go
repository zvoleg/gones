package cpu6502

func acc(cpu *Cpu6502) {
}

func imm(cpu *Cpu6502) {
	// all data was readed by cpu fetchin instruction data
	// but all instructions fetches data by their address
	cpu.operatorAdr = cpu.pc - 1
}

func abs(cpu *Cpu6502) {
	cpu.operatorAdr = cpu.instrData
}

func zp0(cpu *Cpu6502) {
	cpu.operatorAdr = cpu.instrData
}

func zpx(cpu *Cpu6502) {
	adr := cpu.instrData + uint16(cpu.x)
	adr = adr & 0xFF
	cpu.operatorAdr = adr
}

func zpy(cpu *Cpu6502) {
	adr := cpu.instrData + uint16(cpu.y)
	adr = adr & 0xFF
	cpu.operatorAdr = adr
}

func abx(cpu *Cpu6502) {
	baseAdr := cpu.instrData
	cpu.operatorAdr = baseAdr + uint16(cpu.x)
	if cpu.operatorAdr&0xFF00 != baseAdr&0xFF00 {
		cpu.clockCounter += 1
	}
}

func aby(cpu *Cpu6502) {
	baseAdr := cpu.instrData
	cpu.operatorAdr = baseAdr + uint16(cpu.y)
	if cpu.operatorAdr&0xFF00 != baseAdr&0xFF00 {
		cpu.clockCounter += 1
	}
}

func imp(cpu *Cpu6502) {
	// operands are contains directly in the instruction
}

func rel(cpu *Cpu6502) {
	offset := cpu.instrData
	if offset&0x80 != 0 {
		offset = 0xFF00 | offset
	}
	cpu.operatorAdr = cpu.pc + offset
}

func idx(cpu *Cpu6502) {
	zAdr := cpu.instrData
	zAdr = (zAdr + uint16(cpu.x)) & 0xFF // Discard cary
	adrL := cpu.bus.CpuRead(zAdr)
	adrH := cpu.bus.CpuRead((zAdr + 1) & 0xFF) // Both memory locations must be in page zero
	cpu.operatorAdr = uint16(adrH)<<8 | uint16(adrL)
}

func idy(cpu *Cpu6502) {
	zAdr := cpu.instrData
	adrL := cpu.bus.CpuRead(zAdr)
	adrH := cpu.bus.CpuRead((zAdr + 1) & 0xFF) // Both memory locations must be in page zero
	baseAdr := uint16(adrH)<<8 | uint16(adrL)
	cpu.operatorAdr = baseAdr + uint16(cpu.y)
	if cpu.operatorAdr&0xFF00 != baseAdr&0xFF {
		cpu.clockCounter += 1
	}
}

func ind(cpu *Cpu6502) { // JMP[IND] only
	iAdr := cpu.instrData
	adrL := cpu.bus.CpuRead(iAdr)
	iAdrH := iAdr & 0xFF00
	iAdrL := (iAdr + 1) & 0x00FF // Can't cros the memory page
	nextIAdr := iAdrH | iAdrL
	adrH := cpu.bus.CpuRead(nextIAdr)
	cpu.operatorAdr = uint16(adrH)<<8 | uint16(adrL)
}
