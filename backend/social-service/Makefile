GO_MODULE := "github.com/Acad600-TPA/WEB-WE-251/backend/social-service"

.PHONY:protoc-go
protoc-go:
	protoc --proto_path=../../proto --go_opt=module=${GO_MODULE} --go_out=. \
        	--go-grpc_opt=module=${GO_MODULE} --go-grpc_out=. \
    		controller/social_controller.proto dto/social.proto
