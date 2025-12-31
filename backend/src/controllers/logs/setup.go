package logsC

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/ranon-rat/self-hosting-manager/src/domain"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
	"github.com/ranon-rat/self-hosting-manager/src/domain/repositories"
	"github.com/ranon-rat/self-hosting-manager/src/middleware"
)

var connectionsTunnels = domain.NewSecureMap[int, *domain.SecureMap[*websocket.Conn, bool]]()
var logRepo executionlogs.ExecutionLogsRepoDB
var pRepo projectsD.ProjectsRepoDB

func Setup(app *fiber.App, repos repositories.Repositories) {
	logRepo = repos.LogRepo
	pRepo = repos.ProjectRepo

	group := app.Group("/logs")
	group.Get("/", GetLogs)
	group.Get("/ws", middleware.ConnUpgradable(), websocket.New(WebsocketConn))
}
