VERSION:=0.4
EXE:=xccat

all: windows linux_amd64 release

linux_amd64:
	mkdir -p linux
	GOOS="linux" GOARCH="amd64" go build
	mv ${EXE} linux/

windows:
	mkdir windows
	GOOS="windows" GOARCH="amd64" go build
	mv ${EXE}.exe windows/

release:
	zip -r ${EXE}-${VERSION}.zip linux/ windows/

clean:
	@rm -rf linux
	@rm -rf windows
	@rm -f ${EXE}*.zip
