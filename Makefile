all: compile strip

compile:
	go build -ldflags="-s -w" .

strip:
	strip svm

install:
	install -Dvm 755 ./svm /usr/bin/svm