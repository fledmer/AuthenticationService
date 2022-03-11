package sessions

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func HandlerPostSession(ctx echo.Context) error {
	//Берем из тела логин
	bodyData := make(map[string]string)
	if err := json.NewDecoder(ctx.Request().Body).Decode(&bodyData); err != nil {
		log.Println("postSession can't decode request", err)
		return ctx.JSON(http.StatusInternalServerError, "Message: can't decode body")
	}
	login, found := bodyData["login"]
	if found != true {
		return ctx.JSON(http.StatusBadRequest, "Message: bad body")
	}
	//Берем старые cookie, чтобы удалить их из бд
	if oldCookie, err := ctx.Cookie("sessionToken"); err == nil {
		DeleteSession(oldCookie.Value)
	}
	//Создаем новую сессию
	if token, err := CreateSessions(login); err == nil {
		cookie := new(http.Cookie)
		cookie.Name = "sessionToken"
		cookie.Value = token
		cookie.Expires = time.Now().Add(time.Hour)
		ctx.SetCookie(cookie)
		return ctx.JSON(http.StatusOK, "Login: "+login)
	} else {
		return ctx.JSON(http.StatusInternalServerError, "Message: can't create session")
	}
}

func HandlerGetSession(ctx echo.Context) error {
	//Получаем куки токен
	cookie, err := ctx.Cookie("sessionToken")
	if err != nil {
		log.Println("Get Auth cookie error: ", err)
		return ctx.JSON(http.StatusOK, "Message: nolog")
	}

	login, err := GetSession(cookie.Value)
	if err != nil {
		return ctx.JSON(http.StatusOK, "Message: notlog")
	} else {
		return ctx.JSON(http.StatusOK, "Login: "+login)
	}
}

func getSessionId() string {
	id := ""
	s := rand.Int31()%15 + 15
	for x := int32(0); x < s; x++ {
		id += string(rand.Int31()%25 + 65)
	}
	return id
}
