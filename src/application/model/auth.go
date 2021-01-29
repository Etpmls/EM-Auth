package model

import (
	"errors"
	"github.com/Etpmls/EM-Auth/src/application"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type Auth struct {

}

func (this *Auth) VerifyBasicRoute(w http.ResponseWriter, httpUri string) bool {
	for _, v := range application.NoAuthRoute {
		b, _ := filepath.Match(v, httpUri)
		if b {
			return true
		}
	}
	return false
}

func (this *Auth) NoVerify(w http.ResponseWriter) {
	em.Micro.Response.Http_Success(w, http.StatusOK, em.SUCCESS_Code, "No Auth", nil)
	return
}

func (this *Auth) BasicVerify(w http.ResponseWriter, token string, lang string) {
	if len(token) == 0 {
		em.Micro.Response.Http_Error(w, http.StatusUnauthorized, em.ERROR_Code, em.I18n.TranslateString("ERROR_MESSAGE_PermissionDenied", lang), errors.New("No token"))
		return
	}

	b, err := em.Micro.Auth.VerifyToken(token)
	if err != nil {
		em.Micro.Response.Http_Error(w, http.StatusInternalServerError, em.ERROR_Code, em.I18n.TranslateString("ERROR_MESSAGE_TokenVerificationFailed", lang), err)
		return
	}
	if !b {
		em.Micro.Response.Http_Error(w, http.StatusUnauthorized, em.ERROR_Code, em.I18n.TranslateString("ERROR_MESSAGE_PermissionDenied", lang), err)
		return
	}

	em.Micro.Response.Http_Success(w, http.StatusOK, em.SUCCESS_Code, em.I18n.TranslateString("SUCCESS_Validate", lang), nil)
	return
}

func (this *Auth) AdvancedVerify(w http.ResponseWriter, token string, lang string, uri string, method string) {
	// Get Claims
	// 获取Claims
	tmp, err := em.Micro.Auth.ParseToken(token)
	tk, ok := tmp.(*jwt.Token)
	if err != nil {
		em.Micro.Response.Http_Error(w, http.StatusInternalServerError, em.ERROR_Code, em.I18n.TranslateString("ERROR_MESSAGE_TokenVerificationFailed", lang), err)
		return
	}
	if !ok {
		em.Micro.Response.Http_Error(w, http.StatusBadRequest, em.ERROR_Code, em.I18n.TranslateString("ERROR_MESSAGE_TokenVerificationFailed", lang), err)
		return
	}

	// Determine whether the role has the corresponding permissions
	// 判断所属角色是否有相应的权限
	if claims,ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		if userId, ok := claims["jti"].(string); ok {
			b, _ := this.PermissionVerify(uri, method, userId)
			if b {
				em.Micro.Response.Http_Success(w, http.StatusOK, em.SUCCESS_Code, em.I18n.TranslateString("SUCCESS_Validate", lang), nil)
				return
			}
		}
	}

	em.Micro.Response.Http_Error(w, http.StatusUnauthorized, em.ERROR_Code, em.I18n.TranslateString("ERROR_MESSAGE_PermissionDenied", lang), err)
	return
}

func (this *Auth) PermissionVerify(httpUri string, httpMethod string, idStr string) (b bool, err error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return b, err
	}

	// 1.获取用户ID
	var u User
	em.DB.Preload("Roles").First(&u, id)
	var ids []uint
	for _, v := range u.Roles {
		// 如果为管理员组
		if v.ID == 1 {
			return true, nil
		}
		ids = append(ids, v.ID)
	}
	// 获取角色相关权限
	var r []Role
	em.DB.Preload("Permissions").Where(ids).Find(&r)

	// 获取当前URL Path
	tmpUri, err := url.Parse(httpUri)
	if err != nil {
		return b, err
	}
	uri := tmpUri.Path

	// Determine whether there is a request permission
	// 判断是否有请求权限
	for _, v := range r {
		for _, subv := range v.Permissions {

			// define an empty slice
			// 定义一个空切片
			var mtd = []string{}
			mtd = strings.Split(subv.Method, ",")

			// Path comparison
			// 路径对比
			b, _ := filepath.Match(subv.Path, uri)
			if b {

				// Method comparison
				// 方法对比
				for _, mtdv := range mtd {
					// If it is ALL, return the permission verification success directly
					// 如果是ALL直接返回权限验证成功
					if mtdv == "ALL" {
						return true, nil
					}
					// If the method is the same as the current request, return the verification success
					// 如果与当前请求方法相同，返回验证成功
					if mtdv == httpMethod {
						return true, nil
					}
				}

			}
		}
	}

	return false, err
}