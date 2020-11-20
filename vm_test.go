package gomachine

import (
	"encoding/binary"
	"testing"
)



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
