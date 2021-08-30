//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 4:35 下午
package user

const DEFAULT_PWD = "123456"

type PermissionList []PermissionId

type PermissionId struct {
	Id uint64 `json:"id"`
}

var UserColumns = struct {
	ID        string
	AppID     string
	Account   string
	Password  string
	Nickname  string
	TrueName  string
	Mobile    string
	Mail      string
	GroupID   string
	PwdStatus string
	IsDelete  string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	AppID:     "app_id",
	Account:   "account",
	Password:  "password",
	Nickname:  "nickname",
	TrueName:  "true_name",
	Mobile:    "mobile",
	Mail:      "mail",
	GroupID:   "group_id",
	PwdStatus: "pwd_status",
	IsDelete:  "is_delete",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
