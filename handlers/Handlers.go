package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"main/systemComponents/deeds"
	"main/systemComponents/sessions"
	"main/systemComponents/users"
	"net/http"
)

func GetUser(ctx echo.Context) error {
	requestID := ctx.Param("ID")
	//Проверка токена
	token, err := ctx.Cookie("sessionToken")
	if err != nil || token.Value == "" {
		return ctx.JSON(http.StatusUnauthorized, "")
	}
	ID, err := sessions.GetSession(token.Value)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, "")
	}
	//Получаем данные о юзере
	userData, err := users.GetUserByID(requestID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "")
	}
	if ID == requestID {
		return ctx.JSON(http.StatusOK, userData)
	} else {
		return ctx.JSON(http.StatusOK,
			"First Name: "+userData.FirstName+
				"Second Name: "+userData.LastName)
	}
}

func AuthenticateUser(ctx echo.Context) error {
	//Берем из тела логин и пароль
	user := users.User{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&user); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, "Message: can't decode body")
	}
	if user.Mail == "" || user.Password == "" {
		return ctx.JSON(http.StatusBadRequest, "Message: no email/pass")
	}
	if _, token, err := user.Authentication(); err != nil {
		log.Println("Registration error: ", err)
		return ctx.JSON(http.StatusBadRequest, "Message: "+err.Error())
	} else {
		ctx.SetCookie(sessions.RecreateCookie(token, ctx))
		return ctx.JSON(http.StatusOK, "Message: "+" Вход в аккаунт выполнен!")
	}
}

func RegistrateUser(ctx echo.Context) error {
	//Берем логин из тела
	user := users.User{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&user); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusInternalServerError, "Message: can't decode body")
	}
	if user.Mail == "" || user.Password == "" {
		return ctx.JSON(http.StatusBadRequest, "Message: no email/pass")
	}
	if token, err := user.Registration(); err != nil {
		log.Println("Registration error: ", err)
		return ctx.JSON(http.StatusBadRequest, "Message: "+err.Error())
	} else {
		ctx.SetCookie(sessions.RecreateCookie(token, ctx))
		return ctx.JSON(http.StatusOK, "Message: "+"You have been register!")
	}
}

func UpdateUser(ctx echo.Context) error {
	ID := ctx.Param("ID")
	//Берем логин из тела
	newUserData := users.User{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&newUserData); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, "Message: can't decode body")
	}

	token, err := ctx.Cookie("sessionToken")
	if err != nil || token.Value == "" {
		return ctx.JSON(http.StatusUnauthorized, "")
	}

	if tokenId, err := sessions.GetSession(token.Value); err != nil {
		return ctx.JSON(http.StatusUnauthorized, "")
	} else {
		//У текущего запрос отправителя есть права?
		if tokenId == ID {
			if err = newUserData.Update(); err == nil {
				return ctx.JSON(http.StatusOK, "Message: successful update")
			} else {
				return ctx.JSON(http.StatusNotFound, "Message: "+err.Error())
			}
		}
	}
	return ctx.JSON(http.StatusBadRequest, "Message: wtf?")
}

func GetDeedByID(ctx echo.Context) error {
	ID := ctx.Param("ID")

	if deed, err := deeds.GetDeedsByID(ID); err != nil {
		return ctx.JSON(http.StatusNotFound, "")
	} else {
		return ctx.JSON(http.StatusOK, deed)
	}
}

func GetAllDeed(ctx echo.Context) error {
	if deeds, err := deeds.GetAllDeeds(); err != nil {
		return ctx.JSON(http.StatusNotFound, "")
	} else {
		return ctx.JSON(http.StatusOK, deeds)
	}
}

func GetDeedByUser(ctx echo.Context) error {
	ID := ctx.Param("ID")

	if deeds, err := deeds.GetDeedsByUserID(ID); err != nil {
		return ctx.JSON(http.StatusNotFound, "")
	} else {
		return ctx.JSON(http.StatusOK, deeds)
	}

}

func CreateDeed(ctx echo.Context) error {
	newDeed := deeds.Deed{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&newDeed); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, "Message: can't decode body")
	}
	if err := newDeed.Registration(); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusInternalServerError, "Error: "+err.Error())
	} else {
		return ctx.JSON(http.StatusOK, "Message: Deed has been registered")
	}
}
