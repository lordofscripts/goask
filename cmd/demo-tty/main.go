package main

import (
	"fmt"
	"time"

	"github.com/lordofscripts/goask/tty"
)

func main() {
	tty.Clear()
	tty.Cursor(15, 15)
	fmt.Println("At line 15")
	tty.CursorUp(5)
	tty.Bold()
	tty.Yellow("This is Yellow\n")
	tty.BoldOff()
	tty.Purple("This is Purple\n")
	tty.BrightRed("Alert")
	fmt.Print("This is a long line of text")
	tty.CursorLeft(5)
	fmt.Println("Normal")

	const MAX int = 100
	for style := range [3]tty.ProgressStyle{tty.ProgressStyle1, tty.ProgressStyle2, tty.ProgressStyle3} {
		for i := range MAX + 1 {
			tty.ShowProgressAt(tty.ProgressStyle(style), 20+style, "Something", uint(i))
			time.Sleep(50 * time.Millisecond)
		}
	}
}
