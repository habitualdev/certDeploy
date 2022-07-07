package main

import "os"

func checkAdmin() {
	// check for root
	if os.Getuid() != 0 {
		panic("you must be root to run this")
	}
}
