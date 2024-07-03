package handlers

import (
	"strconv"

	m "github.com/aj-2000/shc-backend/models"
	"github.com/aj-2000/shc-backend/services"
	"github.com/gofiber/fiber/v3"
)

// why we made this struct?  ✅
type CheckOtpRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

func VerifyOtpAndGetTokens(c fiber.Ctx, as *services.AppService) error {
	req := new(CheckOtpRequest)

	//what is the meaning of binding body? ✅
	if err := c.Bind().Body(req); err != nil {
		return err
	}

	otp, err := strconv.Atoi(req.Otp)

	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid OTP"}
	}

	// read it from the service ✅
	if err = as.AuthService.VerifyOtp(req.Email, otp); err != nil {
		return &fiber.Error{Code: fiber.StatusUnauthorized, Message: err.Error()}
	}

	// does this mean we have created a pointer to the user? ✅
	var user *m.User

	// read it from the service ✅
	u, err := as.UserService.FindUserByEmail(req.Email)

	// why we are doing this, does it mean we are verifying if user exists or not?
	//if it is so then why we are not returning error if user does not exist? why we are making a new user?? ✅
	if err != nil {
		user, err = as.UserService.CreateUser(&m.User{Email: req.Email, Name: req.Name})
	} else {
		// what is the meaning of update user? what are we exaclty updating in user? read it from service ✅
		user, err = as.UserService.UpdateAUser(&m.User{ID: u.ID, Email: req.Email, Name: req.Name})
	}

	if err != nil {
		return err
	}

	//read it from the service ✅
	tokens, err := as.AuthService.GenerateTokens(user.ID, user.Name, user.Email)

	if err != nil {
		return err
	}

	// what is fiber.Map? and what does c.JSON do? why we have to return it ? and where this c.JSON is returning? ✅
	return c.JSON(fiber.Map{
		"refresh_token": tokens.RefreshToken,
		"access_token":  tokens.AccessToken,
		"name":          user.Name,
		"email":         user.Email,
		"id":            user.ID,
	})
}
