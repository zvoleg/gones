package cpu6502

import "fmt"

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

	opcode byte
	amAdr  uint16

	bus          Bus6502
	signalLine   chan Signal
	dmaEnabled   bool
	clockCounter uint
}

func New(bus Bus6502, interruptLine chan Signal) Cpu6502 {
	cpu := Cpu6502{bus: bus, signalLine: interruptLine}
	cpu.reset()
	return cpu
}

func (cpu *Cpu6502) Clock() {
	select {
	case signal := <-cpu.signalLine:
		switch signal {
		case Nmi:
			cpu.interrupt(nmiVector)
		case Irq:
			if cpu.getFlag(I) == 0 {
				cpu.interrupt(irqVector)
			}
		case Res:
			cpu.reset()
		case DmaEnable:
			cpu.dmaEnabled = true
		case DmaDisable:
			cpu.dmaEnabled = false
		}
	default:
	}
	if cpu.dmaEnabled {
		return
	}

	opcode := cpu.readPc()
	instr := instructionTable[opcode]
	log := fmt.Sprintf("pc:%04X opcode:%02X\t%s", cpu.pc, opcode, instr.name)
	cpu.opcode = opcode
	cpu.clockCounter = instr.clocks
	instr.am(cpu)
	instr.handler(cpu)
	fmt.Printf("%s addr:%04X\n", log, cpu.amAdr)
	clocks := clock{cpu.clockCounter}
	clocks.waitExecution()
}

func (cpu *Cpu6502) reset() {
	cpu.s = 0xFD // stack pointer decrements by 2
	pcL := cpu.bus.CpuRead(resVector)
	pcH := cpu.bus.CpuRead(resVector + 1)
	cpu.pc = uint16(pcH)<<8 | uint16(pcL)
	cpu.status = 0x34 // set flags U, B, Iconst irqVecL uint16 = 0xFFFE
}

func (cpu *Cpu6502) interrupt(vector uint16) {
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
	clocks := clock{7}
	clocks.waitExecution()
}

func (cpu *Cpu6502) incrementPc() {
	cpu.pc += 1
}

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

func (cpu *Cpu6502) push(data byte) {
	cpu.bus.CpuWrite(0x0100|uint16(cpu.s), data)
	cpu.s -= 1
}

func (cpu *Cpu6502) pop() byte {
	cpu.s += 1
	return cpu.bus.CpuRead(0x0100 | uint16(cpu.s))
}

func (cpu *Cpu6502) fetch() byte {
	return cpu.bus.CpuRead(cpu.amAdr)
}
