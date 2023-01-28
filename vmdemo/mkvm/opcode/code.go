package opcode

type Instructions []byte
type OpCode byte

const (
	OpConstant OpCode = iota
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*Definition{}
