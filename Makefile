RED  =  "\e[31;1m"
GREEN = "\e[32;1m"

.PHONY: idl
.ONESHELL:
idl: $(idlFile)
	@ python3 ./gen.py


.ONESHELL:
run:
	./run.sh
	go install github.com/golang/mock/mockgen@v1.6.0
	docker compose up
	go run ./service/api/main.go
	go run ./service/user/main.go


.PHONY: build
build:
	go build  -o ./build/api ./service/api/.
	go build -o ./build/user  ./service/user/.
	go build -o ./build/video ./service/video/.
