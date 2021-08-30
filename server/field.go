package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lenora-Z/low-code/service/form"
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

func (ds *defaultServer) getMembers(value interface{}, srv user.UserService) (error, []string) {
	if srv == nil {
		srv = user.NewUserService(ds.db)
	}
	userIdInterface, ok := value.(bson.A)
	if !ok {
		return errors.New(""), []string{}
	}
	userIds := make([]uint64, 0, cap(userIdInterface))
	for _, x := range userIdInterface {
		id, status := x.(float64)
		if !status {
			continue
		}
		userIds = append(userIds, uint64(id))

	}
	err, u := srv.GetMultiUser(userIds)
	if err != nil {
		return err, nil
	}
	return nil, u.TrueNames()
}

func (ds *defaultServer) getOrganizations(value interface{}, srv user.UserService) (error, []string) {
	if srv == nil {
		srv = user.NewUserService(ds.db)
	}
	orgs, ok := value.(bson.A)
	if !ok {
		return errors.New(""), []string{}
	}
	org := make([]uint64, 0, cap(orgs))
	for _, o := range orgs {
		id, status := o.(float64)
		if !status {
			continue
		}
		org = append(org, uint64(id))
	}
	err, organization := srv.GetOrganizationByIdGroup(org)
	if err != nil {
		return err, nil
	}
	return nil, organization.Names()
}

func (ds *defaultServer) getTime(value interface{}) *time.Time {
	d, ok := value.(float64)
	if !ok {
		return nil
	}
	t := utils.Time(d)
	return &t
}

type fileList []fileFormat

func (ds *defaultServer) getFile(value interface{}) *fileList {
	byteData, _ := json.Marshal(value)
	fs := make(fileList, 0)
	if err := json.Unmarshal(byteData, &fs); err != nil {
		return nil
	}
	return &fs
}

func (list fileList) path() []string {
	ret := make([]string, 0, cap(list))
	for _, x := range list {
		ret = append(ret, x.Path)
	}
	return ret
}

func (ds *defaultServer) getArea(value interface{}) []string {
	v, ok := value.(bson.A)
	if !ok {
		return []string{}
	}
	areas := make([]string, 0, cap(v))
	for _, x := range v {
		area, status := x.(string)
		if !status {
			continue
		}
		areas = append(areas, area)
	}
	return areas
}

func (ds *defaultServer) getItemData(types string, userSrv user.UserService, value interface{}) interface{} {
	var ret interface{}
	switch types {
	case "created_at":
		t, ok := value.(primitive.DateTime)
		if ok {
			ret = t.Time().Format(utils.TimeFormatStr)
		} else {
			ret = ds.getTime(value).Format(utils.TimeFormatStr)
		}
	case form.MEMBER:
		e, names := ds.getMembers(value, userSrv)
		if e != nil && e.Error() != "" {
			logrus.Info("get user error:", e.Error())
			return ""
		} else {
			ret = strings.Join(names, ",")
		}
	case form.ORGANIZATION:
		e, names := ds.getOrganizations(value, userSrv)
		if e != nil && e.Error() != "" {
			logrus.Info("get organization error:", e.Error())
			return ""
		} else {
			ret = strings.Join(names, ",")
		}
	case form.DATETIME:
		if value == nil {
			return ""
		} else {
			ret = ds.getTime(value).Format(utils.TimeFormatStr)
		}
	case form.DATETIME_RANGE:
		dr, ok := value.(bson.A)
		if !ok {
			return ""
		}
		ranges := make([]string, 0, cap(dr))
		for _, d := range dr {
			ranges = append(ranges, ds.getTime(d).Format(utils.TimeFormatStr))
		}
		ret = fmt.Sprintf("%s~%s", ranges[0], ranges[1])
	case form.FILE:
		files := ds.getFile(value)
		if files == nil {
			return ""
		}
		ret = strings.Join(files.path(), ",")
	case form.AUTOGRAPH:
		files := ds.getFile(value)
		if files == nil {
			return ""
		}
		ret = strings.Join(files.path(), ",")
	case form.AREA:
		area := ds.getArea(value)
		ret = strings.Join(area, ",")
	default:
		ret = value
	}
	return ret
}
