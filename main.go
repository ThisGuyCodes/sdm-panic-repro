package main

import (
	_ "github.com/athisguycodes/sdm-panic-repro/biglock"
	"github.com/awnumar/memguard"
)

func init() {
	memguard.CatchInterrupt()
}

func main() {
}
