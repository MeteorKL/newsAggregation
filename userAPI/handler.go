package userAPI

import (
	"net/http"

	"github.com/MeteorKL/koala"
)

func UserHandlers() {
	// http://localhost:1123/api/register?mail=2&nickname=2&password=2
	koala.Get("/api/register", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var status int
		var message string
		var data interface{}
		var mail, nickname, password string
		if arr, ok := p.ParamGet["password"]; ok && arr[0] != "" {
			password = arr[0]
		} else {
			status = 1
			message = "请输入密码"
		}
		if arr, ok := p.ParamGet["nickname"]; ok && arr[0] != "" {
			nickname = arr[0]
		} else {
			status = 2
			message = "请输入昵称"
		}
		if arr, ok := p.ParamGet["mail"]; ok && arr[0] != "" {
			mail = arr[0]
		} else {
			status = 3
			message = "请输入邮箱"
		}
		if message != "fail" {
			if register(mail, nickname, password) {
				status = 0
				message = "注册成功"
			} else {
				status = 4
				message = "邮箱或昵称已经存在"
			}
		}
		koala.WriteJSON(w, map[string]interface{}{
			"status":  status,
			"message": message,
			"data":    data,
		})
	})

	koala.Get("/api/checkLogin", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if session := koala.PeekSession(r, sessionID); session != nil {
			koala.WriteJSON(w, map[string]interface{}{
				"status":  0,
				"message": "success",
				"data": map[string]interface{}{
					"nickname": session.Values["nickname"],
					"tags":     session.Values["tags"],
				},
			})
		} else {
			koala.WriteJSON(w, map[string]interface{}{
				"status":  1,
				"message": "fail",
				"data":    nil,
			})
		}
	})
	// http://localhost:1123/api/login?nickname=2&password=2
	koala.Get("/api/login", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var status int
		var message string
		var data interface{}
		var nickname, password string
		if koala.ExistSession(r, sessionID) {
			status = 1
			message = "你已经登录了啊"
		} else {
			if arr, ok := p.ParamGet["password"]; ok && arr[0] != "" {
				password = arr[0]
			} else {
				status = 2
				message = "请输入密码"
			}
			if arr, ok := p.ParamGet["nickname"]; ok && arr[0] != "" {
				nickname = arr[0]
			} else {
				status = 3
				message = "请输入昵称"
			}
			if message != "fail" {
				if user := loginCheck(nickname, password); user != nil {
					session := koala.GetSession(r, w, sessionID)
					session.Values["nickname"] = nickname
					session.Values["password"] = password
					session.Values["tags"] = user["tags"]
					status = 0
					message = "登陆成功"
					data = map[string]interface{}{
						"nickname": nickname,
						"tags":     session.Values["tags"],
					}
				} else {
					status = 4
					message = "密码错误"
				}
			}
		}
		koala.WriteJSON(w, map[string]interface{}{
			"status":  status,
			"message": message,
			"data":    data,
		})
	})

	// http://localhost:1123/api/logout
	koala.Get("/api/logout", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if session := koala.PeekSession(r, sessionID); session != nil {
			session.Destory()
			koala.WriteJSON(w, map[string]interface{}{
				"status":  0,
				"message": "注销成功",
				"data":    nil,
			})
		} else {
			koala.WriteJSON(w, map[string]interface{}{
				"status":  1,
				"message": "你还没有登录",
				"data":    nil,
			})
		}
	})

	// http://localhost:1123/api/updateTags?tags=国内&tags=军事
	koala.Get("/api/updateTags", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var status int
		var message string
		var data interface{}
		var tags []string
		var nickname, password string
		if arr, ok := p.ParamGet["tags"]; ok {
			tags = arr
		}
		if session := koala.PeekSession(r, sessionID); session != nil {
			if s, ok := session.Values["nickname"].(string); ok {
				nickname = s
			} else {
				status = 1
				message = "session 格式错误"
			}
			if s, ok := session.Values["password"].(string); ok {
				password = s
			} else {
				status = 2
				message = "session 格式错误"
			}
			if updateTags(nickname, password, tags) {
				session.Values["tags"] = tags
				status = 0
				message = "修改关注成功"
				data = tags
			} else {
				status = 3
				message = "修改关注失败"
			}
		} else {
			status = 4
			message = "你还没有登录"
		}
		koala.WriteJSON(w, map[string]interface{}{
			"status":  status,
			"message": message,
			"data":    data,
		})
	})
}
