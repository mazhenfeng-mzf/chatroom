package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	LogServer string = "logServer"
	LogClient string = "logClient"
)

type LOG struct {
	dir     string
	logFile string
	f       *os.File
}

var MyLOG LOG

func (this *LOG) Init(dir string) {
	//dir: LogServer|LogClient
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	fileTimeStr := strings.Replace(strings.Replace(timeStr, ":", "-", -1), " ", "-", -1)
	this.logFile = fmt.Sprintf("C://Users/MZF/go/src/chatroom/main/%s/clientLog-%s.txt", dir, fileTimeStr)

	if checkFileIsExist(this.logFile) { //如果文件存在
		this.f, _ = os.OpenFile(this.logFile, os.O_APPEND, 0666) //打开文件
		fmt.Printf("日志文件存在 %s\n", this.logFile)
	} else {
		fmt.Println("日志文件不存在")
		this.f, _ = os.Create(this.logFile)                      //创建文件
		this.f, _ = os.OpenFile(this.logFile, os.O_APPEND, 0666) //打开文件
		fmt.Printf("打开日志文件 %s\n", this.logFile)
	}
}

func (this *LOG) End() {
	this.f.Close()
}

func (this *LOG) writeFile(writeString string) {
	io.WriteString(this.f, writeString) //写入文件(字符串)

}

func (this *LOG) Log(formating string, args ...interface{}) {

	////////////////////
	srcfilename, line, funcname := "???", 0, "???"
	pc, srcfilename, line, ok := runtime.Caller(1)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
		srcfilename = filepath.Base(srcfilename)     // /full/path/basename.go => basename.go
	}
	/////////////////////

	writeString := fmt.Sprintf("<info>-<%s><%d><func %s>: %s\n", srcfilename, line, funcname, fmt.Sprintf(formating, args...))
	fmt.Println(writeString)
	this.writeFile(writeString)
}

func (this *LOG) ErrLog(formating string, args ...interface{}) {

	////////////////////
	srcfilename, line, funcname := "???", 0, "???"
	pc, srcfilename, line, ok := runtime.Caller(1)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
		srcfilename = filepath.Base(srcfilename)     // /full/path/basename.go => basename.go
	}
	/////////////////////

	writeString := fmt.Sprintf("<error>-<%s><%d><func %s>: %s\n", srcfilename, line, funcname, fmt.Sprintf(formating, args...))
	fmt.Println(writeString)
	this.writeFile(writeString)
}

// func (this *LOG) Log2(formating string, args ...interface{}) {
// 	var funcname, srcfilename string
// 	var codeline int
// 	for skip := 1; true; skip++ {
// 		pc, srcfilename, line, ok := runtime.Caller(skip)
// 		if !ok {
// 			// 不ok，函数栈用尽了
// 			break
// 			// auto.Code = prevCode
// 			// auto.Func = prevFunc
// 			// return auto
// 		} else {
// 			funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
// 			funcname = filepath.Ext(funcname)            // .foo
// 			funcname = strings.TrimPrefix(funcname, ".") // foo
// 			srcfilename = filepath.Base(srcfilename)     // /full/path/basename.go => basename.go
// 			codeline = line
// 		}
// 	}
// 	writeString := fmt.Sprintf("%s:%d:%s: %s\n", srcfilename, codeline, funcname, fmt.Sprintf(formating, args...))
// 	fmt.Println(writeString)
// 	this.writeFile(writeString)
// }

func (this *LOG) WritePanic(err error) {
	if err == nil {
		return
	}

	////////////////////
	srcfilename, line, funcname := "???", 0, "???"
	pc, srcfilename, line, ok := runtime.Caller(2)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
		srcfilename = filepath.Base(srcfilename)     // /full/path/basename.go => basename.go
	}
	/////////////////////

	writeString := fmt.Sprintf("!!!!!! Panic:\n%s:%d:%s: \n%s\n", srcfilename, line, funcname, fmt.Sprintf("%v", err))
	fmt.Println(writeString)
	this.writeFile(writeString)
	panic(err)
}
