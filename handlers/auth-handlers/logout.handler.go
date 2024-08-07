package handlers

import (
	"strings"

	"github.com/aj-2000/shc-backend/services"
	"github.com/gofiber/fiber/v3"
)

func Logout(c fiber.Ctx, as *services.AppService) error {
	//why we need refreshToken for loging out what is the need of refresh token for loging out? ✅
	refreshToken := string(c.Request().Header.Peek("Authorization"))
	if refreshToken == "" {
		refreshToken = c.Cookies("__shc_refresh_token")
	}

	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	//NU
	claim, err := as.AuthService.VerifyRefreshToken(refreshToken)

	if err != nil {
		return c.SendStatus(401)
	}
	//NU
	err = as.SessionService.DeleteSession(claim.SessionId)

	if err != nil {
		return err
	}

	return c.SendStatus(200)
}
