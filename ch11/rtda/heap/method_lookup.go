package heap

//继承中找指定方法
func LookupMethodInClass(class *Class, name, descriptor string) *Method {

	//找到符合的就返回
	for fatherClass := class; fatherClass != nil; fatherClass = fatherClass.superClass {

		for _, method := range fatherClass.methods {
			if method.name == name && method.descriptor == descriptor {
				return method
			}
		}
	}

	return nil
}

//接口中找
func lookupMethodInInterfaces(ifaces []*Class, name, descriptor string) *Method {

	//遍历所有实现的接口
	for _, iface := range ifaces {
		//此接口所有方法
		for _, method := range iface.methods {
			if method.name == name && method.descriptor == descriptor {
				return method
			}
		}
		//此接口的接口
		method := lookupMethodInInterfaces(iface.interfaces, name, descriptor)

		if method != nil {
			return method
		}
	}

	return nil
}
