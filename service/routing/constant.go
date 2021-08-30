// Package routing
// Created by GoLand
// @User: lenora
// @Date: 2021/7/27
// @Time: 16:02

package routing

var RouterColumns = struct {
	ID            string
	CreatedAt     string
	UpdatedAt     string
	NavID         string
	AppID         string
	Title         string
	ParentID      string
	Key           string
	FormID        string
	Icon          string
	Status        string
	OrderNum      string
	Action        string
	ActionContent string
}{
	ID:            "id",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	NavID:         "nav_id",
	AppID:         "app_id",
	Title:         "title",
	ParentID:      "parent_id",
	Key:           "key",
	FormID:        "form_id",
	Icon:          "icon",
	Status:        "status",
	OrderNum:      "order_num",
	Action:        "action",
	ActionContent: "action_content",
}

var TRUE, FALSE = true, false

type MAPS map[string]uint8

const (
	MENU uint8 = iota + 1
	OPEN
	LINK
)

var actionType = MAPS{
	"menu": MENU,
	"open": OPEN,
	"link": LINK,
}

func (m MAPS) searchKey(key string) uint8 {
	for k, v := range m {
		if key == k {
			return v
		}
	}
	return 0
}
