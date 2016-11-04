projectpath = ${PWD}
glidepath = ${PWD}/vendor/github.com/Masterminds/glide
wrkpath = ${PWD}/vendor/github.com/wg/wrk
wrk2path = ${PWD}/vendor/github.com/giltene/wrk2
protobufcpath = ${PWD}/vendor/github.com/google/protobuf

target:
	@go build

test:
	@go test

bench:
	@go test -benchmem -bench=.

$(projectpath)/glide:
	git clone https://github.com/Masterminds/glide.git $(glidepath)
	cd $(glidepath);make build
	cp $(glidepath)/glide .

wrk:
	if [ ! -d "$(wrkpath)" ]; then git clone https://github.com/wg/wrk.git $(wrkpath); fi
	cd $(wrkpath);make
	cp $(wrkpath)/wrk ./wrk

wrk2:
	if [ ! -d "$(wrk2path)" ]; then git clone https://github.com/giltene/wrk2.git $(wrk2path); fi
	cd $(wrk2path);make
	cp $(wrk2path)/wrk ./wrk2

$(protobufcpath)/bin/protoc:
	go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
	if [ ! -d "$(protobufcpath)" ]; then mkdir -p $(protobufcpath); fi
	wget https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-osx-x86_64.zip -P $(protobufcpath)
	unzip $(protobufcpath)/protoc-3.1.0-osx-x86_64.zip -d $(protobufcpath)

$(projectpath)/grpc/contract.pb.go:
	$(protobufcpath)/bin/protoc -I grpc/ grpc/contract.proto --go_out=plugins=grpc:grpc

deps: $(projectpath)/glide wrk wrk2 $(protobufcpath)/bin/protoc $(projectpath)/grpc/contract.pb.go
	$(projectpath)/glide install
