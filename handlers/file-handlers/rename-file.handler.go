package filehandlers

import (
	"github.com/aj-2000/shc-backend/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// why we are using struct here?
// what is the purpose of this struct?
type RenameFileDto struct {
	Name string `json:"name"`
}

func RenameFile(c fiber.Ctx, as *services.AppService) error {

	fileIdString := c.Params("id")

	fileId, err := uuid.Parse(fileIdString)

	if err != nil {
		return err
	}

	userIdString := string(c.Request().Header.Peek("user_id"))

	userId, err := uuid.Parse(userIdString)

	if err != nil {
		return err
	}

	body := new(RenameFileDto)

	if err := c.Bind().Body(body); err != nil {
		return err
	}

	// read it from service ✅
	f, err := as.FileService.RenameFile(userId, fileId, body.Name)

	if err != nil {
		return err
	}

	//why we have to return c.JSON it?
	return c.JSON(fiber.Map{
		"id":         f.ID,
		"r2_path":    f.R2Path,
		"name":       f.Name,
		"size":       f.Size,
		"mime_type":  f.MimeType,
		"extension":  f.Extension,
		"user_id":    f.UserId,
		"updated_at": f.UpdatedAt,
	})
}
