// ws服务端
package controllers

import (
	"encoding/json"
	"io/ioutil"
	"life/consts"
	"life/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/govalidator"
)

type WSController struct {
	BaseController
}

// 接收消息结构体
type receiveMessage struct {
	Action string
	//Token  string
}

// 返回消息结构体
type returnMessage struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Action string `json:"action"`
	//Token  string      `json:"token"`
	Data interface{} `json:"data"`
}

var (
	conn     *websocket.Conn
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (c *WSController) Index(ctx *gin.Context) {
	var err error

	conn, err = upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//log.Error("websocket读取客户端数据错误", err)
		//http.NotFound(ctx.Writer, ctx.Request)
		//return
		goto ERR
	}

	//启动协程
	// go func() {
	// 	//主动向客户端发心跳
	// 	for {
	// 		err = conn.WriteMessage(websocket.TextMessage, []byte("~H#S~"))
	// 		if err != nil {
	// 			return //退出循环，并且代码不会再执行后面的语句
	// 		}
	// 		//心跳每1秒发送1次
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	for {
		//循环读取客户端送来的数据
		_, data, err := conn.ReadMessage()
		if err != nil {
			//log.Error("websocket连接错误", err)
			//return //退出循环，并且代码不会再执行后面的语句
			goto ERR
		}

		//处理接收到消息
		if string(data) == "~H#C~" { //接收到客户端心跳包
			//回复一个心跳包
			conn.WriteMessage(websocket.TextMessage, []byte("~H#S~"))
		} else {
			var revMsg receiveMessage
			json.Unmarshal(data, &revMsg)

			//参数验证
			rules := govalidator.MapData{}
			rules["action"] = []string{"required"}
			rules["token"] = []string{"required"}
			messages := govalidator.MapData{}
			messages["action"] = []string{"required:action 不能为空"}
			messages["token"] = []string{"required:token 不能为空"}
			opts := govalidator.Options{
				Data:            &revMsg,
				Rules:           rules,
				Messages:        messages,
				RequiredDefault: false,
			}
			valid := govalidator.New(opts)
			e := valid.ValidateStruct()
			if len(e) > 0 {
				for _, err := range e {
					conn.WriteMessage(websocket.TextMessage, []byte("invalid request, received->"+err[0]))
				}
			}

			if revMsg.Action == "hardware" { //系统信息
				//cpu使用率
				cpu_usage, _ := utils.GetCpuPercent()
				//cpu温度
				cpu_temp, _ := utils.GetCpuTemp()
				//内存使用率
				ram_usage, _ := utils.GetRamPercent()

				//获取wifi信息
				wifi_status := make(map[string]interface{})
				wifi_status["name"] = "乙太网路"
				wifi_status["value"] = 100

				//组装数据
				resp := make(map[string]interface{}) //创建1个空集合
				resp["cpu_usage"] = cpu_usage
				resp["ram_usage"] = ram_usage
				resp["cpu_temp"] = cpu_temp
				resp["wifi_status"] = wifi_status
				c.Success(revMsg.Action, "success", resp)
			} else if revMsg.Action == "weather" { //天气信息
				token := viper.GetString("caiyun_token")
				url := "https://api.caiyunapp.com/v2.6/" + token + "/121.4159,31.0281/weather?alert=true"
				res, err := http.Get(url)
				if err != nil || res.StatusCode != http.StatusOK {
					c.Error(revMsg.Action, -1, "请求错误", nil)
				} else {
					body, _ := ioutil.ReadAll(res.Body)
					var info interface{}
					json.Unmarshal(body, &info)
					c.Success(revMsg.Action, "success", &info)
				}
			}
		}
	}

ERR:
	log.Error("websocket", err)
	conn.Close()
}

// 成功返回
func (c *WSController) Success(action string, msg string, data interface{}) {
	if msg == "" {
		msg = "success"
	}
	if data == nil || data == "" {
		data = struct{}{}
	}
	res := &returnMessage{
		consts.SUCCESS, msg, action, data, //0=成功
	}
	jsonString, _ := json.Marshal(res)
	conn.WriteMessage(websocket.TextMessage, []byte(jsonString)) //发送消息
}

// 失败返回
func (c *WSController) Error(action string, code int, msg string, data interface{}) {
	if code >= 0 {
		code = consts.UNKNOWN_ERROR_STATUS
	}
	if msg == "" {
		msg = "error"
	}
	if data == nil || data == "" {
		data = struct{}{}
	}
	res := &returnMessage{
		code, msg, action, data, //0=成功
	}
	jsonString, _ := json.Marshal(res)
	conn.WriteMessage(websocket.TextMessage, []byte(jsonString)) //发送消息
}
