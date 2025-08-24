package app

import (
	"PVZ/internal/cache/blackList"
	"PVZ/internal/config"
	"PVZ/internal/logger"
	pvzRepo "PVZ/internal/repository/pvz"
	receptionRepo "PVZ/internal/repository/reception"
	userRepo "PVZ/internal/repository/user"
	"PVZ/internal/service"
	authService "PVZ/internal/service/auth"
	pvzService "PVZ/internal/service/pvz"
	receptionService "PVZ/internal/service/reception"
	"PVZ/internal/utils"
	"PVZ/pkg/handler"
	"PVZ/pkg/routes"
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"time"
)

var logLevel = flag.String("1", "info", "log level")

type Servers struct {
	HTTP       *http.Server
	Redis      *redis.Client
	Prometheus *http.Server
	DB         *pgxpool.Pool
}

func SetupServers(ctx context.Context, cfg *config.Config) (*Servers, error) {
	logger.Init(getCore(getAtomicLevel()))

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	pool, err := initDB(ctx, dsn)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	redisConn, err := newRedisClient(ctx, cfg.Redis)
	if err != nil {
		logger.Error("Failed to create Redis client", zap.Error(err))
		return nil, err
	}

	tokenUtils := createTokenUtils(redisConn, cfg.JWT.SecretKey, cfg.JWT.Expiration)
	authSrv := createAuthService(redisConn, tokenUtils, pool)
	pvzSrv := createPVZService(pool)
	receiptSrv := createReceiptService(pool)

	initHandler := handler.NewPVZHandler(pvzSrv, receiptSrv, authSrv)
	ginEngine, err := routes.SetupRoutes(initHandler, tokenUtils)
	if err != nil {
		logger.Error("Failed to setup HTTP router", zap.Error(err))
		return nil, err
	}

	prometheusServer := &http.Server{
		Addr:        ":9000",
		Handler:     promhttp.Handler(),
		ReadTimeout: 15 * time.Second,
	}

	return &Servers{
		HTTP: &http.Server{
			Addr:    cfg.Server.Port,
			Handler: ginEngine,
		},
		Prometheus: prometheusServer,
		Redis:      redisConn,
		DB:         pool,
	}, nil
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}
	return zap.NewAtomicLevelAt(level)
}

func initDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Info("database ping failed:", zap.Error(err))
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func newRedisClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Проверяем подключение
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

func createTokenUtils(redisConn *redis.Client, secretKey string, expiration time.Duration) utils.TokenUtils {
	return utils.NewTokenService(blackList.NewBlackList(redisConn), secretKey, expiration)
}

func createAuthService(redisConn *redis.Client, tokenUtils utils.TokenUtils, pool *pgxpool.Pool) service.AuthService {
	return authService.NewAuthService(
		blackList.NewBlackList(redisConn),
		tokenUtils,
		userRepo.NewUserRepository(pool),
	)
}
func createPVZService(pool *pgxpool.Pool) service.PVZService {
	return pvzService.NewPVZService(pvzRepo.NewPVZRepository(pool), receptionRepo.NewReceiptRepository(pool))
}

func createReceiptService(pool *pgxpool.Pool) service.ReceptionService {
	return receptionService.NewReceptionService(receptionRepo.NewReceiptRepository(pool))
}
