package service

import (
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/model"
	em "github.com/Etpmls/Etpmls-Micro"
	"net/http"
	"net/url"
	"path/filepath"
)

type ServiceAuth struct {

}

func (this ServiceAuth) Check(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	var (
		token = r.Header.Get("Token")
		lang = r.Header.Get("Language")
		fwd_uri = r.Header.Get("X-Forwarded-Uri")
		method = r.Method
		auth model.Auth
	)

	// Get path without parameters
	// 获取无参数的path
	u, err := url.Parse(fwd_uri)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		em.Micro.Response.Http_Error(w, http.StatusInternalServerError, em.ERROR_Code, em.I18n.TranslateString("ERROR_Validate", lang), err)
		return
	}
	uri := u.Path

	// If it is a basic route, no verification is required
	// 如果是基础路由，无需验证
	b := auth.VerifyBasicRoute(w, uri)
	if b {
		em.Micro.Response.Http_Success(w, http.StatusOK, em.SUCCESS_Code, "No Auth", nil)
		return
	}


	// Get all permission
	var permission model.Permission
	ps, err := permission.GetAll()
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		em.Micro.Response.Http_Error(w, http.StatusInternalServerError, em.ERROR_Code, em.I18n.TranslateString("ERROR_MESSAGE_GetPermission", lang), err)
		return
	}

	// Search URI in permission
	var p model.Permission
	for _, v := range ps {
		b, _ := filepath.Match(v.Path, uri)
		if b {
			p = v
			break
		}
	}

	// If the route is not registered, the token will be verified by default
	// 如果路由没有注册，则默认验证token
	if p.ID == 0 {
		auth.BasicVerify(w, token, lang)
		return
	} else {
		// If the route is registered, verify the request according to the verification method of the route
		// 如果路由已注册，根据路由的验证方法去验证请求
		// No auth
		switch p.Auth {
		case application.Auth_NoVerify:
			auth.NoVerify(w)
			return
		case application.Auth_BasicVerify:
			auth.BasicVerify(w, token, lang)
			return
		case application.Auth_AdvancedVerify:
			auth.AdvancedVerify(w, token, lang, uri, method)
			return
		}
	}
}
