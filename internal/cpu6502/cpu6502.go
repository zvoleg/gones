package cpu6502

import "fmt"

type Bus6502 interface {
	CpuRead(address uint16) uint8
	CpuWrite(address uint16, data uint8)
}

const nmiVector uint16 = 0xFFFA
const resVector uint16 = 0xFFFC
const irqVector uint16 = 0xFFFE

type Interrupt int

const (
	Irq Interrupt = iota
	Nmi
	Res
)

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

	bus             Bus6502
	interruptSignal chan Interrupt
	clockCounter    uint
}

func New(bus Bus6502, interruptLine chan Interrupt) Cpu6502 {
	cpu := Cpu6502{bus: bus, interruptSignal: interruptLine}
	cpu.reset()
	return cpu
}

func (cpu *Cpu6502) Clock() {
	if cpu.getFlag(I) == 1 {
		select {
		case interrupt := <-cpu.interruptSignal:
			switch interrupt {
			case Nmi:
				cpu.interrupt(nmiVector)
			case Irq:
				cpu.interrupt(irqVector)
			case Res:
				cpu.reset()
			}
			return
		default:
		}
	}
	opcode := cpu.readPc()
	instr := instructionTable[opcode]
	fmt.Println(instr.name)
	cpu.opcode = opcode
	cpu.clockCounter = instr.clocks
	instr.am(cpu)
	instr.handler(cpu)
}

func (cpu *Cpu6502) reset() {
	cpu.s = 0xFD // stack pointer decrements by 2
	pcL := cpu.bus.CpuRead(resVector)
	pcH := cpu.bus.CpuRead(resVector + 1)
	cpu.pc = uint16(pcH)<<8 | uint16(pcL)
	cpu.s = 0x34 // set flags U, B, Iconst irqVecL uint16 = 0xFFFE
}

func (cpu *Cpu6502) interrupt(vector uint16) {
	// save cpu backup on the stack
	pcH := uint8(cpu.pc >> 8)
	cpu.push(pcH)
	pcL := uint8(cpu.pc)
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
