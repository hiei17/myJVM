package base
import "jvmgo/ch11/rtda"
type Instruction interface {//所有指令要实现的接口
	//从字节码中提取操作数
	FetchOperands(reader *BytecodeReader)
	//执行指令
	Execute(frame *rtda.Frame)
}

//**********************************************以下是可以复用的指令父类*****************************************************///
//空指令
type NoOperandsInstruction struct {}
func (self *NoOperandsInstruction) FetchOperands(reader *BytecodeReader) {
	// nothing to do
}

//跳转指令
type BranchInstruction struct {
	Offset int
}
func (self *BranchInstruction) FetchOperands(reader *BytecodeReader) {
	//字节码中读取一个uint16整数，转成int后赋给Offset字段
	self.Offset = int(reader.ReadInt16())
}

//根据索引存取局部变量表
type Index8Instruction struct {
	Index uint//局部变量表索引
}
func (self *Index8Instruction) FetchOperands(reader *BytecodeReader) {
	//从字节码中读取一个int8整数，转成uint后赋给Index字段
	self.Index = uint(reader.ReadUint8())
}


//访问运行时常量池
type Index16Instruction struct {
	Index uint//常量池索引由两字节操作数给出
}
func (self *Index16Instruction) FetchOperands(reader *BytecodeReader) {
	self.Index = uint(reader.ReadUint16())
}