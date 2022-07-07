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
	err = os.WriteFile("/usr/local/share/ca-certificates/TrustedCert.crt", data, 0400)
	if err != nil {
		panic(err)
	}
	err = exec.Command("update-ca-certificates").Run()
	if err != nil {
		panic(err)
	}
}
