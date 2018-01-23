package heap

//本类是不是可以赋值给 other类
func (self *Class) isAssignableFrom(other *Class) bool {
	s, t := other, self
	//1是同一个
	if s == t {
		return true
	}

	if s.IsArray() {
		//s数组 t不是数组
		if !t.IsArray() {
			//t是不是Object
			if !t.IsInterface() {
				return t.isJlObject()
			} else {//t是不是Cloneable和Serializable 接口
				return t.isJlCloneable() || t.isJioSerializable()
			}
		} else {//s t 都是数组 看成员是不是可以
			sc := s.ComponentClass()
			tc := t.ComponentClass()
			return sc == tc || tc.isAssignableFrom(sc)
		}
	} else {
		//s是接口
		if s.IsInterface() {
			//t不是接口
			if !t.IsInterface() {//t只能是Object
				return t.isJlObject()
			} else {//t也是接口 只能t是s的父接口
				return t.isSuperInterfaceOf(s)
			}
		} else {
			//二个都不是接口 看看是不是other的子类
			if !t.IsInterface() {
				return s.IsSubClassOf(t)
			} else {//目标t是接口 s不是 看s有没有实现t
				return s.IsImplements(t)
			}
		}
	}

	return  false;
}

//继承链往上找 到底都找不到就算了
func (self *Class) IsSubClassOf(other *Class) bool {
	for c := self.superClass; c != nil; c = c.superClass {
		if c == other {
			return true
		}
	}
	return false
}

func (self *Class) IsImplements(iface *Class) bool {
	//继承链往上找
	for c := self; c != nil; c = c.superClass {
		//每一个接口
		for _, i := range c.interfaces {
			//是这个接口 或其子接口
			if i == iface || i.isSubInterfaceOf(iface) {
				return true
			}
		}
	}
	return false
}

func (self *Class) isSubInterfaceOf(iface *Class) bool {
	//遍历所有接口
	for _, superInterface := range self.interfaces {

		//是这个接口 是子接口
		if superInterface == iface || superInterface.isSubInterfaceOf(iface) {
			return true
		}
	}
	return false
}

// c extends self
func (self *Class) IsSuperClassOf(other *Class) bool {
	return other.IsSubClassOf(self)
}

// iface extends self
func (self *Class) isSuperInterfaceOf(iface *Class) bool {
	return iface.isSubInterfaceOf(self)
}