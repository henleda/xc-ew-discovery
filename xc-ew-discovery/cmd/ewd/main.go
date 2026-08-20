// Command ewd is the development entry point. It runs the pipeline locally so
// each milestone has a runnable definition of done.
//
//	ewd sensor --stdout      M1
//	ewd sensor --graph       M2
//	ewd infer --out ./specs  M3
//	ewd inventory            M4 and M5
//	ewd collector            M6
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: ewd <sensor|infer|inventory|collector>")
		os.Exit(2)
	}
	// TODO(M1): wire subcommands as milestones land.
	fmt.Fprintf(os.Stderr, "subcommand %q not implemented yet, see docs/MILESTONES.md\n", os.Args[1])
	os.Exit(1)
}
