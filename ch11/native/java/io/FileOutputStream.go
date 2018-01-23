package io

import "os"
import "unsafe"
import "jvmgo/ch11/native"
import "jvmgo/ch11/rtda"

const fos = "java/io/FileOutputStream"

func init() {
	native.Register(fos, "writeBytes", "([BIIZ)V", writeBytes)
}


//最终打印的方法
// private native void writeBytes(byte b[], int off, int len, boolean append) throws IOException;
// ([BIIZ)V
func writeBytes(frame *rtda.Frame) {
	vars := frame.LocalVars()
	//this := vars.GetRef(0)
	b := vars.GetRef(1)
	off := vars.GetInt(2)
	len := vars.GetInt(3)
	//append := vars.GetBoolean(4)

	//Java语言中byte是有符号类型，在Go语言中byte则是无符号类型。
	//所以这里需要先把Java的字节数组转换成Go的[]byte变量
	jBytes := b.Data().([]int8)
	goBytes := castInt8sToUint8s(jBytes)

	goBytes = goBytes[off : off+len]

	//写到控制台。
	os.Stdout.Write(goBytes)
}



func castInt8sToUint8s(jBytes []int8) (goBytes []byte) {
	ptr := unsafe.Pointer(&jBytes)
	goBytes = *((*[]byte)(ptr))
	return
}

