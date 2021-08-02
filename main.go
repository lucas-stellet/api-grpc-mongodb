package main

import (
	"fmt"
	"lucas-stellet/api-grpc-mongodb/pkg/db"
	"lucas-stellet/api-grpc-mongodb/pkg/env"
	"lucas-stellet/api-grpc-mongodb/pkg/grpc"
	"lucas-stellet/api-grpc-mongodb/pkg/logger"
)

func init() {
	env.Load(".env")

	err := db.StartConnection(env.DSN, "internal-gateway-v2")

	if err != nil {
		logger.Write(logger.FATAL, fmt.Sprintf("error when starting mongo connection :: %v", err), "file")
	}

	logger.Write(logger.INFO, "mongodb connected", "stdout")
}

func main() {
	grpc.Start()
}
