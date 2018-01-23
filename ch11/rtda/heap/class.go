package heap

import (
	"jvmgo/ch11/classfile"
	"strings"
)
//放在方法区 实例们共有的类信息
// 从class-film获取的类信息
//第一次使用这个class时 加载后把类信息放入方法区

// name, superClassName and interfaceNames are all binary names(jvms8-4.2.1)
type Class struct {//从classFile拿过来 名字一样的都差不多 就加了对本class引用

	accessFlags       uint16

	//name、superClassName和interfaceNames字段分别存放类名、超类名和接口名。
        // 都是完全限定名，具有java/lang/Object的形式
	//3个都是存常量池坐标
	name              string // thisClassName
	superClassName    string
	interfaceNames    []string

	constantPool      *ConstantPool//运行时常量池

	fields            []*Field//类字段
	methods           []*Method//方法表


	loader            *ClassLoader//类加载器

	superClass        *Class//父类
	interfaces        []*Class//接口

	staticSlotCount   uint//类变量 几个

	instanceSlotCount uint//实例变量 几个 加载类文件的时间算出来的 LoadClass

	staticVars        Slots//类变量值

	initStarted       bool

	jClass *Object // java.lang.Class 实例 通过jClass字段，每个Class结构体实例都与一个类对象关联。

	sourceFile string
}

//把ClassFile结构体转换成Class结构体
//TODO 这些数据都由ClassFile里面拿来
func newClass(cf *classfile.ClassFile) *Class {

	class := &Class{}

	//直接get字段拿过来
	class.accessFlags = cf.AccessFlags()//access_flags 是访问标志的含义
	//3个只存个名字 完全限定名 //不存坐标了 存字符串
	class.name = cf.ClassName()
	class.superClassName = cf.SuperClassName()
	class.interfaceNames = cf.InterfaceNames()


	// 不再转来转去存索引了 直接存
	class.constantPool = newConstantPool(class, cf.ConstantPool())


	//除了拷贝过来的 还加了一个字段 class
	//new出来以后每个字段有	访问标志 名字 描述  值在常量池的索引
	class.fields = newFields(class, cf.Fields())
	//new出来以后每个方法有	访问标志 名字 描述  Code属性(MaxStack MaxLocals 字节码
	class.methods = newMethods(class, cf.Methods())

	class.sourceFile = getSourceFile(cf)
	return class
}

func getSourceFile(cf *classfile.ClassFile) string {
	if sfAttr := cf.SourceFileAttribute(); sfAttr != nil {
		return sfAttr.FileName()
	}

	//并不是每个class文件中都有源文件信息，这个因编译时的编译器选项而异
	return "Unknown"
}

func (self *Class) StartInit() {
	self.initStarted = true
}
func (self *Class) InitStarted() bool {
	return self.initStarted
}

func (self *Class) GetClinitMethod() *Method {
	return self.GetStaticMethod("<clinit>", "()V")
}

//8个方法差不多 都是和配置比较 访问标志 某位有没有被设置
func (self *Class) IsPublic() bool {
	//TODO 检查 本标志位 是否被 设置
	return 0 != self.accessFlags&ACC_PUBLIC
}
func (self *Class) IsFinal() bool {
	return 0 != self.accessFlags&ACC_FINAL
}
func (self *Class) IsSuper() bool {
	return 0 != self.accessFlags&ACC_SUPER
}
func (self *Class) IsInterface() bool {
	return 0 != self.accessFlags&ACC_INTERFACE
}
func (self *Class) IsAbstract() bool {
	return 0 != self.accessFlags&ACC_ABSTRACT
}
func (self *Class) IsSynthetic() bool {
	return 0 != self.accessFlags&ACC_SYNTHETIC
}
func (self *Class) IsAnnotation() bool {
	return 0 != self.accessFlags&ACC_ANNOTATION
}
func (self *Class) IsEnum() bool {
	return 0 != self.accessFlags&ACC_ENUM
}


// jvms 5.4.4类访问权限只有public和无(除非是内部类
func (self *Class) isAccessibleTo(other *Class) bool {
	//public 或者同包
	return self.IsPublic() || self.GetPackageName() == other.GetPackageName()
}

//去掉最后一个/后面的
func (self *Class) GetPackageName() string {
	//比如类名是java/lang/Object，则它的包名就是java/lang。
	// 如果类定义在默认包中，它的包名是空字符串。
	if i := strings.LastIndex(self.name, "/"); i >= 0 {
		return self.name[:i]
	}
	return ""
}

//如java/lang/Object转java.lang.Object
//得.分隔的类名
func (self *Class) JavaName() string {
	return strings.Replace(self.name, "/", ".", -1)
}

//new指令调用
func (self *Class) NewObject() *Object {
	return newObject(self)
}



//遍历类里面的所有方法 名字和描述符都对上 就是 了
func (self *Class) GetMainMethod() *Method {
	return self.GetStaticMethod("main", "([Ljava/lang/String;)V")
}

func (self *Class) GetStaticMethod(name, descriptor string) *Method {
	for _, method := range self.methods {
		if method.IsStatic() &&
			method.name == name && method.descriptor == descriptor {
			return method
		}
	}
	return nil
}

//判断类是否是基本类型的类，
func (self *Class) IsPrimitive() bool {
	_, ok := primitiveTypes[self.name]
	return ok
}

// getters
func (self *Class) AccessFlags() uint16 {
	return self.accessFlags
}
func (self *Class) Name() string {
	return self.name
}
func (self *Class) ConstantPool() *ConstantPool {
	return self.constantPool
}
func (self *Class) Fields() []*Field {
	return self.fields
}
func (self *Class) Methods() []*Method {
	return self.methods
}
func (self *Class) Loader() *ClassLoader {
	return self.loader
}
func (self *Class) SuperClass() *Class {
	return self.superClass
}
func (self *Class) Interfaces() []*Class {
	return self.interfaces
}
func (self *Class) StaticVars() Slots {
	return self.staticVars
}
func (self *Class) JClass() *Object {
	return self.jClass
}

func (self *Class) SourceFile() string {
	return self.sourceFile
}
func (self *Class) GetRefVar(fieldName, fieldDescriptor string) *Object {
	field := self.getField(fieldName, fieldDescriptor, true)
	return self.staticVars.GetRef(field.slotId)
}
func (self *Class) GetInstanceMethod(name, descriptor string) *Method {
	return self.getMethod(name, descriptor, false)
}
func (self *Class) getMethod(name, descriptor string, isStatic bool) *Method {
	for c := self; c != nil; c = c.superClass {
		for _, method := range c.methods {
			if method.IsStatic() == isStatic &&
				method.name == name &&
				method.descriptor == descriptor {

				return method
			}
		}
	}
	return nil
}
func (self *Class) SetRefVar(fieldName, fieldDescriptor string, ref *Object) {
	field := self.getField(fieldName, fieldDescriptor, true)
	self.staticVars.SetRef(field.slotId, ref)
}

