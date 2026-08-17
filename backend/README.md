MakeFile Protoc Command
---
SETUP FILE (in backend)
cd /Documents/work/TPA/web/WEB-WE-251/backend/
mkdir brainrot-service
cd brainrot-service
mkdir config internal proto proto/controller proto/dto utils internal/adapters internal/adapters/inbound/controller 
internal/adapters/inbound/dto internal/adapters/inbound/grpc internal/adapters/outbound/authentication 
internal/adapters/outbound/cache internal/adapters/outbound/db internal/adapters/outbound/repository app/domain/ 
app/helper ports/in ports/out
touch Dockerfile Makefile .air.toml main.go
--- 
*in makefile:*

GO_MODULE := "github.com/Wenev/Survace/*"
.PHONY:protoc-go
protoc-go:
    protoc --go_opt=module={GO_MODULE} --go_out=. \
        --go-grpc_opt=module=${GO_MODULE} --go-grpc_out=. \
        ./proto/*_controller.proto ./proto/*.proto

---
GORM & Go Setup
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest****
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres

---
export GOPRIVATE=github.com/Wenev/Survace/proto
export GONOSUMDB=github.com/Wenev/Survace/proto
go get github.com/Wenev/Survace/proto