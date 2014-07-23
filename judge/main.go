package main

import (
	"github.com/ggaaooppeenngg/OJ/app/judge"
)

func main() {
	go judge.GetHandledCodeLoop()
	judge.HandleCodeLoop()
}
