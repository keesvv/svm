all: compile strip

compile:
	go build -ldflags="-s -w" .

strip:
	strip svm