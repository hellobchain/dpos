# Implement blockchain use DPoS consensus

## Build
go build -o build/dpos  cmd/main.go
## RUN 
```
git clone https://github.com/wsw365904/dpos.git
cd dpos
go build cmd/main.go
```

connect multi peer 
```
./dpos new --port 7051 --secio
```
## Vote
```
./dpos vote -name wsw -v 50
```

