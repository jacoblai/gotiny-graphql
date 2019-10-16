package resultor

import (
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"reflect"
)

var ErrList = map[string]string{
	"1001": "没有收到提交内容",
	"1002": "防注入限制生效",
	"1003": "参数错误或objectid类型错误",
	"1004": "账号密码错误",
	"1005": "内容错误",
	"1006": "token已过期，请重点登陆",
	"1007": "删除操作失败",
	"1008": "记录已存在",
	"1009": "密码设置长度过短",
	"1010": "此密码在服务器端已发生改变",
	//"1011": "请使用代理商账号",
	"1011": "权限不足",
	"1012": "代理商级别限制生效",
	"1013": "该代理商已有上属",
	"1014": "最多只能添加10个收货地址",
	"1015": "当前账号不是被添加人",
	"1016": "代理商oid不存在",
	"1017": "依赖主键不存在",
	"1018": "用户信息与提交信息不符",
	"1019": "用户账号被禁用",
	"1020": "账号状态不允许操作",
	"1021": "账号没有供货人",
	"1022": "提交内容错误",
	"1023": "非可操作用户",
	"1024": "只能公司账号操作",
	"1025": "公司级Agent账号被删除",
	"1026": "该记录已被使用不允许删除",
	"1027": "支付密码错误",
	"1028": "对接代理商必须上传付款凭证",
	"1029": "时间参数错误",
	"1030": "用户已离开房间",
	"1031": "客服状态无改变",
	"1032": "对不起,当前不在工作时间,请在工作时间咨询(上午9:00-12:00;下午13:00-18:00)",
	"1033": "亲,非常抱歉,当前不在工作时间,请在工作时间咨询。(上午9:00-12:00;下午13:00-18:00)。",
}

func RetChanges(w http.ResponseWriter, changes int64) {
	_, _ = fmt.Fprintf(w, `{"ok":%v,"changes":%v}`, true, changes)
}

func RetOk(w http.ResponseWriter, result interface{}, changes int) {
	resValue := reflect.ValueOf(result)
	if result == nil {
		RetChanges(w, 0)
		return
	}
	if resValue.Kind() == reflect.Ptr {
		resValue = resValue.Elem()
	}
	bytes, _ := ffjson.Marshal(result)
	_, _ = fmt.Fprintf(w, `{"ok":%v,"changes":%v,"data":%v}`, true, changes, string(bytes))
}

func RetErr(w http.ResponseWriter, errmsg interface{}) {
	if _, ok := ErrList[errmsg.(string)]; ok {
		_, _ = fmt.Fprintf(w, `{"ok":%v, "errcode":%v, "errmsg":"%v"}`, false, errmsg, ErrList[errmsg.(string)])
	} else {
		_, _ = fmt.Fprintf(w, `{"ok":%v, "errmsg":"%v"}`, false, errmsg)
	}
}

func RetErrList(w http.ResponseWriter, result interface{}, changes int) {
	resValue := reflect.ValueOf(result)
	if result == nil {
		RetChanges(w, 0)
		return
	}
	if resValue.Kind() == reflect.Ptr {
		resValue = resValue.Elem()
	}
	var res interface{}
	if resValue.Kind() == reflect.Array || resValue.Kind() == reflect.Slice {
		res = result
	} else {
		res = []interface{}{result}
	}
	bytes, _ := ffjson.Marshal(res)
	_, _ = fmt.Fprintf(w, `{"ok":%v,"changes":%v,"data":%v}`, false, changes, string(bytes))
}
