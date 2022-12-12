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

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Debug("Loaded configuration")

	conn, err := pgx.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Debug("Connected to PostgreSQL")

	userRepo := adapter.NewUserRepository(conn)

	userUC := usecase.NewUserUseCase(userRepo)

	userController := controller.NewUserController(userUC)

	app := fiber.New(fiber.Config{
		ErrorHandler: controller.ErrHandler,
	})

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: encryptcookie.GenerateKey(),
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
		).Info()
		return nil
	})

	router := app.Group("/api")

	userRouter := router.Group("/users")

	userController.RegisterRoutes(userRouter)

	go func() {
		err = app.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	log.Debug("Application has started")

	exit := make(chan os.Signal)

	signal.Notify(exit, os.Interrupt)

	<-exit

	err = app.Shutdown()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Debug("Application has been shut down")
}
