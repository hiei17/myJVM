package classpath
import "os"
import "strings"

//路径分隔符
const pathListSeparator = string(os.PathListSeparator)

//Entry接口
//代表类路径
// 有4个实现，分别是DirEntry、ZipEntry、CompositeEntry和WildcardEntry。
type Entry interface {
	//参数是class文件的相对路径，路径之间用斜线（/）分隔，文件名有.class后缀
	//比如要读取java.lang.Object类，传入的参数应该是java/lang/Object.class。
	// 返回值是读取到的字节数据、最终定位到class文件的Entry，以及错误信息//Go的函数或方法允许返回多个值，按照惯例，可以使用最后一个返回值作为错误信息
	readClass(className string) ([]byte, Entry, error)//寻找和加载class文件
	String() string//相当于Java中的toString（），用于返回变量的字符串表示
}


//newEntry（）函数根据参数创建不同类型的Entry实例  go没有指定构造函数 随便我们自己指定 一般用这样命名
//调用子类构造方法
func newEntry(path string) Entry {

	//;隔开的多个
	if strings.Contains(path, pathListSeparator) {
		return newCompositeEntry(path)
	}

	//*结尾 匹配目录下全部jar包
	if strings.HasSuffix(path, "*") {
		return newWildcardEntry(path)
	}

	//压缩包
	if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".ZIP") {
		return newZipEntry(path)
	}

	//目录
	return newDirEntry(path)
}

