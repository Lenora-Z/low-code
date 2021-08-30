package routing

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Router 路由表
type Router struct {
	ID            uint64    `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
	NavID         uint64    `gorm:"index:nav;column:nav_id;type:int(11) unsigned;not null;default:0" json:"nav_id"` // 路由id
	AppID         uint64    `gorm:"column:app_id;type:int(11);not null;default:0" json:"app_id"`                    // app id
	Title         string    `gorm:"column:title;type:varchar(100);not null;default:''" json:"title"`                // 路由名称
	ParentID      uint64    `gorm:"column:parent_id;type:int(11) unsigned;not null;default:0" json:"parent_id"`     // 父级路由
	Key           string    `gorm:"index:route_key;column:key;type:varchar(20);not null;default:0" json:"key"`      // key值
	FormID        uint64    `gorm:"column:form_id;type:int(11) unsigned;not null;default:0" json:"form_id"`         // 对应的表单id
	Icon          string    `gorm:"column:icon;type:varchar(255)" json:"icon"`                                      // 路由icon
	Status        *bool     `gorm:"column:status;type:tinyint(1) unsigned;not null;default:0" json:"status"`        // 生效状态 0-未生效 1-已生效
	OrderNum      uint32    `gorm:"column:order_num;type:int(5) unsigned;not null;default:0" json:"order_num"`      // 排序号
	Action        uint8     `gorm:"column:action;type:tinyint(3) unsigned;not null;default:1" json:"action"`        // 触发事件
	ActionContent string    `gorm:"column:action_content;type:text" json:"action_content"`                          // 事件内容
}

// TableName get sql table name.获取数据库表名
func (m *Router) TableName() string {
	return "router"
}

var tableName = (&Router{}).TableName()

func createRouter(db *gorm.DB, m Router) (error, *Router) {
	err := db.Table(tableName).Create(&m).Error
	return err, &m
}

func updateRouter(db *gorm.DB, m Router) (error, *Router) {
	m.UpdatedAt = time.Now()
	err := db.Table(tableName).Where("id = ? and "+RouterColumns.AppID+" = ?", m.ID, m.AppID).Update(&m).Error
	return err, &m
}

func DeleteWithIds(db *gorm.DB, ids []uint64) error {
	return db.Where("id in (?)", ids).Delete(&Router{}).Error
}

type RouterList []Router

func ListRouterWithAppId(db *gorm.DB, appId uint64) (RouterList, error) {
	RouterList := make(RouterList, 0, 100)
	sql := db.Table(tableName).
		Select("*").Where(RouterColumns.AppID+" = ?", appId)

	err := sql.Order(RouterColumns.OrderNum + " asc").Scan(&RouterList).Error
	return RouterList, err
}

func (uList RouterList) Ids() []uint64 {
	ids := make([]uint64, 0, len(uList))
	for _, v := range uList {
		ids = append(ids, v.ID)
	}
	return ids
}

func GetRouterGroupWithIds(db *gorm.DB, ids []uint64) (error, RouterList) {
	list := RouterList{}
	err := db.Table(tableName).
		Where("status = 1").
		Where("id in (?)", ids).
		Order(RouterColumns.OrderNum + " asc").
		Scan(&list).
		Error
	return err, list
}

func GetRouterList(db *gorm.DB, params map[string]interface{}) RouterList {
	var record Router
	var list RouterList
	sql := db.Table(record.TableName())
	for k, v := range params {
		sql = sql.Where(k+"=?", v)
	}
	sql.Order(RouterColumns.OrderNum + " asc").Scan(&list)
	return list
}
