package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"news-app-api/config"
	"news-app-api/internal/adapter"
	"news-app-api/internal/controller"
	"news-app-api/internal/usecase"
	"os"
	"os/signal"
)

func Run() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Debug("Loaded configuration")

	conn, err := pgx.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Debug("Connected to PostgreSQL")

	userRepo := adapter.NewUserRepository(conn)
	mediaRepo := adapter.NewMediaRepository(conn)
	audioFileRepo, err := adapter.NewAudioFileRepository()
	if err != nil {
		log.Fatal(err.Error())
	}
	imageFileRepo, err := adapter.NewImageFileRepository()
	if err != nil {
		log.Fatal(err.Error())
	}
	videoFileRepo, err := adapter.NewVideoFileRepository()
	if err != nil {
		log.Fatal(err.Error())
	}

	userUC := usecase.NewUserUseCase(userRepo)
	mediaUC := usecase.NewMediaUseCase(mediaRepo)
	newsUC := usecase.NewNewsUseCase(
		func() adapter.NewsRepository {
			return adapter.NewNewsRepository(conn)
		},
		mediaRepo,
		audioFileRepo,
		imageFileRepo,
		videoFileRepo,
	)
	feedUC := usecase.NewFeedUseCase(func() adapter.NewsRepository {
		return adapter.NewNewsRepository(conn)
	})

	middleware := controller.NewMiddleware()

	userController := controller.NewUserController(userUC)
	mediaController := controller.NewMediaController(mediaUC)
	newsController := controller.NewNewsController(newsUC)
	feedController := controller.NewFeedController(feedUC)
	favoriteController := controller.NewFavoriteController(newsUC)

	app := fiber.New(fiber.Config{
		ErrorHandler:          controller.ErrHandler,
		DisableStartupMessage: true,
		BodyLimit:             4 * 1024 * 1024 * 1024,
	})

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: cfg.Secret,
	}))

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000, http://127.0.0.1:3000, http://0.0.0.0:3000",
	}))

	app.Use(func(ctx *fiber.Ctx) error {
		err = ctx.Next()
		if err != nil {
			err = ctx.App().ErrorHandler(ctx, err)
			if err != nil {
				return err
			}
		}

		log.WithField(
			"status", ctx.Response().StatusCode(),
		).WithField(
			"method", ctx.Method(),
		).WithField(
			"path", ctx.Path(),
		).Info("Request")
		return nil
	})

	router := app.Group("api")

	userRouter := router.Group("users")
	mediaRouter := router.Group("media")
	newsRouter := router.Group("news")
	feedRouter := router.Group("feed")
	favoriteRouter := router.Group("favorites")

	userController.RegisterRoutes(userRouter, middleware)
	mediaController.RegisterRoutes(mediaRouter, middleware)
	newsController.RegisterRoutes(newsRouter, middleware)
	feedController.RegisterRoutes(feedRouter, middleware)
	favoriteController.RegisterRoutes(favoriteRouter, middleware)

	go func() {
		err = app.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	log.Info("Application has started")

	exit := make(chan os.Signal)

	signal.Notify(exit, os.Interrupt)

	<-exit

	err = app.Shutdown()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Application has been shut down")
}
