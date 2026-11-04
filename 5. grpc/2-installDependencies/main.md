1. 在 mac 中保存 protoc-36.2-osx-aarch_64 文件夹

2. 执行命令(一次性终端): 
    export PATH="/Users/vvyh/Desktop/Go/5. grpc/protoc-36.2-osx-aarch_64/bin:/Users/vvyh/go/bin:$PATH"

3. 检查是否可用: protoc --version  正常会输出 libprotoc 36.2

4. 到 .proto 文件所在的文件夹下执行命令将 .proto 文件转换成 Go 代码: 
   protoc --go_out=. --go_opt=paths=source_relative hello.proto