package env

import (
	"fmt"
	"lucas-stellet/api-grpc-mongodb/pkg/logger"
	"os"

	"github.com/joho/godotenv"
)

var (
	// Port is the PORT environment variable or 8080 if missing.
	// Used to open the tcp listener for our web server.
	Port string
	// DSN is the DSN environment variable or mongodb://localhost:27017 if missing.
	// Used to connect to the mongodb.
	DSN string
	// Enviroment is sets the enviroment mode that application is running.
	// If DEV, the environment variables will be shown in terminal.
	Enviroment string
	// Token ...
	Token string
)

// Load ...
func Load(envFile string) {
	if fileExists(envFile) {
		logger.Write(logger.INFO, "Loading environment variables from file", "stdout")

		if err := godotenv.Load(envFile); err != nil {
			logger.Write(logger.FATAL, fmt.Sprintf("error loading environment variables from [%s]: %v", envFile, err), logger.FILE)
			panic("error loading environment variables")
		}
	}

	parse()

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func parse() {
	Port = getDefault("PORT")
	DSN = getDefault("DSN")
	Enviroment = getDefault("ENVIRONMENT")
	Token = getDefault("TOKEN")

	switch Enviroment {
	case "DEV":
		logger.Write(logger.INFO, fmt.Sprintf("• Port=%s", Port), "stdout")
		logger.Write(logger.INFO, fmt.Sprintf("• DSN=%s", DSN), "stdout")
		logger.Write(logger.INFO, fmt.Sprintf("• Enviroment=%s", Enviroment), "stdout")
		logger.Write(logger.INFO, fmt.Sprintf("• Token=%s", Token), "stdout")
	default:
		logger.Write(logger.INFO, fmt.Sprintf("• Port=%s", Port), "file")
		logger.Write(logger.INFO, fmt.Sprintf("• DSN=%s", DSN), "file")
		logger.Write(logger.INFO, fmt.Sprintf("• Enviroment=%s", Enviroment), "file")
	}
}

func getDefault(key string) string {
	value := os.Getenv(key)

	if key == "ENVIRONMENT" {
		return "PRODUCTION"
	}

	return value
}
