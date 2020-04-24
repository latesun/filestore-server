## go-micro测试grpc通信

### 安装依赖工具包

```
sudo apt-get -y install autoconf automake libtool
```

### 安装protobuf

```
mkdir ./tmp && cd ./tmp
git clone https://github.com/google/protobuf
cd protobuf
./autogen.sh
./configure
make
sudo make install
```

### 安装go的grpc相关包

```
# github.com/micro/micro测试grpc通信
# go get github.com/micro/go-web
go get -v github.com/micro/protobuf/{proto,protoc-gen-go}
go get -v github.com/micro/protoc-gen-micro

export LD_LIBRARY_PATH=/usr/local/lib
export PATH=$GOPATH/bin:$PATH
```

### 生成go版的proto

假设在`/data/test/proto/` 下有一个proto文件：

```proto
syntax = "proto3";

package go.micro.service.user;

service UserService {
    // 用户注册
    rpc Signup(ReqSignup) returns (RespSignup) {}
    // 用户登录
    rpc Signin(ReqSignin) returns (RespSignin) {}
}

message ReqSignup {
    string username = 1;
    string password = 2;
}

message RespSignup {
    int32 code = 1;
    string message = 2;
}

message ReqSignin {
    string username = 1;
    string password = 2;
}

message RespSignin {
    int32 code = 1;
    string token = 2;
    string message = 3;
}
```

尝试以下操作来生成golang版本的代码:

```
protoc --proto_path=proto --proto_path=/data/test/ --micro_out=proto --go_out=proto proto/user.proto
```

正常情况，会生成两个go文件:

```
xiaomo@xiaomo:/data/test$ ls proto/ -l
总用量 20
-rw-r--r-- 1 xiaomo xiaomo 3252 3月   8 23:41 user.micro.go
-rw-r--r-- 1 xiaomo xiaomo 8272 3月   8 23:41 user.pb.go
-rw-r--r-- 1 xiaomo xiaomo  520 3月   8 23:40 user.proto

```