# Implement blockchain use DPoS consensus

## Build
### 二进制
sh compile/compile.sh bin
### 镜像
sh compile/compile.sh docker
## RUN 
```
git clone https://github.com/hellobchain/dpos.git
cd dpos
go run cmd/main.go
```

connect multi peer 
```
./dpos new --port 7051 --secio
```
## Vote
```
./dpos vote -name wsw -v 50
```

