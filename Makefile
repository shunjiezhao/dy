RED  =  "\e[31;1m"
GREEN = "\e[32;1m"

.PHONY: Idl

idlPbFile = $(shell find . -name '*.proto') # idl path
idlTrFile = $(shell find . -name '*.thrift') # idl path
moudle = "first" # go.mod name
# 生成的文件
genFile = build.sh main.go kitex.yaml handler.go ./script
# 将生成的文件移动到目录下
mvGenFile = cd  ./service/$(1)/ && rm $(genFile) -rf && cd - &&  mv $(genFile) ./service/$(1)/
# a/b.c -> bname = b, file = b.c
# protobuf 生成命令
protoGen = (kitex -module $(moudle) -type protobuf -service $(basename $(notdir $(1))) -I $(shell dirname $(1)) $(1)   \
			&& $(call mvGenFile,$(basename $(notdir $(1)))) \
			&& echo  $(GREEN)$(1)" kitex generate success.") \
		|| echo  $(RED)$(1)" kitex generate fail." ;
#TODO: thrift 生成命令
thriftGen = echo  $(RED)"$(notdir $(1)) kitex generate fail." ;


.ONESHELL:
Idl: $(idlFile)
	@$(foreach name,$(idlPbFile), $(call protoGen, $(name)))
	@$(foreach name,$(idlTrFile), $(call thriftGen, $(name)))
cleanidl:
	@$(foreach name,$(idlPbFile), $(call mvGenFile,$(basename $(notdir $(name)))))


http:
	go run ./service/api/main.go
