BUILD=go build


SOURCES=scorsh.go \
types.go \
config.go \
spooler.go \
commits.go \
workers.go

all: scorsh


scorsh: $(SOURCES) 
	$(BUILD) scorsh.go types.go config.go spooler.go commits.go workers.go

clean:
	rm scorsh
