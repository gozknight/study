package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {

	// go run code-user/main.go
	cmd := exec.Command("go", "run", "code-user/main.go")
	var out, stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr
	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
		return
	}
	io.WriteString(pipe, "23 11\n")
	// 根据输入的测试案例，拿到输入结果和标准输出结果对比
	if err := cmd.Run(); err != nil {
		log.Fatalln(err, stdErr.String())
		return
	}
	fmt.Println(out.String() == "34\n")
}
