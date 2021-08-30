//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 8:01 下午
package user

import "github.com/Lenora-Z/low-code/utils"

func (srv *userService) CreateUser(account, nickName, trueName, mobile, mail string, appId, groupId uint64) (error, *User) {
	pwd, err := utils.CryptoPassword(DEFAULT_PWD)
	if err != nil {
		return err, nil
	}
	return createUser(srv.db, &User{
		AppID:     appId,
		Account:   account,
		Password:  pwd,
		Nickname:  nickName,
		TrueName:  trueName,
		Mobile:    mobile,
		Mail:      mail,
		GroupID:   groupId,
		IsDelete:  &FALSE,
		PwdStatus: &FALSE,
	})
}

func (srv *userService) GetMultiUser(ids []uint64) (error, UserList) {
	return getMultiUserDetail(srv.db, ids)
}

func (srv *userService) GetUserListByGroupId(groupId, appId uint64) (error, uint32, UserList) {
	group := []uint64{groupId}
	return getUserList(srv.db, 0, utils.MAX_LIMIT, appId, "", group, "")
}

func (srv *userService) GetUserList(page, limit uint32, appId uint64, name string, group []uint64) (error, uint32, UserList) {
	offset := (page - 1) * limit
	return getUserList(srv.db, offset, limit, appId, name, group, "")
}

func (srv *userService) GetAllUser(appId uint64, name string, group ...uint64) (error, uint32, UserList) {
	return getUserList(srv.db, 0, utils.MAX_LIMIT, appId, "", group, name)
}

func (srv *userService) UpdateUser(id uint64, account, nickName, trueName, mobile, mail string, groupId uint64) (error, *User) {
	return updateUser(srv.db, &User{
		ID:       id,
		Account:  account,
		Nickname: nickName,
		TrueName: trueName,
		Mobile:   mobile,
		Mail:     mail,
		GroupID:  groupId,
	})
}

func (srv *userService) ResetUserPwd(id uint64, pwd ...string) (error, *User) {
	var password = DEFAULT_PWD
	if len(pwd) > 0 {
		password = pwd[0]
	}
	password, err := utils.CryptoPassword(password)
	if err != nil {
		return err, nil
	}
	return updateUser(srv.db, &User{
		ID:        id,
		Password:  password,
		PwdStatus: &TRUE,
	})
}

func (srv *userService) UpdatePwdStatus(id uint64, status *bool) (error, *User) {
	return updateUser(srv.db, &User{
		ID:        id,
		PwdStatus: status,
	})
}

func (srv *userService) DeleteUser(id uint64) error {
	err, _ := updateUser(srv.db, &User{
		ID:       id,
		IsDelete: &TRUE,
	})
	return err
}

func (srv *userService) GetUserByAccount(name string, appId uint64) (bool, *User) {
	column := make(map[string]interface{}, 0)
	column[UserColumns.Account] = name
	column[UserColumns.IsDelete] = 0
	column[UserColumns.AppID] = appId
	return userDetailByColumn(srv.db, column)
}

func (srv *userService) GetUserByTrueName(name string, appId uint64) (bool, *User) {
	column := make(map[string]interface{}, 0)
	column[UserColumns.TrueName] = name
	column[UserColumns.IsDelete] = 0
	column[UserColumns.AppID] = appId
	return userDetailByColumn(srv.db, column)
}

func (srv *userService) GetUser(id uint64) (bool, *User) {
	return userDetail(srv.db, id)
}
