package classpath
import "io/ioutil"
import "path/filepath"

//表示目录形式的类路径
type DirEntry struct {
	absDir string//目录的绝对路径  构造方法内传入
}
//Go结构体不需要显示实现接口，只要方法匹配即可


func newDirEntry(path string) *DirEntry {
	absDir, err := filepath.Abs(path)//传入的路径 转为绝对路径

	//如果转换过程出现错误
	if err != nil {
		panic(err)//调用panic（）函数终止程序执行
	}
	return &DirEntry{absDir}//创建DirEntry实例并返回
}

//接口要求的方法
//有用到self 是个方法
func (self *DirEntry) readClass(className string) ([]byte, Entry, error) {

	//把目录和class文件名拼成一个完整的路径
	fileName := filepath.Join(self.absDir, className)

	data, err := ioutil.ReadFile(fileName)//调用ioutil包提供的ReadFile（）函数读取class文件内容

	//返回值是读取到的
	// 1.字节数据、
	// 2.最终定位到class文件的Entry
	// 3.错误信息
	return data, self, err
}

func (self *DirEntry) String() string {
	return self.absDir//直接返回目录
}