package main

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"os"
	"time"
)

//go:embed deployers/certDeploy.exe
var windowsExe []byte

//go:embed deployers/certDeployLinux
var linuxExe []byte

func PatchDeployer(exeFile []byte, cert []byte, outputFile string) {
	originalSize := bytes.NewBuffer(nil)
	size := uint64(len(exeFile))
	binary.Write(originalSize, binary.LittleEndian, size)

	output := append(exeFile, cert...)
	output = append(output, []byte("payload ")...)
	output = append(output, originalSize.Bytes()...)

	err := os.WriteFile("builds/"+outputFile, output, 0777)
	if err != nil {
		fmt.Println("error writing output file:", err)
		return
	}
}

var certBytes []byte
var exe []byte

func main() {
	cert := binding.BindBytes(&certBytes)
	_, err := os.Stat("builds")
	if err != nil {
		err = os.Mkdir("builds", 0777)
		if err != nil {
			fmt.Println("error creating builds directory:", err)
			os.Exit(1)
		}
	}
	a := app.New()
	w := a.NewWindow("Certificate Deployer")
	w.Resize(fyne.NewSize(800, 600))
	certButton := widget.NewButton("Select Certificate", func() {
		fileDialog := dialog.NewFileOpen(func(r fyne.URIReadCloser, _ error) {
			if r == nil {
				return
			}
			data, _ := ioutil.ReadAll(r)
			p, _ := pem.Decode(data)
			if p == nil {
				dialog.ShowError(errors.New("Invalid Certificate"), w)
				return
			}
			err := cert.Set(data)
			if err != nil {
				panic(err)
			}
		}, w)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".pem", ".crt", ".cer"}))
		fileDialog.Show()

	})
	buildButton := widget.NewButton("Build", func() {
		if certBytes == nil {
			dialog.ShowError(errors.New("Please select a certificate"), w)
			return
		}
		certData, _ := cert.Get()

		curDate := time.Now().Format("2006-01-02")

		PatchDeployer(windowsExe, certData, curDate+"-certDeploy-patched.exe")
		PatchDeployer(linuxExe, certData, curDate+"-certDeployLinux-patched")
		dialog.NewInformation("Cert Deployer", "Build Complete!", w).Show()
	})
	w.SetContent(container.NewVBox(certButton, buildButton))
	w.ShowAndRun()

}
