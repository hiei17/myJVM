package extended

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/instructions/loads"
	"jvmgo/ch11/instructions/stores"
	"jvmgo/ch11/instructions/math"
	"jvmgo/ch11/rtda"
)

/*
加载类指令、存储类指令、ret指令和iinc指令
需要按索引访问局部变量表，索引以uint8的形式存在字节码中。
对于大部分方法来说，局部变量表大小都不会超过256，所以用一字节来表示索引就
够了。但是 超过这限制,wide指令来扩展前述指令。
*/

// Extend local variable index by additional bytes
type WIDE struct {
	//存放被 改变的指令
	modifiedInstruction base.Instruction
}

//不同的地方是 变成ReadUint16 读2字节
func (self *WIDE) FetchOperands(reader *base.BytecodeReader) {
	opcode := reader.ReadUint8()
	switch opcode {
	case 0x15:
		inst := &loads.ILOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x16:
		inst := &loads.LLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x17:
		inst := &loads.FLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x18:
		inst := &loads.DLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x19:
		inst := &loads.ALOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x36:
		inst := &stores.ISTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x37:
		inst := &stores.LSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x38:
		inst := &stores.FSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x39:
		inst := &stores.DSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x3a:
		inst := &stores.ASTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x84:
		inst := &math.IINC{}
		inst.Index = uint(reader.ReadUint16())
		inst.Const = int32(reader.ReadInt16())
		self.modifiedInstruction = inst
	case 0xa9: // ret 暂时没实现
		panic("Unsupported opcode: 0xa9!")
	}
}

/*wide指令只是增加了索引宽度，并不改变子指令操作，所以其
Execute（）方法只要调用子指令的Execute（）方法即可，*/
//执行原来指令
func (self *WIDE) Execute(frame *rtda.Frame) {
	self.modifiedInstruction.Execute(frame)
}
