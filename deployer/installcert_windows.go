package main

import (
	"github.com/gonutz/payload"
	"os"
	"os/exec"
)

func installCert() {
	data, err := payload.Read()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("TrustedCert.pem", data, 0666)
	if err != nil {
		panic(err)
	}
	err = exec.Command("certutil", "-addstore", "Root", "TrustedCert.pem").Run()
	if err != nil {
		panic(err)
	}
}
