BUILD=go build

SOURCES=scorsh.go \
types.go \
config.go \
spooler.go \
commits.go \
workers.go

all: scorsh

deps:
	go get 'github.com/fsnotify/fsnotify'
	go get 'github.com/libgit2/git2go'
	go get 'github.com/go-yaml/yaml'
	go get 'golang.org/x/crypto/openpgp'

scorsh: $(SOURCES)
	$(BUILD) scorsh.go types.go config.go spooler.go commits.go workers.go

clean:
	rm scorsh
