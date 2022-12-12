package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
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

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Debugln("Loaded configuration")

	conn, err := pgx.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Debugln("Connected to PostgreSQL")

	userRepo := adapter.NewUserRepository(conn)

	userUC := usecase.NewUserUseCase(userRepo)

	userController := controller.NewUserController(userUC)

	app := fiber.New(fiber.Config{
		ErrorHandler: controller.ErrHandler,
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

	log.Debugln("Application has started")

	exit := make(chan os.Signal)

	signal.Notify(exit, os.Interrupt)

	<-exit

	err = app.Shutdown()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Debugln("Application has been shut down")
}
