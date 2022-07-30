package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"social-network/internal/app/middleware"
	"social-network/internal/controller"
	"social-network/internal/repo"
	"social-network/internal/service"
	"social-network/pkg/logger"
	"social-network/pkg/sqlite/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db        sqlite.DB
	topRouter *http.ServeMux
}

func (a *App) Run(port int, path string) error {
	storage, err := a.db.Connect(path)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	a.initRoutes(storage)

	logger.InfoLogger.Println("Starting the application at port:", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), middleware.Cors(middleware.NewResponseHeader(a.topRouter, "Content-Type", "application/json")))
}

func (a *App) initRoutes(storage *sql.DB) {
	a.topRouter = http.NewServeMux()

	WsController := controller.WsController{}

	notificationService := service.NewNotificationService(repo.NewNotificationRepo(storage), &WsController)
	userService := service.NewUserService(repo.NewUserRepo(storage))
	followerService := service.NewFollowerService(repo.NewFRepo(storage), repo.NewUserRepo(storage))
	postService := service.NewPostService(repo.NewPostRepo(storage, *repo.NewUserRepo(storage)), repo.NewUserRepo(storage))
	chatService := service.NewChatService(repo.NewChatRepo(storage), repo.NewGroupRepo(storage))
	groupService := service.NewGroupService(repo.NewGroupRepo(storage), repo.NewFRepo(storage))

	a.topRouter.Handle("/user/", http.StripPrefix("/user", controller.UsersHandler(userService, notificationService, chatService)))
	a.topRouter.Handle("/follower/", http.StripPrefix("/follower", controller.FollowerHandler(followerService, notificationService)))
	a.topRouter.Handle("/post/", http.StripPrefix("/post", controller.PostHandler(postService, notificationService)))
	a.topRouter.Handle("/ws/", http.StripPrefix("/ws", controller.WsHandler(notificationService, chatService)))
	a.topRouter.Handle("/chat/", http.StripPrefix("/chat", controller.ChatHandler(chatService)))
	a.topRouter.Handle("/group/", http.StripPrefix("/group", controller.GroupHandler(groupService, notificationService)))
}
