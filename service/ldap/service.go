// Package ldap
// Created by GoLand
// @User: lenora
// @Date: 2021/8/26
// @Time: 10:16

package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/sirupsen/logrus"
	"time"
)

func Login(name, password string) {
	filter := fmt.Sprintf("(%s=%s)", USER_NAME_KEY, name)
	attr := []string{
		USER_NAME_KEY,
		EMAIL_KEY,
		NAME_KEY,
		"objectClass",
		"description",
	}

	conn, err := ldap.Dial("tcp", HOST+":"+PORT)
	if err != nil {
		logrus.Error(err)
		return
	}
	conn.SetTimeout(5 * time.Second)
	defer conn.Close()

	sql := ldap.NewSearchRequest(
		BASE_ON,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		attr,
		nil,
	)

	cur, err := conn.Search(sql)
	if err != nil {
		logrus.Error(err)
		return
	}

	if len(cur.Entries) == 0 {
		logrus.Error("no user")
		return
	}

	entry := cur.Entries[0]

	logrus.Println(entry.GetAttributeValues(USER_NAME_KEY))

	err = conn.Bind(entry.DN, password)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("success:", entry.Attributes)
}
