package heap
import "unicode/utf16"


//用map来表示字符串池，key是Go字符串，value是Java字符串。
var internedStrings = map[string]*Object{}

//给一个string 还一个string的java对象
//得内容是goStr的 String实例
func JString(loader *ClassLoader, goStr string) *Object {

	//已经在池中，直接返回即可
	if internedStr, ok := internedStrings[goStr]; ok {
		return internedStr
	}

	//Go字符串（UTF8格式）转换成Java字符数组（UTF16格式）
	chars := stringToUtf16(goStr)
	charArrClass := loader.LoadClass("[C")
	jChars := &Object{//一个char[]实例
		charArrClass,
		chars,nil}

	//string类
	stringClass := loader.LoadClass("java/lang/String")
	//new 一个String对象
	stringObject := stringClass.NewObject()
	//把它的value字段设置成刚刚转换而来的字符数组
	stringObject.SetRefVar("value", "[C", jChars)

	//放入池中
	internedStrings[goStr] = stringObject

	return stringObject
}

func stringToUtf16(s string) []uint16 {
	//Go语言字符串在内存中是UTF8编码的，先把它强制转成 UTF32
	runes := []rune(s) // utf32
	//调用utf16包的Encode（）函数编码成UTF16
	return utf16.Encode(runes)
}

//java的string对象 转 go字符串
func GoString(stringObject *Object) string {
	//拿到String对象的value变量值 是一个char[]
	charArr := stringObject.GetRefVar("value", "[C")
	//字符数组转换成Go字符串
	return utf16ToString(charArr.Chars())
}


func utf16ToString(s []uint16) string {
	//UTF16数据转换成UTF8编码
	runes := utf16.Decode(s) // utf8
	//强制转换成Go字符串
	return string(runes)
}


func InternString(jStr *Object) *Object {
	goStr := GoString(jStr)
	if internedStr, ok := internedStrings[goStr]; ok {
		return internedStr
	}
	//如果字符串还没有入池，把它放入并返回该字符串
	internedStrings[goStr] = jStr
	return jStr
}


