package logsC

import (
	"log"
	"slices"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/ranon-rat/self-hosting-manager/src/domain"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
	executionerServices "github.com/ranon-rat/self-hosting-manager/src/services/executioner"
)

// GET /logs?first_id=0
func GetLogs(c *fiber.Ctx) error {
	firstIDQ := c.QueryInt("id")
	paginated, err := logRepo.Get(firstIDQ)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	oldID := 0

	slices.Reverse(paginated)
	if len(paginated) > 0 {
		oldID = paginated[0].ID
	}
	return c.Status(200).JSON(executionlogs.PaginatedLogs{
		OldestID: oldID,
		HasMore:  len(paginated) >= domain.LIMIT_PAGE,
		Logs:     paginated,
	})
}

// WS /ws/?id-project=0
func WebsocketConn(c *websocket.Conn) {
	idStr := c.Query("id-project")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Close()
		return
	}
	if project, err := pRepo.GetByID(id); err != nil || project.IsPaused {
		c.Close()
		return
	}
	connections, exist := connectionsTunnels.Get(id)
	if !exist {
		connections = domain.NewSecureMap[*websocket.Conn, bool]()
		connectionsTunnels.Set(id, connections)
	}
	connections.Set(c, true)
	defer connections.Delete(c)
	for {
		// aqui me da igual esto
		if _, _, err := c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
	}

}
func CloseAll(id int) {
	connections, exist := connectionsTunnels.Get(id)
	if !exist {
		return
	}
	connections.Range(func(conn *websocket.Conn, _ bool) bool {
		conn.Close()
		return true
	})

	connectionsTunnels.Delete(id)
}
func Messager(id int) {
	channel, exist := executionerServices.OutputChannels.Get(id)
	if !exist {
		CloseAll(id)
		return
	}
	for out := range channel {
		connections, exist := connectionsTunnels.Get(id)
		if !exist {
			break
		}

		connections.Range(func(conn *websocket.Conn, _ bool) bool {
			if err := conn.WriteJSON(fiber.Map{"content": out}); err != nil {
				conn.Close()
				connections.Delete(conn)
			}
			return true
		})
	}
}
