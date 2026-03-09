package account

import (
	"Monitoring-Pressure/session"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ProcessRequest(ctx *gin.Context) {

	var userSession session.Session

	defer func() {
		if userSession == nil {
			userSession, _ = session.CreateSession()
		}

		ctx.Set(MercurySessionName, userSession)
	}()

	// 从cookie中获取session_id
	cookie, err := ctx.Request.Cookie(CookieSessionId)
	if err != nil {
		// 不存在的话，设置user_id=0，login_status为0，表示未登陆状态
		ctx.Set(MercuryUserId, int64(0))
		ctx.Set(MercuryUserLoginStatus, int64(0))
		return
	}

	sessionId := cookie.Value
	if len(sessionId) == 0 {
		ctx.Set(MercuryUserId, int64(0))
		ctx.Set(MercuryUserLoginStatus, int64(0))
		return
	}

	// 根据sessionId，获取用户的session
	userSession, err = session.Get(sessionId)
	if err != nil {
		ctx.Set(MercuryUserId, int64(0))
		ctx.Set(MercuryUserLoginStatus, int64(0))

		return
	}

	tmpUserId, err := userSession.Get(MercuryUserId)
	if err != nil {
		ctx.Set(MercuryUserId, int64(0))
		ctx.Set(MercuryUserLoginStatus, int64(0))
		return
	}

	userId, ok := tmpUserId.(int64)
	if !ok || userId == 0 {
		ctx.Set(MercuryUserId, int64(0))
		ctx.Set(MercuryUserLoginStatus, int64(0))
		return
	}
	ctx.Set(MercuryUserId, int64(userId))
	ctx.Set(MercuryUserLoginStatus, int64(1))
	return
}

func GetUserId(ctx *gin.Context) (userId int64, err error) {

	tmpUserId, exists := ctx.Get(MercuryUserId)
	if !exists {
		err = errors.New("user id not exists")
		return
	}
	userId, ok := tmpUserId.(int64)
	if !ok || userId == 0 {
		err = errors.New("user id covert to int64 failed")
		return
	}

	return
}

func IsLogin(ctx *gin.Context) (login bool) {

	tmpUserStuatus, exists := ctx.Get(MercuryUserLoginStatus)
	if !exists {
		return
	}
	loginStatus, ok := tmpUserStuatus.(int64)
	if !ok {
		return
	}

	if loginStatus == 0 {
		return
	}

	login = true
	return
}

func SetUserId(userId int64, ctx *gin.Context) {

	var userSession session.Session
	tempSession, exists := ctx.Get(MercurySessionName)
	if !exists {
		return
	}

	userSession, ok := tempSession.(session.Session)
	if !ok {
		return
	}

	if userSession == nil {
		return
	}
	userSession.Set(MercuryUserId, userId)
}

func ProcessResponse(ctx *gin.Context) {

	var userSession session.Session
	tempSession, exists := ctx.Get(MercurySessionName)
	if !exists {
		return
	}

	userSession, ok := tempSession.(session.Session)
	if !ok {
		return
	}

	if userSession == nil {
		return
	}

	if userSession.IsModify() == false {
		return
	}

	err := userSession.Save()
	if err != nil {
		return
	}

	sessionId := userSession.Id()
	cookie := &http.Cookie{
		Name:     CookieSessionId,
		Value:    sessionId,
		MaxAge:   CookieMaxAge,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(ctx.Writer, cookie)
	return
}
