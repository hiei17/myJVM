package main
//导入go的内置包
import "os"//os包定义了一个Args变量，其中存放传递给命令行的全部参数
import "flag"//帮助我们处理命令行选项
import "fmt"


type Cmd struct {//定义结构体
	helpFlag        bool
	versionFlag     bool
	cpOption        string//用户类路径
	class           string//类名
	args            []string//main传参
	XjreOption      string//指定jre目录的位置//启动类路径和扩展类路径

	verboseClassFlag bool//是否把类加载信息输出到控制台
	verboseInstFlag bool//是否把指令执行信息输出到控制台
}
//Go语言有函数（Function）和方法（Method）之分，方法调用需要receiver，函数调用则不需要。
//这是函数
func parseCmd() *Cmd {
	//new 一个对象
	cmd := &Cmd{}

	//设置 如果Parse（）函数解析失败，它就调用printUsage打印
	flag.Usage = printUsage

	//调用flag包提供的各种Var（）函数
	// 设置	需要解析的选项
	flag.BoolVar(&cmd.helpFlag, "help", false, "print help message")
	flag.BoolVar(&cmd.helpFlag, "?", false, "print help message")
	flag.BoolVar(&cmd.versionFlag, "version", false, "print version and exit")
	flag.StringVar(&cmd.cpOption, "classpath", "", "classpath")
	flag.StringVar(&cmd.cpOption, "cp", "", "classpath")//用户类路径
	flag.StringVar(&cmd.XjreOption, "Xjre", "", "path to jre")//启动类路径和扩展类路径

	flag.BoolVar(&cmd.verboseClassFlag, "verbose", false, "enable verbose output")
	flag.BoolVar(&cmd.verboseClassFlag, "verbose:class", false, "enable verbose output")
	flag.BoolVar(&cmd.verboseInstFlag, "verbose:inst", false, "enable verbose output")

	//解析以上设置的选项
	flag.Parse()

	//如果解析成功， 捕获其他没有被解析的参数
	args := flag.Args()
	if len(args) > 0 {
		//参数放入Cmd结构体
		cmd.class = args[0]//第一个参数就是主类名
		cmd.args = args[1:]//剩下的是要传递给主类的参数
	}

	return cmd
}

func printUsage() {
	fmt.Printf("Usage: %s [-options] class [args...]\n", os.Args[0])
}


