package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("sandbox", "--lang=c", "--time=1000", "--memory=10000", "-c", "-s code/0c8a5c307c5911e4b3163c970e1ef66b/tmp.c", "-b code/0c8a5c307c5911e4b3163c970e1ef66b/tmp", "-i problem/afa58e907c4611e4a6703c970e1ef66b/inputTest", " -o problem/afa58e907c4611e4a6703c970e1ef66b/outputTest")
	out, e := cmd.CombinedOutput()
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Printf("11 %s", out)
	}
	fmt.Println("hello world")
}
