protoc -I ./redirect/ -I ../../../ --go_out=plugins=grpc:./redirect/ ./redirect/redirect.proto
