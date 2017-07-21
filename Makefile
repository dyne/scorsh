BUILD=go build

SERVER_SOURCES=scorshd.go \
types.go \
config.go \
spooler.go \
commits.go \
workers.go \
exec.go

all: scorshd

deps:
	go get 'github.com/fsnotify/fsnotify'
	go get 'github.com/libgit2/git2go'
	go get 'gopkg.in/yaml.v2'
	go get 'golang.org/x/crypto/openpgp'

scorshd: $(SERVER_SOURCES)
	$(BUILD) $(SERVER_SOURCES)

clean:
	rm scorshd
