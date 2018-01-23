package classpath
import "os"
import "path/filepath"
type Classpath struct {
	//用jre目录确定的
	bootClasspath Entry//启动类路径 jre/lib/*
	extClasspath  Entry//扩展类路径 jre/lib/ext/*

	//随便指定 可目录 也压缩包 可* 可多个;隔开
	userClasspath Entry//用户类路径 默认是当前路径
}
func Parse(jreOption, cpOption string) *Classpath {
	oneClassPath := &Classpath{}
	//得到jre目录 给启动类路径 扩展类路径 赋值
	oneClassPath.parseBootAndExtClasspath(jreOption)//启动类路径和扩展类路径
	//用户有传cpOption就用 没有就用当前目录
	oneClassPath.parseUserClasspath(cpOption)
	return oneClassPath
}
func (self *Classpath) parseBootAndExtClasspath(jreOption string) {
	//1优先使用用户输入的-Xjre选项作为jre目录。
	// 2当前目录下寻找jre目录。
	// 3用JAVA_HOME环境变量。
	jreDir := getJreDir(jreOption)

	//启动类路径
	jreLibPath := filepath.Join(jreDir, "lib", "*")// jre/lib/*
	self.bootClasspath = newWildcardEntry(jreLibPath)

	//扩展类路径
	jreExtPath := filepath.Join(jreDir, "lib", "ext", "*")// jre/lib/ext/*
	self.extClasspath = newWildcardEntry(jreExtPath)
}

func (self *Classpath) parseUserClasspath(cpOption string) {
	//如果用户没有提供-classpath/-cp选项，则使用当前目录作为用户类路径
	if cpOption == "" {
		cpOption = "."
	}
	self.userClasspath = newEntry(cpOption)
}


func getJreDir(jreOption string) string {
	//优先使用用户输入的-Xjre选项作为jre目录。
	if jreOption != "" && exists(jreOption) {
		return jreOption
	}

	// 如果没有输入该选项，则在当前目录下寻找jre目录。
	if exists("./jre") {
		return "./jre"
	}

	// 如果找不到，尝试使用JAVA_HOME环境变量
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		return filepath.Join(jh, "jre")
	}
	panic("Can not find jre folder!")
}

//目录是不是存在
func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//依次从启动类路径、扩展类路径和用户类路径中搜索class文件
func (self *Classpath) ReadClass(className string) ([]byte, Entry, error) {
	className = className + ".class"

	//启动
	data, entry, err := self.bootClasspath.readClass(className);
	if  err == nil {
		return data, entry, err
	}

	//扩展
	if data, entry, err := self.extClasspath.readClass(className); err == nil {
		return data, entry, err
	}

	//用户
	return self.userClasspath.readClass(className)
}

func (self *Classpath) String() string {
	return self.userClasspath.String()
}