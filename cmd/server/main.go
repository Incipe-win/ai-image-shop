// @title AI T-Shirt Shop API
// @version 1.0
// @description AI驱动的T恤设计商店API服务，提供用户认证、AI设计生成等功能
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for API authentication

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Incipe-win/ai-tshirt-shop/internal/handler"
	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
	"github.com/Incipe-win/ai-tshirt-shop/pkg/logger"
	"github.com/spf13/viper"
	_ "github.com/Incipe-win/ai-tshirt-shop/docs"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func main() {
	initConfig()

	env := viper.GetString("server.env")
	if env == "" {
		env = "development"
	}

	logger.Init(env)
	defer logger.Sync()

	db, err := repository.InitDatabase()
	if err != nil {
		logger.Fatal("Failed to initialize database", err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	// Try to auto migrate with the current connection
	err = db.AutoMigrate(&model.User{}, &model.Design{})
	if err != nil {
		logger.Error("Failed to auto migrate tables", err)
		logger.Fatal("Please ensure PostgreSQL is running and execute the following SQL commands manually:\n\n" +
			"CREATE USER tshirt WITH PASSWORD 'tshirt';\n" +
			"CREATE DATABASE tshirt_db;\n" +
			"GRANT ALL PRIVILEGES ON DATABASE tshirt_db TO tshirt;\n" +
			"\\c tshirt_db\n" +
			"GRANT ALL PRIVILEGES ON SCHEMA public TO tshirt;\n" +
			"GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO tshirt;")
	} else {
		logger.Info("Database tables migrated successfully")
	}

	// Initialize design repository after database is ready
	handler.InitDesignRepository(db)

	r := handler.InitRouter(env)

	port := viper.GetString("server.port")
	if port == "" {
		port = ":8080"
	}

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		logger.Info("Starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("Server exited")
}
