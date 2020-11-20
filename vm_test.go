package gomachine

import (
	"encoding/binary"
	"testing"
	"time"
)

func TestVM_Execute_Blank(t *testing.T) {
	vm := NewVM(0, 0)
	if err := vm.Execute([]byte{}); err != nil {
		t.Fatal(err)
	}
}

func TestVM_CPUTimeExhaustion(t *testing.T) {
	vm := NewVM(0, time.Millisecond)
	err := vm.Execute([]byte{InstructionJmp, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != CPUTimeExhausted {
		t.Fatal("expected cpu time exhausted error, got:", err)
	}
}

func TestVM_Execute_StoreLoadRAM(t *testing.T) {
	vm := NewVM(2, 0)
	if err := vm.Execute([]byte{
		InstructionUint8Load, 0x0A, InstructionUint8Dump,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		InstructionUint8Load, 0x00, InstructionMemoryUint8Load,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}); err != nil {
		t.Fatal(err)
	}
	if vm.Registers[0] != 0x0A {
		t.Fatal("register not 0x0A:", vm.Registers[0])
	}
	if vm.Memory[1] != 0x0A {
		t.Fatal("RAM not 0x0A:", vm.Memory[1])
	}
}

func BenchmarkVM_Execute_Add10000000Numbers(b *testing.B) {
	x := make([]byte, 4)
	binary.LittleEndian.PutUint32(x, 10000000)
	instructions := []byte{
		InstructionUint32Load,
		x[0], x[1], x[2], x[3],
		InstructionMoveR1ToR3,
		InstructionUint8Load,
		0x01,
		InstructionMoveR1ToR2,
		InstructionUint8Load,
		0x00,
		InstructionUnsignedAdd,
		InstructionJmpIfNe,
		0x0B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	vm := NewVM(0, 0)
	b.ResetTimer()
	b.ReportAllocs()
	if err := vm.Execute(instructions); err != nil {
		b.Fatal(err)
	}
	if vm.Registers[0] != 10000000 {
		b.Fatal("not 10000000:", vm.Registers[0])
	}
}

func BenchmarkVM_Execute_Add10000000Numbers_CPUTimeCheck(b *testing.B) {
	x := make([]byte, 4)
	binary.LittleEndian.PutUint32(x, 10000000)
	instructions := []byte{
		InstructionUint32Load,
		x[0], x[1], x[2], x[3],
		InstructionMoveR1ToR3,
		InstructionUint8Load,
		0x01,
		InstructionMoveR1ToR2,
		InstructionUint8Load,
		0x00,
		InstructionUnsignedAdd,
		InstructionJmpIfNe,
		0x0B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	vm := NewVM(0, time.Second)
	b.ResetTimer()
	b.ReportAllocs()
	if err := vm.Execute(instructions); err != nil {
		b.Fatal(err)
	}
	if vm.Registers[0] != 10000000 {
		b.Fatal("not 10000000:", vm.Registers[0])
	}
}
