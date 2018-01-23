package heap

import (
	"jvmgo/ch11/classpath"
	"fmt"
	"jvmgo/ch11/classfile"
)

/*
class names:
    - primitive types: boolean, byte, int ...
    - primitive arrays: [Z, [B, [I ...
    - non-array classes: java/lang/Object ...
    - array classes: [Ljava/lang/Object; ...
*/
type ClassLoader struct {

	classPath   *classpath.Classpath //找文件的地方

					 //mark 方法区
					 // 记录已经加载的类数据，key是类的完 全限定名
	classMap    map[string]*Class    // loaded classes
	//是否把指令执行信息输出到控制台
	verboseFlag bool
}

//目前只main里面调用一次
// 类加载器构造方法
func NewClassLoader(myClassPath *classpath.Classpath, verboseFlag bool) *ClassLoader {
	loader := &ClassLoader{
		classPath: myClassPath,
		verboseFlag: verboseFlag,
		classMap: make(map[string]*Class),
	}

	//处理Class类
	loader.loadBasicClasses()
	//和数组类一样，基本类型的类也是由Java虚拟机在运行期间生成的
	loader.loadPrimitiveClasses()
	return loader
}

func (self *ClassLoader) loadPrimitiveClasses() {

	for primitiveType, _ := range primitiveTypes {//遍历每种基础类型

		self.loadPrimitiveClass(primitiveType)
	}
	//注意 "void"也有一个
}

//生成这个基础类型的class扔入方法区
func (self *ClassLoader) loadPrimitiveClass(className string) {
	class := &Class{
		accessFlags: ACC_PUBLIC, // todo
		name:        className,
		loader:      self,
		initStarted: true,

		//基本类型的类没有超类，没有实现接口
	}
	newClassObject:=self.classMap["java/lang/Class"].NewObject()
	newClassObject.extra=class

	class.jClass = newClassObject

	self.classMap[className] = class
}

func (self *ClassLoader) loadBasicClasses() {

	//会触发 java.lang.Object等类和接口的加载
	jlClassClass := self.LoadClass("java/lang/Class")

	//遍历方法区内所有类 处理class类
	for _, class := range self.classMap {
		if class.jClass != nil {
			continue;
		}
		//new 一个空白的class实例
		newClassObject := jlClassClass.NewObject()
		//类里面有指向它的class对象的
		class.jClass = newClassObject
		//class对象里面有指向它代表的类的
		newClassObject.extra = class
	}
}

func (self *ClassLoader) LoadClass(name string) *Class {

	//方法区已有的
	if class, ok := self.classMap[name]; ok {
		return class // already loaded
	}

	//第一次出现 要加载
	var class *Class
	if name[0] == '[' {
		// array class 数组
		class = self.loadArrayClass(name)
	} else {
		//普通类
		class = self.loadNonArrayClass(name)
	}

	//有class类
	if jlClassClass, ok := self.classMap["java/lang/Class"]; ok {
		//new 一个空白的class实例
		aClassObject := jlClassClass.NewObject()
		//类里面有指向它的class对象的
		class.jClass = aClassObject
		//class对象里面有指向它代表的类的
		aClassObject.extra = class

	}
	return class
}
/*
主要的变动是粗体部分。在类加载完之后，看java.lang.Class是
否已经加载。如果是，则给类关联类对象。这样，在
loadBasicClasses（）和LoadClass（）方法的配合之下，所有加载到方
法区的类都设置好了jClass字段。

*/


func (self *ClassLoader) loadArrayClass(name string) *Class {
	//运行时自主生成 不通过classfile
	class := &Class{
		accessFlags: ACC_PUBLIC, // todo
		name: name,
		loader: self,
		initStarted: true,//数组类不要初始化
		superClass: self.LoadClass("java/lang/Object"),//超类是java.lang.Object
		//都实现这2个接口
		interfaces: []*Class{
			self.LoadClass("java/lang/Cloneable"),
			self.LoadClass("java/io/Serializable"),
		},
	}
	self.classMap[name] = class
	return class
}


func (self *ClassLoader) loadNonArrayClass(name string) *Class {

	//找到class文件 读进来原始[]byte
	data, entry := self.readClass(name)

	//mark 解析成class_file 从而得到class  里面同时放了classMap
	class := self.defineClass(data)

	link(class)
	if self.verboseFlag {
		fmt.Printf("[Loaded %s from %s]\n", name, entry)
	}
	return class
}


//只是调用了Classpath的ReadClass（）方法
func (self *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {

	//classpath 遍历找到[]byte
	data, entry, err := self.classPath.ReadClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException: " + name)
	}

	//返回值  []byte 类路径
	return data, entry
}

// jvms 5.3.5
//得本class 父类和接口都加载了 链上了
func (self *ClassLoader) defineClass(data []byte) *Class {

	//[]byte→classfile→class
	class := parseClass(data)

	class.loader = self
	//class填入父类
	resolveSuperClass(class)

	//class填入所有接口
	resolveInterfaces(class)

	//缓存class
	self.classMap[class.name] = class
	return class
}


////[]byte→classfile→class
func parseClass(data []byte) *Class {

	cf, err := classfile.Parse(data)

	if err != nil {
		//panic("java.lang.ClassFormatError")
		panic(err)
	}
	return newClass(cf)
}

// jvms 5.4.3.1
func resolveSuperClass(class *Class) {

	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}
func resolveInterfaces(class *Class) {

	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}

func link(class *Class) {
	//虚拟机规范要求了 暂时没做
	//verify(class) // todo

	//mark 给类变量分配空间并给予初始值
	prepare(class)
}

// jvms 5.4.2//分配slot 算共几个  类变量赋初值
func prepare(class *Class) {

	//实例变量 发编号 算共多少个
	calcInstanceFieldSlotIds(class)

	//类变量static发编号 算共多少个
	calcStaticFieldSlotIds(class)

	//static final赋初值
	allocAndInitStaticVars(class)
}

//给实例发编号 计算有多少个
func calcInstanceFieldSlotIds(class *Class) {
	slotCount := uint(0)
	if class.superClass != nil {
		//父类继承的实例变量
		slotCount = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() {
			field.slotId = slotCount
			slotCount++
			if field.isLongOrDouble() {
				slotCount++
			}
		}
	}
	class.instanceSlotCount = slotCount
}

func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}


/*因为Go语言会保证新创建的Slot结构体有默认值（num字段是
0，ref字段是nil），而浮点数0编码之后和整数0相同，所以不用做任
何操作就可以保证静态变量有默认初始值（数字类型是0，引用类型
是null）。
如果静态变量属于基本类型或String类型，有final修饰符，
且它的值在编译期已知，则该值存储在class文件常量池中。*/
func allocAndInitStaticVars(class *Class) {

	class.staticVars = newSlots(class.staticSlotCount)

	//类的所有类变量
	for _, field := range class.fields {

		//static final 值在编译时已经确定了 在常量池
		if field.IsStatic() && field.IsFinal() {
			//从常量池里面拿到字段的具体值 然后拿到上一步给字段分配的slot编号
			// 放到class.staticVars相应编号的slot里面
			initStaticFinalVar(class, field)
		}

	}
}

func initStaticFinalVar(class *Class, field *Field) {
	//放类变量的slots
	vars := class.staticVars
	//常量池
	cp := class.constantPool
	//值的常量池编号
	cpIndex := field.ConstValueIndex()

	//准备拿来放的slot编号
	slotId := field.SlotId()

	//不同类型 拿出来转类型 放到指定编号的slot里面
	if cpIndex <=0 {
		return
	}

	//运行时常量池这个编号的常量拿出来 给相应staticVars[slotId]
	//各类型不同的set get
	switch field.Descriptor() {
		case "Z", "B", "C", "S", "I":
			val := cp.GetConstant(cpIndex).(int32)
			vars.SetInt(slotId, val)
		case "J":
			val := cp.GetConstant(cpIndex).(int64)
			vars.SetLong(slotId, val)
		case "F":
			val := cp.GetConstant(cpIndex).(float32)
			vars.SetFloat(slotId, val)
		case "D":
			val := cp.GetConstant(cpIndex).(float64)
			vars.SetDouble(slotId, val)
		case "Ljava/lang/String;":
			goStr := cp.GetConstant(cpIndex).(string)
			jStr := JString(class.Loader(), goStr)
			vars.SetRef(slotId, jStr)
	}

}

