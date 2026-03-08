package cpu6502

type Bus6502 interface {
	CpuRead(address uint16) byte
	CpuWrite(address uint16, data byte)
}

const nmiVector uint16 = 0xFFFA
const resVector uint16 = 0xFFFC
const irqVector uint16 = 0xFFFE

type Signal int

const (
	Irq Signal = iota
	Nmi
	Res
	DmaEnable
	DmaDisable
)

type Cpu6502 struct {
	a byte
	x byte
	y byte

	pc     uint16
	s      byte
	status byte

	opcode      byte
	instrData   uint16
	operatorAdr uint16

	bus          Bus6502
	clockCounter uint
}

func New(bus Bus6502) Cpu6502 {
	cpu := Cpu6502{bus: bus}
	cpu.Reset()
	return cpu
}

func (cpu *Cpu6502) SetPC(address uint16) {
	cpu.pc = address
}

func (cpu *Cpu6502) Clock() {
	if cpu.clockCounter == 0 {
		// opcodeAddr := cpu.pc
		opcode := cpu.readPc()
		instr := instructionTable[opcode]
		cpu.opcode = opcode
		cpu.clockCounter = instr.clocks
		cpu.fetch(instr.am.size)
		instr.am.exec(cpu)
		instr.op.exec(cpu)
		// var data byte
		// switch instr.am.name {
		// case "ZP0", "ZPX", "ZPY", "ABS", "ABX", "ABY", "IDX", "IDY":
		// 	if cpu.operatorAdr < 0x2000 || cpu.operatorAdr >= 0x4020 {
		// 		data = cpu.bus.CpuRead(cpu.operatorAdr)
		// 	}
		// }
		// log.Printf("a=%02X x=%02X Y=%02X st=%08b pc=%04X st_ptr=%02X | opcode=%02X %s", cpu.a, cpu.x, cpu.y, cpu.status, opcodeAddr, cpu.s, cpu.opcode, instr.disassebly(cpu.instrData, data))
	} else {
		cpu.clockCounter -= 1
	}
}

func (cpu *Cpu6502) Reset() {
	cpu.s = 0xFD // stack pointer decrements by 2
	pcL := cpu.bus.CpuRead(resVector)
	pcH := cpu.bus.CpuRead(resVector + 1)
	cpu.pc = uint16(pcH)<<8 | uint16(pcL)
	cpu.status = 0x24 // set flags U, I
}

func (cpu *Cpu6502) Interrupt(signal Signal) {
	vector := resVector
	switch signal {
	case Irq:
		if cpu.getFlag(I) != 0 {
			return
		}
		vector = irqVector
	case Nmi:
		vector = nmiVector
	}
	cpu.clockCounter = 7 // interrupt handler clocks
	// save cpu backup on the stack
	pcH := byte(cpu.pc >> 8)
	cpu.push(pcH)
	pcL := byte(cpu.pc)
	cpu.push(pcL)
	cpu.setFlag(B, false)
	cpu.push(cpu.status)
	// prepare cpu to handling interrupt
	cpu.setFlag(I, true)
	pcL = cpu.bus.CpuRead(vector)
	pcH = cpu.bus.CpuRead(vector + 1)
	cpu.pc = uint16(pcH)<<8 | uint16(pcL)
}

func (cpu *Cpu6502) incrementPc() {
	cpu.pc += 1
}

// Reading a byte by current program counter and increments it
func (cpu *Cpu6502) readPc() byte {
	data := cpu.bus.CpuRead(cpu.pc)
	cpu.incrementPc()
	return data
}

func (cpu *Cpu6502) getFlag(f flag) byte {
	bit := cpu.status & (1 << f)
	return bit >> f
}

func (cpu *Cpu6502) setFlag(f flag, set bool) {
	var mask byte = 1 << f
	if set {
		cpu.status |= mask
	} else {
		cpu.status &= ^mask
	}
}

// this function pushes a byte to the stack and executes the stack pointer calculations
func (cpu *Cpu6502) push(data byte) {
	cpu.bus.CpuWrite(0x0100|uint16(cpu.s), data)
	cpu.s -= 1
}

// this function popes a byte to the stack and executes the stack pointer calculations
func (cpu *Cpu6502) pop() byte {
	cpu.s += 1
	return cpu.bus.CpuRead(0x0100 | uint16(cpu.s))
}

// Fetch instruction data in amount of n bytes.
//
// There is assumed that n can't be more than 2
//
// First read byte is represented as LSB and the second read byte is represented as MSB
func (cpu *Cpu6502) fetch(n uint16) {
	var data uint16 = 0
	for i := uint16(0); i < n; i += 1 {
		b := uint16(cpu.readPc())
		b <<= 8 * i
		data |= b
	}
	cpu.instrData = data
}
