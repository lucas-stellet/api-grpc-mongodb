# Compile stage

FROM golang:alpine AS build_stage
WORKDIR /go/src/app
RUN apk add make protoc
COPY . .
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
  go get google.golang.org/protobuf/cmd/protoc-gen-go && \
  make gen && \
  make gobuild

# Final stage

FROM alpine:latest
EXPOSE 50051
WORKDIR /
COPY --from=build_stage /go/src/app/build/api .

CMD ["/api"]