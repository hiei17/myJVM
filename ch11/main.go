//main是一个特殊的包，这个包所在的目录（可以叫作任何名字）会被编译为可执行文件
package main



//Go程序的入口也是main（）函数，但是不接收任何参数
func main() {
	//命令行返回的对象cmd
	cmd := parseCmd()
	if cmd.versionFlag {
		//输出（一个滥竽充数的）版本信息
		println("version 0.0.1")
	} else if cmd.helpFlag || cmd.class == "" {
		printUsage()
	} else {
		//启动Java虚拟机
		newJVM(cmd).start()
	}
}



/*func printClassInfo(cf *classfile.ClassFile) {

	fmt.Printf("version: %v.%v\n", cf.MajorVersion(), cf.MinorVersion())
	fmt.Printf("constants count: %v\n", len(cf.ConstantPool()))
	fmt.Printf("access flags: 0x%x\n", cf.AccessFlags())
	fmt.Printf("this class: %v\n", cf.ClassName())
	fmt.Printf("super class: %v\n", cf.SuperClassName())
	fmt.Printf("interfaces: %v\n", cf.InterfaceNames())
	fmt.Printf("fields count: %v\n", len(cf.Fields()))
	for _, f := range cf.Fields() {
		fmt.Printf(" %s\n", f.Name())
	}
	fmt.Printf("methods count: %v\n", len(cf.Methods()))
	for _, m := range cf.Methods() {
		fmt.Printf(" %s\n", m.Name())
	}

}*/





