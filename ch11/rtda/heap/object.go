package heap

//slot的组成单位
type Object struct {
	class  *Class //本类信息
	data interface{}

	extra interface{}//额外信息 如:本对象是class对象的 就记Class
}

//class里面调用
// create normal (non-array) object
func newObject(class *Class) *Object {//new 的时候用
	return &Object{
		class: class,
		data: newSlots(class.instanceSlotCount),
		//instanceSlotCount是加载类的时候就已经算好的 和动态连接一起
		//每个字段都分配好索引了
	}
}



func (self *Object) IsInstanceOf(class *Class) bool {

	return class.isAssignableFrom(self.class)
}
func (self *Class) isJlObject() bool {
	return self.name == "java/lang/Object"
}
func (self *Class) isJlCloneable() bool {
	return self.name == "java/lang/Cloneable"
}
func (self *Class) isJioSerializable() bool {
	return self.name == "java/io/Serializable"
}

// getters
func (self *Object) Class() *Class {

	return self.class
}
func (self *Object) Fields() Slots {
	return self.data.(Slots)//转型
}
func (self *Object) Extra() interface{} {
	return self.extra
}
func (self *Object) SetExtra(extra interface{}) {
	self.extra = extra
}
func (self *Object) Data() interface{} {
	return self.data
}

// 给 对象的 这个名字这个描述符  的字段 赋值为 refValue
func (self *Object) SetRefVar(name, descriptor string, refValue *Object) {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	slots.SetRef(field.slotId, refValue)
}

//根据 [字段名 描述符 是不是类变量] 查找字段，
func (self *Class) getField(name, descriptor string, isStatic bool) *Field {
	for c := self; c != nil; c = c.superClass {
		for _, field := range c.fields {
			if field.IsStatic() == isStatic &&field.name == name && field.descriptor == descriptor {
				return field
			}
		}
	}
	return nil
}

func (self *Object) GetRefVar(name, descriptor string) *Object {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	return slots.GetRef(field.slotId)
}

