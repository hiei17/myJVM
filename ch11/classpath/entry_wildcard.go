package classpath
import "os"
import "path/filepath"
import "strings"

// 用*通配符  匹配目录下全部jar包
func newWildcardEntry(path string) CompositeEntry {
	baseDir := path[:len(path)-1] // remove *

	compositeEntry := []Entry{}

	//遍历baseDir创建ZipEntry
	walkFn := func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		//跳过子目录
		if info.IsDir() && path != baseDir {
			return filepath.SkipDir
		}

		//选出后缀是.jar的
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") {
			jarEntry := newZipEntry(path)
			compositeEntry = append(compositeEntry, jarEntry)
		}
		return nil
	}
	//执行
	filepath.Walk(baseDir, walkFn)

	return compositeEntry
}



