RED  =  "\e[31;1m"
GREEN = "\e[32;1m"

.PHONY: Idl
#TODO:完成 main.go handler.go 存在的情况
mvFile =  (test -r $(1) || mv $(1))


idlPbFile = $(shell find . -name '*.proto') # idl path
idlTrFile = $(shell find . -name '*.thrift') # idl path
Moudle = "first" # go.mod name
# 生成的文件
genFile = build.sh kitex.yaml  ./script
# 将生成的文件移动到目录下
mvGenFile = cd  ./service/$(1)/ && rm $(genFile) -rf && cd - &&  mv $(genFile) ./service/$(1)/
# protobuf 生成命令
protoGen = (kitex -module $(Moudle) -type protobuf -service $(basename $(notdir $(1))) -I $(shell dirname $(1)) $(1)   \
			&& $(call mvGenFile,$(basename $(notdir $(1)))) \
			&& echo  $(GREEN)$(1)" kitex generate success.") \
		|| echo  $(RED)$(1)" kitex generate fail." ;
#TODO: thrift 生成命令
thriftGen = echo $(RED)"$(notdir $(1)) TODO: thrift 生成命令." ;


.ONESHELL:
Idl: $(idlFile)
	$(foreach name,$(idlPbFile), $(call protoGen, $(name)))
	@$(foreach name,$(idlTrFile), $(call thriftGen, $(name)))
cleanidl:
	@$(foreach name,$(idlPbFile), $(call mvGenFile,$(basename $(notdir $(name)))))


.ONESHELL:
run:
	go install github.com/golang/mock/mockgen@v1.6.0
	docker compose up
	go run ./service/api/main.go
	go run ./service/user/main.go


.PHONY: build
build:
	go build ./service/api/.
	go build ./service/user/.
