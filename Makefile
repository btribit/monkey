build:
	go build -o monkey main/main.go

extensions:
	go build -buildmode=plugin -o extensions/hello.so extensions/hello.go 

evaluator: build
	./monkey

compiler: build
	./monkey compiler