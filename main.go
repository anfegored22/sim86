package main

import (
	"bufio"
	"emulator/emulator"
	"flag"
	"fmt"
	"os"
)

// Encoded as regreg
// Should I do something like wregreg
func main() {
	fn := flag.String("file", "", "Name of the assembly file")
	flag.Parse()
	if *fn == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: -file")
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.ReadFile(*fn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cpu := emulator.NewCPU()
	cpu.Print()

	instructions := emulator.Decoder(f)
	reader := bufio.NewReader(os.Stdin)
	for _, inst := range instructions {
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := cpu.Execute(inst); err != nil {
			os.Exit(1)
		}
		cpu.Update(inst) // This only prints it doesn't change state
	}
}
