//ZipEntry表示ZIP或JAR文件形式的类路径
package classpath
import "archive/zip"
import "errors"
import "io/ioutil"
import "path/filepath"
type ZipEntry struct {
	absPath string//存放ZIP或JAR文件的绝对路径
}

 //和DirEntry 一样
func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)//转绝对路径
	if err != nil {
		panic(err)
	}
	return &ZipEntry{absPath}
}
//和DirEntry 一样
func (self *ZipEntry) String() string {
	return self.absPath
}

//从ZIP文件中提取class文件
func (self *ZipEntry) readClass(className string) ([]byte, Entry, error) {
	//首先打开ZIP文件，如果这一步出错的话，直接返回
	r, err := zip.OpenReader(self.absPath)
	if err != nil {
		return nil, nil, err
	}

	//使用了defer语句来确保打开的文件得以关闭
	defer r.Close()

	//遍历ZIP压缩包里的文件，看能否找到class文件
	for _, f := range r.File {

		if f.Name != className {
			continue;
		}

		//如果能找到，则打开class文件，把内容读取出来，并返回

		rc, err := f.Open()//打开
		//返回错误信息
		if err != nil {
			return nil, nil, err
		}

		// 确保打开的文件得以关闭
		defer rc.Close()

		data, err := ioutil.ReadAll(rc)//todo 读到data里面
		//返回错误信息
		if err != nil {
			return nil, nil, err
		}
		//数据 本对象 异常
		return data, self, nil

	}
	//如果找不到
	return nil, nil, errors.New("class not found: " + className)
}

