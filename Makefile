BINARIES = scribe
GOBIN=$(or ${GOPATH},~/go)/bin
CMD_PACKAGE_PATH = $(PWD)/cmd

all: $(BINARIES)

$(BINARIES): setup
	@mkdir -p bin
	go build -o bin/$@ $(CMD_PACKAGE_PATH)/$@
	cp bin/$@ $(GOBIN)/$@

setup:
	if [ ! -d $(GOBIN) ] ; then mkdir $(GOBIN); fi

clean:
	rm -rf bin/
