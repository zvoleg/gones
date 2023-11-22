package cpu6502

type Bus6502 interface {
	CpuRead(address uint16) uint8
	CpuWrite(address uint16, data uint8)
}

type Cpu6502 struct {
	a uint8
	x uint8
	y uint8

	pc     uint16
	s      uint8
	status uint8

	opcode uint8
	amAdr  uint16
	amOpr  uint8

	bus          Bus6502
	clockCounter uint
}

func New(bus Bus6502) Cpu6502 {
	return Cpu6502{bus: bus}
}

func (cpu *Cpu6502) Clock() {
	opcode := cpu.readPc()
	instr := instructionTable[opcode]
	cpu.opcode = opcode
	cpu.clockCounter = instr.clocks
	instr.am(cpu)
	instr.handler(cpu)
}

func (cpu *Cpu6502) incrementPc() {
	cpu.pc += 1
}

func (cpu *Cpu6502) readPc() uint8 {
	data := cpu.bus.CpuRead(cpu.pc)
	cpu.incrementPc()
	return data
}

func (cpu *Cpu6502) getFlag(f flag) uint8 {
	bit := cpu.status & 1 << f
	return bit >> f
}

func (cpu *Cpu6502) setFlag(f flag, set bool) {
	var mask uint8 = 1 << f
	if set {
		cpu.status |= mask
	} else {
		cpu.status &= ^mask
	}
}

func (cpu *Cpu6502) push(data uint8) {
	cpu.bus.CpuWrite(0x0100|uint16(cpu.s), data)
	cpu.s -= 1
}

func (cpu *Cpu6502) pop() uint8 {
	cpu.s += 1
	return cpu.bus.CpuRead(0x0100 | uint16(cpu.s))
}
