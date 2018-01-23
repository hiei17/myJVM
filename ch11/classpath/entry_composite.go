package classpath
import "errors"
import "strings"
//;隔开的复合路径
type CompositeEntry []Entry//更小的Entry组成

//构造函数把参数（路径列表）按分隔符分成小路径，然后把每个小路径都转换成具体的Entry实例，代码如下：
func newCompositeEntry(pathList string) CompositeEntry {
	compositeEntry := []Entry{}
	for _, path := range strings.Split(pathList, pathListSeparator) {
		entry := newEntry(path)
		compositeEntry = append(compositeEntry, entry)
	}
	return compositeEntry
}

func (self CompositeEntry) readClass(className string) ([]byte, Entry, error) {
	for _, entry := range self {
		data, from, err := entry.readClass(className)
		//报错不管它

		//有个不报错就成了 返回
		if err == nil {
			return data, from, nil
		}
	}
	//遍历完都不成 返回错误
	return nil, nil, errors.New("class not found: " + className)
}

//调用每一个子路径的String（）方法，然后把得到的字符串用路径分隔符拼接起来即可
func (self CompositeEntry) String() string {
	strs := make([]string, len(self))
	for i, entry := range self {
		strs[i] = entry.String()
	}
	return strings.Join(strs, pathListSeparator)
}
