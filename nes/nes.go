package nes

import (
	"goNES/bus"
	"goNES/cpu"
	"goNES/cpu_interrupts"
	"goNES/cpubus"
	"goNES/dma"
	"goNES/ppu"
)

type Rom struct {
	ProgramRom         []byte
	CharacterRom       []byte
	isHorizontalMirror bool
	mapper             uint8
}

type Nes struct {
	Rom          Rom
	Ram          bus.Ram
	characterMem bus.Ram
	ProgramRom   bus.Rom
	PpuBus       bus.PpuBus
	Interrupts   cpu_interrupts.Interrupts
	Dma          dma.Dma
	Ppu          ppu.Ppu
	CpuBus       cpubus.CpuBus
	Cpu          cpu.Cpu
}

func New(data []byte) Nes {
	return Nes{Rom: parse(data)}
}

func (nes *Nes) Load() {
	nes.Ram = bus.NewRam(2048)
	nes.characterMem = bus.NewRam(0x4000)
	for i := 0; i < len(nes.Rom.CharacterRom); i++ {
		nes.characterMem.Write(uint16(i), nes.Rom.CharacterRom[i])
	}

	nes.ProgramRom = bus.NewRom(nes.Rom.ProgramRom)
	nes.PpuBus = bus.NewPpuBus(&nes.characterMem)
	nes.Interrupts = cpu_interrupts.NewInterrupts()

	nes.Ppu = ppu.NewPpu(nes.PpuBus, nes.Interrupts, nes.Rom.isHorizontalMirror)

	nes.Dma = dma.NewDma(nes.Ram, nes.Ppu)

	nes.CpuBus = cpubus.NewCpuBus(
		nes.Ram,
		nes.ProgramRom,
		nes.Ppu,
		nes.Dma,
	)

	nes.Cpu = cpu.NewCpu(nes.CpuBus, nes.Interrupts)
	nes.Cpu.Reset()
}

func (nes *Nes) frame() {
	for {
		cycle := 0
		if nes.Dma.IsDmaProcessing() {
			nes.Dma.RunDma()
			cycle = 514
		}
		cycle += nes.Cpu.Run()

		/* todo
		    $renderingData = $this->ppu->run($cycle * 3);
            if ($renderingData) {
                $this->cpu->bus->keypad->fetch();
                $this->renderer->render($renderingData);
                break;
            }
		 */
	}
}

func (nes Nes) Start() {
	for {
		nes.frame()
	}
}