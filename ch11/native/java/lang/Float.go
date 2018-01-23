package lang

import "math"
import "jvmgo/ch11/native"
import "jvmgo/ch11/rtda"
/*

Math类在初始化时需要调用
Float.floatToRawIntBits（）和
Double.doubleToRawLongBits（）
package java.lang;
public final class Math {
	// Use raw bit-wise conversions on guaranteed non-NaN arguments.
	private static long negativeZeroFloatBits = Float.floatToRawIntBits(-0.0f);
	private static long negativeZeroDoubleBits = Double.doubleToRawLongBits(-0.0d);
}
*/

const jlFloat = "java/lang/Float"

func init() {
	native.Register(jlFloat, "floatToRawIntBits", "(F)I", floatToRawIntBits)
	native.Register(jlFloat, "intBitsToFloat", "(I)F", intBitsToFloat)
}

// public static native int floatToRawIntBits(float value);
func floatToRawIntBits(frame *rtda.Frame) {
	value := frame.LocalVars().GetFloat(0)
	bits := math.Float32bits(value)
	frame.OperandStack().PushInt(int32(bits))
}

// public static native float intBitsToFloat(int bits);
// (I)F
func intBitsToFloat(frame *rtda.Frame) {
	bits := frame.LocalVars().GetInt(0)
	value := math.Float32frombits(uint32(bits)) // todo
	frame.OperandStack().PushFloat(value)
}
