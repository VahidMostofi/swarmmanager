package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

// PythonPath is the path to python interpretor
const PythonPath = "/home/vahid/envs/data/bin/python"

// ScriptPath is the path to python script
const ScriptPath = "mobo.py"

type InteractiveInterface struct {
}

var wg sync.WaitGroup
var counter int

func (i *InteractiveInterface) Write(p []byte) (int, error) {
	// fmt.Println("from mobo:", string(p))
	v := &struct {
		Value int `json:"value"`
	}{}
	err := json.Unmarshal(p, v)
	if err != nil {
		// panic(err)
		fmt.Println(string(p))
	}
	fmt.Println(v.Value)
	v.Value = v.Value * v.Value

	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(stdin, string(b)+"\n")
	if err != nil {
		panic(err)
	}
	if v.Value == 81 {
		// time.Sleep(2 * time.Second)
		wg.Done()
	}
	return len(p), nil
}

var cmd *exec.Cmd
var stdin io.WriteCloser
var stdout io.ReadCloser

func main() {
	// var outValue []byte
	ctx, _ := context.WithCancel(context.Background())
	cmd = exec.CommandContext(ctx, PythonPath, ScriptPath)
	var err error

	stdin, err = cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer stdin.Close()

	// stdout, err = cmd.StdoutPipe()
	// if err != nil {
	// 	panic(err)
	// }
	// defer stdout.Close()

	cmd.Stdout = &InteractiveInterface{}
	cmd.Stderr = cmd.Stderr

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	// ===================

	wg.Add(1)
	// go func() {
	// 	outValue, _ = capture(os.Stdout, stdout)
	// 	wg.Done()
	// }()

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(outValue))
}
