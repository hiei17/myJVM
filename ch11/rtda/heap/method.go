package heap

import "jvmgo/ch11/classfile"

type Method struct {
	ClassMember
	/*
	type ClassMember struct {
				  //3个都直接从classfile.MemberInfo里面拿
		accessFlags uint16
		name        string
		descriptor  string

				  //Class结构体指针，这样可以通过字段或方法访问到它所属的类。
				  //构造函数传进来赋值
		class       *Class//所属类
	}
	*/
	maxStack  uint
	maxLocals uint
	code      []byte//字节码

	//参数个数
	argSlotCount uint

	exceptionTable ExceptionTable

	lineNumberTable *classfile.LineNumberTableAttribute
}

/*//根据class文件中的方法信息创建Method表
func newMethods(class *Class, cfMethods []*classfile.MemberInfo) []*Method {

	methods := make([]*Method, len(cfMethods))


	for i, cfMethod := range cfMethods {

		methods[i] = &Method{}
		methods[i].class = class

		//父类方法 拿到: 访问标志 名字 描述符
		methods[i].copyMemberInfo(cfMethod)

		//code属性 里面有: MaxStack MaxLocals 字节码
		methods[i].copyCodeAttributes(cfMethod)
		methods[i].calcArgSlotCount()
	}

	return methods
}*/

//根据class文件中的方法信息 创建class的Method表
func newMethods(class *Class, cfMethods []*classfile.MemberInfo) []*Method {
	methods := make([]*Method, len(cfMethods))
	//classfile中每个方法
	for i, cfMethod := range cfMethods {
		methods[i] = newMethod(class, cfMethod)
	}
	return methods
}

func newMethod(class *Class, cfMethod *classfile.MemberInfo) *Method {
	method := &Method{}
	method.class = class

	//父类方法 拿到: 访问标志 名字 描述符
	method.copyMemberInfo(cfMethod)

	//code属性 里面有: MaxStack MaxLocals 字节码
	method.copyAttributes(cfMethod)

	//方法描述符
	methodDescriptor := parseMethodDescriptor(method.descriptor)

	//根据方法描述符 得到传参个数
	method.calcArgSlotCount(methodDescriptor.parameterTypes)//argSlotCount

	if method.IsNative() {
		//如果是本地方法，则注入字节码和其他信息。
		method.injectCodeAttribute(methodDescriptor.returnType)
	}

	return method
}

//本地方法在class文件中没有Code属性 都要自己写
func (self *Method) injectCodeAttribute(returnType string) {

	//本地方法帧的操作数栈至少要能容纳返回值，为了简化代码，暂时给maxStack字段赋值为4
	self.maxStack = 4

	//本地方法帧的局部变量表只用来存放参数值
	self.maxLocals = self.argSlotCount

	//code字段，也就是本地方法的字节码
	//第一条指令都是0xFE 第二条指令 根据函数的返回值 相应的返回指令。
	switch returnType[0] {
		case 'V': self.code = []byte{0xfe, 0xb1} // return
		case 'D': self.code = []byte{0xfe, 0xaf} // dreturn
		case 'F': self.code = []byte{0xfe, 0xae} // freturn
		case 'J': self.code = []byte{0xfe, 0xad} // lreturn
		case 'L', '[': self.code = []byte{0xfe, 0xb0} // areturn
		default: self.code = []byte{0xfe, 0xac} // ireturn
	}
}



func (self *Method) copyAttributes(cfMethod *classfile.MemberInfo) {
	codeAttr := cfMethod.CodeAttribute();
	if codeAttr == nil {
		return
	}
	self.maxStack = codeAttr.MaxStack()
	self.maxLocals = codeAttr.MaxLocals()
	self.code = codeAttr.Code()
	self.exceptionTable = newExceptionTable(codeAttr.ExceptionTable(),
		self.class.constantPool)

	self.lineNumberTable = codeAttr.LineNumberTableAttribute()
}

func (self *Method) calcArgSlotCount(paramTypes []string) {
	
	for _, paramType := range paramTypes {
		self.argSlotCount++
		if paramType == "J" || paramType == "D" {
			self.argSlotCount++
		}
	}
	if !self.IsStatic() {
		self.argSlotCount++
	}
}
func (self *Method) FindExceptionHandler(exClass *Class, pc int) int {
	handler := self.exceptionTable.findExceptionHandler(exClass, pc)
	if handler != nil {
		return handler.handlerPc
	}
	return -1
}
func (self *Method) GetLineNumber(pc int) int {

	//不一定有哦

	if self.IsNative() {
		return -2
	}

	if self.lineNumberTable == nil {
		return -1
	}
	return self.lineNumberTable.GetLineNumber(pc)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++getters
func (self *Method) MaxStack() uint {
	return self.maxStack
}
func (self *Method) MaxLocals() uint {
	return self.maxLocals
}
func (self *Method) Code() []byte {
	return self.code
}
func (self *Method) ArgSlotCount() uint {
	return self.argSlotCount
}

func (self *Method) IsSynchronized() bool {
	return 0 != self.accessFlags&ACC_SYNCHRONIZED
}
func (self *Method) IsBridge() bool {
	return 0 != self.accessFlags&ACC_BRIDGE
}
func (self *Method) IsVarargs() bool {
	return 0 != self.accessFlags&ACC_VARARGS
}
func (self *Method) IsNative() bool {
	return 0 != self.accessFlags&ACC_NATIVE
}
func (self *Method) IsAbstract() bool {
	return 0 != self.accessFlags&ACC_ABSTRACT
}
func (self *Method) IsStrict() bool {
	return 0 != self.accessFlags&ACC_STRICT
}


