package heap

// symbolic reference
type SymRef struct {
	cp        *ConstantPool//本类的本运行时常量池 new的时候放进来的
	//这个全限定名可以唯一确定一个类
	className string//引用的类名  字符串 new的时候放进来的
	class     *Class//引用的类 本来没有 要用了再用className来加载 加载一次就缓存了
}

//返回引用的类
func (self *SymRef) ResolvedClass() *Class {
	if self.class == nil {
		//第一次来
		//拜托类加载器加载这个类(类加载器也有缓存)//访问限制检测 不能访问就报错
		self.resolveClassRef()
	}
	return self.class
}

// jvms8 5.4.3.1
func (self *SymRef) resolveClassRef() {

	thisClass := self.cp.class//本类
	refClass := thisClass.loader.LoadClass(self.className)//引用的类

	//public 或者 同包的class 就能引用  [类访问权限]只有public和无(除非是内部类
	if !refClass.isAccessibleTo(thisClass) {
		//没有权限引用 就异常
		panic("java.lang.IllegalAccessError")
	}

	self.class = refClass
}
