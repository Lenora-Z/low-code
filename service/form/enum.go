package form

type MAPS map[string]uint8

const (
	SINGLE uint8 = iota + 1
	MULTI
	TABLE
)

var multiFormShowType = MAPS{
	"single": SINGLE,
	"multi":  MULTI,
	"table":  TABLE,
}

const (
	RELATION uint8 = iota + 1
	COPY
)

var multiFormImportType = MAPS{
	"relation": RELATION,
	"copy":     COPY,
}

const (
	CARD uint8 = iota + 1
	BACK_FILL
)

var recordsModeType = MAPS{
	"card":     CARD,
	"backfill": BACK_FILL,
}

const (
	PROCESS      uint8 = iota + 1 //触发流程
	CLOSE                         //关闭页面
	SUBMIT_CLOSE                  //提交数据并关闭页面(暂未使用)
	SUBMIT                        //提交数据
	API                           //触发接口
	OPEN                          //跳转弹窗
	LINK                          //打开标准页
	DETAIL                        //查看详情
	UPDATE                        //更新数据
	REMOVE                        //删除
)

var ButtonType = MAPS{
	"process": PROCESS,
	"open":    OPEN,
	"close":   CLOSE,
	"submit":  SUBMIT,
	"link":    LINK,
	"api":     API,
	"detail":  DETAIL,
	"update":  UPDATE,
	"remove":  REMOVE,
}

const (
	DEFAULT uint8 = iota + 1
	CUSTOM
	FIRST_LEVEL
	SECOND_LEVEL
	THIRD_LEVEL
	SECTION
)

var TextSize = MAPS{
	"default":     DEFAULT,
	"custom":      CUSTOM,
	"firstLevel":  FIRST_LEVEL,
	"secondLevel": SECOND_LEVEL,
	"thirdLevel":  THIRD_LEVEL,
	"section":     SECTION,
}

func (m MAPS) searchKey(key string) uint8 {
	for k, v := range m {
		if key == k {
			return v
		}
	}
	return 0
}
