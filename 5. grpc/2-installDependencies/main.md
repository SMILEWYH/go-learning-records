1. 在 mac 中保存 protoc-36.2-osx-aarch_64 文件夹, 和安装以下四个依赖
    go get github.com/golang/protobuf/proto
    go get google.golang.org/grpc
    go install github.com/golang/protobuf/protoc-gen-go
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

2. 执行命令(一次性终端): 
    export PATH="/Users/vvyh/Desktop/Go/5. grpc/protoc-36.2-osx-aarch_64/bin:/Users/vvyh/go/bin:$PATH"

3. 检查是否可用: protoc --version  正常会输出 libprotoc 36.2

4. 到 .proto 文件所在的文件夹下执行命令将 .proto 文件转换成 Go 代码:
   // 定义请求、响应的数据类型, 以及序列化、反序列化所需的支持代码
   protoc --go_out=. --go_opt=paths=source_relative hello.proto
   // 提供服务端、客户端的注册和调用函数方法、还有接口
   protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative hello.proto