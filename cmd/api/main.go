package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"venturex-backend/internal/app/config"
	"venturex-backend/internal/pkg/database"
	"venturex-backend/internal/pkg/jwt"
	"venturex-backend/internal/pkg/logger"
	"venturex-backend/internal/pkg/password"

	"venturex-backend/internal/app/api"
	v1 "venturex-backend/internal/app/api/v1"
	"venturex-backend/internal/app/repositories"
	"venturex-backend/internal/app/service"
	"venturex-backend/internal/app/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.New()

	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "config/app.example.yaml"
	}

	appCfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config file (%s): %v", configPath, err)
	}

	fmt.Printf("Loaded DB config: %+v\n", appCfg.Database)

	db, err := database.New(appCfg.Database.ToDBConfig())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Database connection established successfully")

	txManager := database.NewTransactionManager(db)

	passwordService := password.NewPasswordService(password.DefaultPasswordConfig())

	jwtCfg, err := jwt.DefaultJWTConfig()
	if err != nil {
		log.Fatal("Failed to create default JWT config:", err)
	}
	if appCfg.JWT != nil {
		if appCfg.JWT.AccessTokenDuration > 0 {
			jwtCfg.AccessTokenDuration = appCfg.JWT.AccessTokenDuration
		}
		if appCfg.JWT.RefreshTokenDuration > 0 {
			jwtCfg.RefreshTokenDuration = appCfg.JWT.RefreshTokenDuration
		}
		if appCfg.JWT.Issuer != "" {
			jwtCfg.Issuer = appCfg.JWT.Issuer
		}
		if appCfg.JWT.Audience != "" {
			jwtCfg.Audience = appCfg.JWT.Audience
		}
	}

	jwtService := jwt.NewJWTService(jwtCfg)

	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	loginAttemptRepo := repositories.NewLoginAttemptRepository(db)

	serviceRegistry := service.NewServiceRegistry(
		userRepo,
		sessionRepo,
		loginAttemptRepo,
		jwtService,
		passwordService,
	)

	authUseCase := usecase.NewAuthUseCase(
		serviceRegistry.AuthService(),
		serviceRegistry.JWTService(),
		serviceRegistry.PasswordService(),
		userRepo,
		sessionRepo,
		txManager,
	)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(api.LoggerMiddleware())
	router.Use(api.RequestID())
	router.Use(api.CORSMiddleware())
	router.Use(api.SecurityHeaders())
	router.Use(api.ValidateJSONMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		})
	})

	authHandler := v1.NewAuthHandler(authUseCase, serviceRegistry.AuthService())

	routerRegister := api.NewGinRouterRegisterImpl(router)

	if err := authHandler.Register(routerRegister); err != nil {
		log.Fatal("Failed to register auth handler:", err)
	}

	port := appCfg.Server.Port
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	go func() {
		fmt.Printf("🚀 Server starting on port %s\n", port)
		fmt.Printf("📋 API Documentation: http://localhost:%s/\n", port)
		fmt.Printf("🔍 Health Check: http://localhost:%s/health\n", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("✅ Server exited gracefully")
}

// Architecture layers:
// 1. Handler Layer (api/v1.*Handler) - HTTP handling
// 2. UseCase Layer (usecase.*UseCase) - Business logic
// 3. Service Layer (service.*Service) - Domain services & shared logic
// 4. Repository Layer (repositories.*) - Data access
// 5. Domain Layer (domain.*) - Business entities & value objects
//
// Flow: Handler → UseCase → Service → Repository
//       Handler ← UseCase ← Service ← Repository
//
// Rules:
// - Handler can only call UseCase
// - UseCase can call Service and Repository
// - Service can call Service and Repository
// - Repository can only call Domain
// - All layers communicate through interfaces
// - No layer can call a layer above it

// Production deployment:
// CONFIG_FILE=config/app.example.yaml go run ./cmd/api
