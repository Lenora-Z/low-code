//Created by Goland
//@User: lenora
//@Date: 2021/2/20
//@Time: 2:34 下午
package tritium

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestTritiumService_GetTriCheckModule(t *testing.T) {
	srv := NewTritiumService("192.168.3.22:7685")
	err, arg := srv.GetTriCheckModule("ims")
	logrus.Info(err, arg)
}

func TestTritiumService_SendNotice(t *testing.T) {
	srv := NewTritiumService("192.168.3.22:7685")
	srv.SendNotice(&NoticeArg{})
}
