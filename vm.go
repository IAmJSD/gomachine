package gomachine

import (
	"errors"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

// Defines the CPU instructions.
const (
	// InstructionUint8Load is used to load a uint8 argument into R1.
	InstructionUint8Load = uint8(iota + 1)

	// InstructionUint16Load is used to load a uint16 argument into R1.
	InstructionUint16Load

	// InstructionUint32Load is used to load a uint32 argument into R1.
	InstructionUint32Load

	// InstructionUint64Load is used to load a uint64 argument into R1.
	InstructionUint64Load

	// InstructionUint8Load is used to load a uint8 argument into R1 from the memory location specified.
	InstructionMemoryUint8Load

	// InstructionUint16Load is used to load a uint16 argument into R1 from the memory location specified.
	InstructionMemoryUint16Load

	// InstructionUint32Load is used to load a uint32 argument into R1 from the memory location specified.
	InstructionMemoryUint32Load

	// InstructionUint64Load is used to load a uint64 argument into R1 from the memory location specified.
	InstructionMemoryUint64Load

	// InstructionMoveR1ToR2 is used to move R1 into R2.
	InstructionMoveR1ToR2

	// InstructionMoveR1ToR3 is used to move R1 into R3.
	InstructionMoveR1ToR3

	// InstructionMoveR2ToR1 is used to move R2 into R1.
	InstructionMoveR2ToR1

	// InstructionMoveR2ToR3 is used to move R2 into R3.
	InstructionMoveR2ToR3

	// InstructionMoveR3ToR1 is used to move R3 into R1.
	InstructionMoveR3ToR1

	// InstructionMoveR3ToR2 is used to move R3 into R2.
	InstructionMoveR3ToR2

	// InstructionMoveR4ToR1 is used to move R4 into R1.
	InstructionMoveR4ToR1

	// InstructionMoveR4ToR2 is used to move R4 into R2.
	InstructionMoveR4ToR2

	// InstructionMoveR4ToR3 is used to move R4 into R3.
	InstructionMoveR4ToR3

	// InstructionUint8Dump is used to dump a uint8 argument from R1 into the memory location specified.
	InstructionUint8Dump

	// InstructionUint16Dump is used to dump a uint16 argument from R1 into the memory location specified.
	InstructionUint16Dump

	// InstructionUint32Dump is used to dump a uint32 argument from R1 into the memory location specified.
	InstructionUint32Dump

	// InstructionUint64Dump is used to dump a uint64 argument from R1 into the memory location specified.
	InstructionUint64Dump

	// InstructionUnsignedAdd is used to add R2 to R1 and treat them as unsigned integers. The result is stored in R1.
	InstructionUnsignedAdd

	// InstructionSignedAdd is used to add R2 to R1 and treat them as signed integers. The result is stored in R1.
	InstructionSignedAdd

	// InstructionUnsignedSub is used to subtract R2 from R1 and treat them as unsigned integers. The result is stored in R1.
	InstructionUnsignedSub

	// InstructionSignedSub is used to subtract R2 from R1 and treat them as signed integers. The result is stored in R1.
	InstructionSignedSub

	// InstructionUnsignedDiv is used to divide R1 against R2 and treat them as unsigned integers. The result is stored in R1, and if you try and divide by 0 it returns 1 in R4.
	InstructionUnsignedDiv

	// InstructionSignedDiv is used to divide R1 against R2 and treat them as signed integers. The result is stored in R1, and if you try and divide by 0 it returns 1 in R4.
	InstructionSignedDiv

	// InstructionUnsignedMod is used to mod R1 against R2 and treat them as unsigned integers. The result is stored in R1, and if you try and divide by 0 it returns 1 in R4.
	InstructionUnsignedMod

	// InstructionSignedMod is used to mod R1 against R2 and treat them as signed integers. The result is stored in R1, and if you try and divide by 0 it returns 1 in R4.
	InstructionSignedMod

	// InstructionBitwiseAnd is used to perform bitwise and on R1 with R2. The result is stored in R1.
	InstructionBitwiseAnd

	// InstructionBitwiseOr is used to perform bitwise or on R1 with R2. The result is stored in R1.
	InstructionBitwiseOr

	// InstructionBitwiseXor is used to perform bitwise xor on R1 with R2. The result is stored in R1.
	InstructionBitwiseXor

	// InstructionBitwiseLeftShift is used to shift R1 left the number of bits specified in R2. The result is stored in R1.
	InstructionBitwiseLeftShift

	// InstructionBitwiseRightShift is used to shift R1 right the number of bits specified in R2. The result is stored in R1.
	InstructionBitwiseRightShift

	// InstructionJmp is used to jump to another place in the bytecode.
	InstructionJmp

	// InstructionJmpIfEq is used to jump if R3 is equal to R1.
	InstructionJmpIfEq

	// InstructionJmpIfNe is used to jump if R3 is not equal to R1.
	InstructionJmpIfNe

	// InstructionJmpIfGt is used to jump if R1 is greater than R3.
	InstructionJmpIfGt

	// InstructionJmpIfLt is used to jump if R1 is less than R3.
	InstructionJmpIfLt

	// InstructionJmpIfGtOrEqual is used to jump if R1 is greater than or equal to R3.
	InstructionJmpIfGtOrEqual
	
	// InstructionJmpIfLtOrEqual is used to jump if R1 is less than or equal to R3.
	InstructionJmpIfLtOrEqual

	// InstructionSyscall is used to make a system call with the instruction in R1. System calls are expected to throw errors in R3.
	InstructionSyscall
)

// InvalidInstructionArgument is used when the instruction expects a argument but none is provided.
var InvalidInstructionArgument = errors.New("no argument provided as the instruction expects one")

// InvalidMemoryLocation is an error which is thrown when the memory location is outside of the memory length.
var InvalidMemoryLocation = errors.New("memory location is outside of the maximum array size")

// InvalidSyscall is an error which is thrown when the bytes reference a syscall which doesn't exist.
var InvalidSyscall = errors.New("syscall is invalid")

// CPUTimeExhausted is returned when the amount of CPU time a user has was exhausted.
var CPUTimeExhausted = errors.New("cpu time is exhausted")

// UnknownInstruction is used when the CPU instruction is unknown.
var UnknownInstruction = errors.New("unknown cpu instruction")

// VM is used to represent the virtual machine.
type VM struct {
	// Memory is used to represent the memory of the virtual machine.
	Memory []byte

	// MaxCPUTime is used to say how much CPU time a VM can use. 0 means unlimited.
	MaxCPUTime time.Duration

	// Syscalls is used to define system calls the virtual machine can do.
	// An error being returned here will error the execution of the VM.
	Syscalls map[uint64]func(*VM) error

	// Defines the CPU registers.
	Registers [4]uint64
}

// Execute is used to execute bytecode on the virtual machine.
func (v *VM) Execute(Bytecode []byte) error {
	// Get the bytecode location and length.
	bytecodeLen := uint64(len(Bytecode))
	if bytecodeLen == 0 {
		// Return no errors. No bytecode was executed.
		return nil
	}
	bytecodePtr := (unsafe.Pointer)(&Bytecode[0])

	// Get the virtual memory location and length.
	virtualMemoryLen := uint64(len(v.Memory))
	var virtualMemory uintptr
	if virtualMemoryLen != 0 {
		virtualMemory = (uintptr)((unsafe.Pointer)(&v.Memory[0]))
	}

	// A pointer to the registers array.
	r1 := &v.Registers[0]
	r2 := &v.Registers[1]
	r3 := &v.Registers[2]
	r4 := &v.Registers[3]

	// Defines if we should stop.
	shouldStop := uintptr(0)

	// Defines if we should do time checks and handle them if so.
	doTimeChecks := v.MaxCPUTime != 0
	var timer *time.Timer
	if doTimeChecks {
		timer = time.AfterFunc(v.MaxCPUTime, func() {
			atomic.StoreUintptr(&shouldStop, 1)
		})
	}
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	// Go through the bytecode.
	bytecodeIndex := uint64(0)
	for bytecodeIndex != bytecodeLen {
	s:
		// Do a time check.
		if doTimeChecks {
			if atomic.LoadUintptr(&shouldStop) == 1 {
				return CPUTimeExhausted
			}
		}

		// Run a switch on this byte to get the instruction.
		switch *(*uint8)(bytecodePtr) {
		// Load from bytecode instructions.
		case InstructionUint8Load:
			bytecodeIndex++
			if bytecodeIndex == bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			*r1 = uint64(*(*uint8)(bytecodePtr))
			*r4 = 0
		case InstructionUint16Load:
			bytecodeIndex += 2
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			*r1 = uint64(*(*uint16)(bytecodePtr))
			*r4 = 0
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
		case InstructionUint32Load:
			bytecodeIndex += 4
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			*r1 = uint64(*(*uint32)(bytecodePtr))
			*r4 = 0
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 3)
		case InstructionUint64Load:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			*r1 = *(*uint64)(bytecodePtr)
			*r4 = 0
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)

		// Load from virtual memory instructions.
		case InstructionMemoryUint8Load:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			memoryLocation := *(*uint64)(bytecodePtr)
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)
			if memoryLocation >= virtualMemoryLen {
				return InvalidMemoryLocation
			}
			*r1 = uint64(*(*uint8)(unsafe.Pointer(virtualMemory + uintptr(memoryLocation))))
			*r4 = 0
		case InstructionMemoryUint16Load:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			memoryLocation := *(*uint64)(bytecodePtr)
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)
			if memoryLocation+1 >= virtualMemoryLen {
				return InvalidMemoryLocation
			}
			*r1 = uint64(*(*uint16)(unsafe.Pointer(virtualMemory + uintptr(memoryLocation))))
			*r4 = 0
		case InstructionMemoryUint32Load:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			memoryLocation := *(*uint64)(bytecodePtr)
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)
			if memoryLocation+3 >= virtualMemoryLen {
				return InvalidMemoryLocation
			}
			*r1 = uint64(*(*uint32)(unsafe.Pointer(virtualMemory + uintptr(memoryLocation))))
			*r4 = 0
		case InstructionMemoryUint64Load:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			memoryLocation := *(*uint64)(bytecodePtr)
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)
			if memoryLocation+7 >= virtualMemoryLen {
				return InvalidMemoryLocation
			}
			*r1 = *(*uint64)(unsafe.Pointer(virtualMemory + uintptr(memoryLocation)))
			*r4 = 0

		// Register move instructions.
		case InstructionMoveR1ToR2:
			*r2 = *r1
			*r4 = 0
		case InstructionMoveR1ToR3:
			*r3 = *r1
			*r4 = 0
		case InstructionMoveR2ToR1:
			*r1 = *r2
			*r4 = 0
		case InstructionMoveR2ToR3:
			*r3 = *r2
			*r4 = 0
		case InstructionMoveR3ToR1:
			*r1 = *r3
			*r4 = 0
		case InstructionMoveR3ToR2:
			*r2 = *r3
			*r4 = 0
		case InstructionMoveR4ToR1:
			*r1 = *r4
			*r4 = 0
		case InstructionMoveR4ToR2:
			*r2 = *r4
			*r4 = 0
		case InstructionMoveR4ToR3:
			*r3 = *r4
			*r4 = 0

		// Memory dump instructions.
		case InstructionUint8Dump:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			memoryLocation := *(*uint64)(bytecodePtr)
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)
			if memoryLocation >= virtualMemoryLen {
				return InvalidMemoryLocation
			}
			*(*uint8)(unsafe.Pointer(virtualMemory + uintptr(memoryLocation))) = uint8(*r1)
			*r4 = 0
		case InstructionUint16Dump:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			memoryLocation := *(*uint64)(bytecodePtr)
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)
			if memoryLocation+1 >= virtualMemoryLen {
				return InvalidMemoryLocation
			}
			*(*uint16)(unsafe.Pointer(virtualMemory + uintptr(memoryLocation))) = uint16(*r1)
			*r4 = 0
		case InstructionUint32Dump:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			memoryLocation := *(*uint64)(bytecodePtr)
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)
			if memoryLocation+3 >= virtualMemoryLen {
				return InvalidMemoryLocation
			}
			*(*uint32)(unsafe.Pointer(virtualMemory + uintptr(memoryLocation))) = uint32(*r1)
			*r4 = 0
		case InstructionUint64Dump:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			memoryLocation := *(*uint64)(bytecodePtr)
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)
			if memoryLocation+7 >= virtualMemoryLen {
				return InvalidMemoryLocation
			}
			*(*uint64)(unsafe.Pointer(virtualMemory + uintptr(memoryLocation))) = *r1
			*r4 = 0

		// Addition instructions.
		case InstructionUnsignedAdd:
			*r1 += *r2
			*r4 = 0
		case InstructionSignedAdd:
			*(*int64)(unsafe.Pointer(r1)) += *(*int64)(unsafe.Pointer(r2))
			*r4 = 0

		// Subtraction instructions.
		case InstructionUnsignedSub:
			*r1 -= *r2
			*r4 = 0
		case InstructionSignedSub:
			*(*int64)(unsafe.Pointer(r1)) -= *(*int64)(unsafe.Pointer(r2))
			*r4 = 0

		// Division instructions.
		case InstructionUnsignedDiv:
			if *r2 == 0 {
				*r4 = 1
			} else {
				*r1 /= *r2
				*r4 = 0
			}
		case InstructionSignedDiv:
			if *r2 == 0 {
				*r4 = 1
			} else {
				*(*int64)(unsafe.Pointer(r1)) /= *(*int64)(unsafe.Pointer(r2))
				*r4 = 0
			}

		// Modulo instructions.
		case InstructionUnsignedMod:
			if *r2 == 0 {
				*r4 = 1
			} else {
				*r1 %= *r2
				*r4 = 0
			}
		case InstructionSignedMod:
			if *r2 == 0 {
				*r4 = 1
			} else {
				*(*int64)(unsafe.Pointer(r1)) %= *(*int64)(unsafe.Pointer(r2))
				*r4 = 0
			}

		// Bitwise instructions.
		case InstructionBitwiseAnd:
			*r1 &= *r2
			*r4 = 0
		case InstructionBitwiseOr:
			*r1 |= *r2
			*r4 = 0
		case InstructionBitwiseXor:
			*r1 ^= *r2
			*r4 = 0
		case InstructionBitwiseLeftShift:
			*r1 <<= *r2
			*r4 = 0
		case InstructionBitwiseRightShift:
			*r1 >>= *r2
			*r4 = 0

		// Jump instruction.
		case InstructionJmp:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			location := *(*uint64)(bytecodePtr)
			if location >= bytecodeLen {
				return InvalidMemoryLocation
			}
			bytecodePtr = (unsafe.Pointer)(&Bytecode[location])
			bytecodeIndex = location
			goto s
		case InstructionJmpIfEq:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			if *r1 == *r3 {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
				location := *(*uint64)(bytecodePtr)
				if location >= bytecodeLen {
					return InvalidMemoryLocation
				}
				bytecodePtr = (unsafe.Pointer)(&Bytecode[location])
				bytecodeIndex = location
				goto s
			} else {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 8)
			}
		case InstructionJmpIfNe:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			if *r1 != *r3 {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
				location := *(*uint64)(bytecodePtr)
				if location >= bytecodeLen {
					return InvalidMemoryLocation
				}
				bytecodePtr = (unsafe.Pointer)(&Bytecode[location])
				bytecodeIndex = location
				goto s
			} else {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 8)
			}
		case InstructionJmpIfGt:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			if *r1 > *r3 {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
				location := *(*uint64)(bytecodePtr)
				if location >= bytecodeLen {
					return InvalidMemoryLocation
				}
				bytecodePtr = (unsafe.Pointer)(&Bytecode[location])
				bytecodeIndex = location
				goto s
			} else {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 8)
			}
		case InstructionJmpIfLt:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			if *r1 < *r3 {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
				location := *(*uint64)(bytecodePtr)
				if location >= bytecodeLen {
					return InvalidMemoryLocation
				}
				bytecodePtr = (unsafe.Pointer)(&Bytecode[location])
				bytecodeIndex = location
				goto s
			} else {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 8)
			}
		case InstructionJmpIfGtOrEqual:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			if *r1 >= *r3 {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
				location := *(*uint64)(bytecodePtr)
				if location >= bytecodeLen {
					return InvalidMemoryLocation
				}
				bytecodePtr = (unsafe.Pointer)(&Bytecode[location])
				bytecodeIndex = location
				goto s
			} else {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 8)
			}
		case InstructionJmpIfLtOrEqual:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			if *r1 <= *r3 {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
				location := *(*uint64)(bytecodePtr)
				if location >= bytecodeLen {
					return InvalidMemoryLocation
				}
				bytecodePtr = (unsafe.Pointer)(&Bytecode[location])
				bytecodeIndex = location
				goto s
			} else {
				bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 8)
			}

		// System call instruction.
		case InstructionSyscall:
			bytecodeIndex += 8
			if bytecodeIndex >= bytecodeLen {
				return InvalidInstructionArgument
			}
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
			syscall := *(*uint64)(bytecodePtr)
			bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 7)
			call, ok := v.Syscalls[syscall]
			*r4 = 0
			if ok {
				// Attempt the system call.
				if err := call(v); err != nil {
					return err
				}
			} else {
				// Invalid system call.
				return InvalidSyscall
			}

		// Handle unknown instruction.
		default:
			return UnknownInstruction
		}

		// Add 1 to the pointer and bytecode index.
		bytecodeIndex++
		bytecodePtr = (unsafe.Pointer)((uintptr)(bytecodePtr) + 1)
	}

	// Keep the bytecode alive.
	runtime.KeepAlive(Bytecode)

	// Return no errors.
	return nil
}

// ClearRegisters is used to clear the registers of the virtual CPU.
func (v *VM) ClearRegisters() {
	_ = v.Execute([]byte{
		InstructionUint8Load,  // Insert uint8 into R1.
		0x00,                  // uint8 for 0.
		InstructionMoveR1ToR2, // Move R2 <- [R1].
		InstructionMoveR1ToR3, // Move R3 <- [R1].
	})
}

// ClearMemory is used to clear the memory of a virtual machine.
func (v *VM) ClearMemory() {
	for i := range v.Memory {
		v.Memory[i] = 0
	}
}

// NewVM is used to create a new virtual machine.
func NewVM(MemoryLength uint64, MaxCPUTime time.Duration) *VM {
	return &VM{
		Memory:     make([]byte, MemoryLength),
		MaxCPUTime: MaxCPUTime,
		Syscalls:   map[uint64]func(*VM) error{},
		Registers:  [4]uint64{},
	}
}
