package admin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"newanysock/internal/entity"
	"newanysock/internal/usecase/repo"
	"newanysock/pkg/database/mysql"
	"time"
)

func UpdateLisence(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	licenseExist, rawStruct := repo.MatchRecord("license_key", data["license"], &entity.License{})
	if licenseExist == false {
		return c.JSON(fiber.Map{"error": "License does not exists"})
	}
	if rawStruct == nil {
		return c.JSON(fiber.Map{"error": "can not found your user"})
	}
	object := rawStruct.(*entity.License)

	var enable bool = false

	if data["enable"] == "true" {
		enable = true
	} else {
		enable = false
	}

	newLicense := entity.License{
		Id:           object.Id,
		LicenseKey:   object.LicenseKey,
		HardwareCode: object.HardwareCode,
		Enable:       enable,
	}

	var rs = mysql.DB.Save(&newLicense)

	if rs.Error != nil {
		return c.JSON(fiber.Map{"error": rs.Error.Error()})
	}
	if rs.RowsAffected == 1 {
		return c.JSON(fiber.Map{"success": "true"})
	}
	return nil
}

func RemoveLisence(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	fmt.Println(data["license"])
	licenseExist, rawStruct := repo.MatchRecord("license_key", data["license"], &entity.License{})
	if licenseExist == false {
		return c.JSON(fiber.Map{"error": "License does not exists"})
	}
	if rawStruct == nil {
		return c.JSON(fiber.Map{"error": "can not found your user"})
	}
	object := rawStruct.(*entity.License)
	var rs = mysql.DB.Delete(&object)

	if rs.Error != nil {
		return c.JSON(fiber.Map{"error": rs.Error.Error()})
	}
	if rs.RowsAffected == 1 {
		return c.JSON(fiber.Map{"success": "true"})
	}
	return nil
}

func AddLicense(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	endAtTime, err := time.Parse("2006-01-02 03:04:05", data["endAt"])
	if err != nil {

		return c.JSON(fiber.Map{"error": "Could not parse time"})
	}

	var object = entity.License{
		LicenseKey:   data["license"],
		HardwareCode: data["hardwareCode"],
		Enable:       true,
		StartAt:      time.Now(),
		EndAt:        endAtTime,
		Owner:        data["owner"],
	}

	licenseExist, _ := repo.MatchRecord("license_key", data["license"], &entity.License{})
	if licenseExist == true {
		return c.JSON(fiber.Map{"error": "License Key already exists"})
	}

	result := mysql.DB.Create(&object)

	if result.Error != nil {
		return c.JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.JSON(object)
}
