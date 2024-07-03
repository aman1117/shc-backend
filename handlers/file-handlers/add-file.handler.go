package filehandlers

import (
	"context"
	"strings"

	m "github.com/aj-2000/shc-backend/models"
	"github.com/aj-2000/shc-backend/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// what is the need of this struct?
type AddFileDto struct {
	//these tags are called struct tags
	FileName string `json:"file_name"`
	FileSize uint   `json:"file_size"`
	MimeType string `json:"mime_type"`
}

func AddFileToDb(c fiber.Ctx, as *services.AppService) error {
	userIdString := string(c.Request().Header.Peek("user_id"))

	// why do we have to parse it? ✅
	userId, err := uuid.Parse(userIdString)

	if err != nil {
		return err

	}

	//what is AddFileDto? ✅
	body := new(AddFileDto)

	//what we are doing here? what is the meaning of binding body? ✅
	if err := c.Bind().Body(body); err != nil {
		print("error binding body", err.Error())
		return err
	}

	fileSize := body.FileSize

	//read it from service ✅
	subscription, err := as.SubscriptionService.FindSubscriptionByUserId(userId)

	if err != nil {
		return err
	}

	if subscription.TodayRemainingWrites == 0 {
		return &fiber.Error{Code: fiber.StatusPaymentRequired, Message: "You have exceeded your daily write limit"}
	}

	if subscription.StorageRemainingBytes < fileSize {
		return &fiber.Error{Code: fiber.StatusPaymentRequired, Message: "You have exceeded your storage limit"}
	}

	//what is context?
	ctx := context.Background()

	//what is aws.String? -> creating a pointer to the string from the string.
	// how to use strings.Replace? ✅
	key := aws.String(userId.String() + "/" + uuid.NewString() + "_" + strings.Replace(body.FileName, " ", "_", -1))

	//what is presign?
	//what is s3.PutObjectInput?
	res, err := as.S3Service.S3PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(as.S3Service.BucketName),
		Key:           key,
		ContentType:   aws.String(body.MimeType),
		ContentLength: aws.Int64(int64(body.FileSize)),
		// TODO: add expiration
	})

	if err != nil {
		return nil
	}

	newFile := m.File{
		Name:      body.FileName,
		Size:      fileSize,
		Extension: strings.Split(body.FileName, ".")[1],
		MimeType:  body.MimeType,
		R2Path:    *key,
		UserId:    userId,
	}

	// read it from service ✅
	f, err := as.FileService.CreateFile(&newFile)

	if err != nil {
		return err
	}

	//FIXME: two ID?
	return c.JSON(fiber.Map{
		"file_id":    f.ID,
		"file_name":  f.Name,
		"upload_url": res.URL,
		"is_public":  f.IsPublic,
	})
}
