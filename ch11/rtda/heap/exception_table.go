package heap

import "jvmgo/ch11/classfile"
type ExceptionTable []*ExceptionHandler

//源头是attr_code 里面 从.class文件解析
type ExceptionHandler struct {
	startPc int
	endPc int
	handlerPc int
	catchType *ClassRef
}

//class文件中的异常处理表转换成 ExceptionTable类型
func newExceptionTable(entries []*classfile.ExceptionTableEntry, cp *ConstantPool) ExceptionTable {

	table := make([]*ExceptionHandler, len(entries))

	for i, entry := range entries {

		table[i] = &ExceptionHandler{

			startPc: int(entry.StartPc()),
			endPc: int(entry.EndPc()),
			handlerPc: int(entry.HandlerPc()),
			catchType: getCatchType(uint(entry.CatchType()), cp),
		}
	}
	return table
}

//从运行时常量池中查找类符号引用
func getCatchType(index uint, cp *ConstantPool) *ClassRef {
	if index == 0 {//是表示catch-all
		return nil
	}
	return cp.GetConstant(index).(*ClassRef)
}


func (self ExceptionTable) findExceptionHandler(exClass *Class, pc int) *ExceptionHandler {

	for _, handler := range self {

		if pc >= handler.startPc && pc < handler.endPc {

			//在class文件中是0 表示可以处理所有异常，这是用来实现finally子句的。
			if handler.catchType == nil {
				return handler // catch-all
			}
			catchClass := handler.catchType.ResolvedClass()
			if catchClass == exClass || catchClass.IsSuperClassOf(exClass) {
				return handler
			}
		}
	}
	return nil
}
