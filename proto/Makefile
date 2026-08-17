GO_MODULE := github.com/Acad600-TPA/WEB-WE-251/proto
.PHONY: protoc-go
protoc-go:
	protoc --go_opt=module=$(GO_MODULE) --go_out=. \
		--go-grpc_opt=module=$(GO_MODULE) --go-grpc_out=. \
		controller/galactus_controller.proto \
		controller/notification_controller.proto \
		controller/social_controller.proto \
		controller/stream_controller.proto \
		controller/video_controller.proto \
		dto/galactus.proto \
		dto/notification.proto \
		dto/social.proto \
		dto/video.proto \
		dto/stream.proto \
