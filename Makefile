build:
	cd deployer && GOOS=linux go build -ldflags "-s -w" -o certDeployLinux && GOOS=windows go build -ldflags "-s -w" -o certDeploy.exe
	mv deployer/certDeploy* builder/deployers/
	cd builder && GOOS=linux go build -ldflags "-s -w" -o builderLinux && GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags "-s -w" -o builder.exe
	rm -rf bin
	mkdir bin
	mv builder/builder* bin
	rm -rf builder/deployers/certDeploy*
