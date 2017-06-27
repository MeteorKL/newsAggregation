package userAPI

import (
	"net/http"

	"github.com/MeteorKL/koala"
)

func UserHandlers() {
	// http://localhost:1123/api/register?mail=2&nickname=2&password=2
	koala.Get("/api/register", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var mail, nickname, password string
		var message, data string
		if arr, ok := p.ParamGet["password"]; ok {
			password = arr[0]
		} else {
			message = "fail"
			data = "请输入密码"
		}
		if arr, ok := p.ParamGet["nickname"]; ok {
			nickname = arr[0]
		} else {
			message = "fail"
			data = "请输入昵称"
		}
		if arr, ok := p.ParamGet["mail"]; ok {
			mail = arr[0]
		} else {
			message = "fail"
			data = "请输入邮箱"
		}
		if message != "fail" {
			if register(mail, nickname, password) {
				message = "success"
				data = "注册成功"
			} else {
				message = "fail"
				data = "邮箱或用户名已经存在"
			}
		}
		koala.WriteJSON(w, map[string]interface{}{
			"message": message,
			"data":    data,
		})
	})

	// http://localhost:1123/api/login?nickname=2&password=2
	koala.Get("/api/login", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var nickname, password string
		var message, data string
		if koala.ExistSession(r, sessionID) {
			message = "fail"
			data = "你已经登录了啊"
		} else {
			if arr, ok := p.ParamGet["password"]; ok {
				password = arr[0]
			} else {
				message = "fail"
				data = "请输入密码"
			}
			if arr, ok := p.ParamGet["nickname"]; ok {
				nickname = arr[0]
			} else {
				message = "fail"
				data = "请输入昵称"
			}
			if message != "fail" {
				if loginCheck(nickname, password) {
					session := koala.GetSession(r, w, sessionID)
					session.Values["nickname"] = nickname
					session.Values["password"] = password
					message = "success"
					data = "登陆成功"
				} else {
					message = "fail"
					data = "密码错误"
				}
			}
		}
		koala.WriteJSON(w, map[string]interface{}{
			"message": message,
			"data":    data,
		})
	})

	// http://localhost:1123/api/logout
	koala.Get("/api/logout", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var message, data string
		message = "fail"
		data = "你还没有登录"
		if session := koala.PeekSession(r, sessionID); session != nil {
			session.Destory()
			message = "success"
			data = "注销成功"
		}
		koala.WriteJSON(w, map[string]interface{}{
			"message": message,
			"data":    data,
		})
	})

	// http://localhost:1123/api/updateTag?tags=国内&tags=军事
	koala.Get("/api/updateTag", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var tags []string
		var nickname, password string
		var message, data string
		if arr, ok := p.ParamGet["tags"]; ok {
			tags = arr
		} else {
			message = "fail"
			data = "请选择你关注的新闻"
		}
		if message != "fail" {
			if session := koala.PeekSession(r, sessionID); session != nil {
				if s, ok := session.Values["nickname"].(string); ok {
					nickname = s
				} else {
					message = "fail"
					data = "session 格式错误"
				}
				if s, ok := session.Values["password"].(string); ok {
					password = s
				} else {
					message = "fail"
					data = "session 格式错误"
				}
				if updateTag(nickname, password, tags) {
					message = "success"
					data = "修改关注成功"
				} else {
					message = "fail"
					data = "修改关注失败"
				}
			} else {
				message = "fail"
				data = "你还没有登录"
			}
		}
		koala.WriteJSON(w, map[string]interface{}{
			"message": message,
			"data":    data,
		})
	})
}
