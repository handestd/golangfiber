package admin

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"newanysock/internal/entity/user"
	"newanysock/internal/usecase/repo"
	"newanysock/pkg"
	"newanysock/pkg/database/mysql"
)

func UpdateUser(c *fiber.Ctx) error {

	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	nilValidate := pkg.IsMapContainNil(data, []string{"groupid", "avatar", "role", "lastip"})

	if nilValidate {
		return c.JSON(fiber.Map{"error": "you send empty value"})
	}

	var u = user.Account{}
	if err := c.BodyParser(&u); err != nil {
		return err
	}

	userExist, rawAccount := repo.MatchRecord("id", u.Id, &user.Account{})

	if userExist == false {
		return c.JSON(fiber.Map{"error": "User does not exists"})
	}
	if rawAccount == nil {
		return c.JSON(fiber.Map{"error": "can not found your user"})
	}

	account := rawAccount.(*user.Account)

	newPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	_, ok := data["password"]
	// If the key exists
	if ok {
		data["password"] = string(newPassword)
	}

	rs := mysql.DB.Model(&account).Updates(data)
	if rs.Error != nil {
		return c.JSON(fiber.Map{"error": rs.Error.Error()})
	}
	if rs.RowsAffected == 1 {
		return c.JSON(fiber.Map{"success": "true"})
	}

	return nil
}
func DeleteUser(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	// comment this line
	usernameExist, rawAccount := repo.MatchRecord("username", data["username"], &user.Account{})
	if usernameExist == false {
		return c.JSON(fiber.Map{"error": "Username does not exists"})
	}
	if rawAccount == nil {
		return c.JSON(fiber.Map{"error": "can not found your user"})
	}
	account := rawAccount.(*user.Account)
	var rs = mysql.DB.Delete(&account)

	if rs.Error != nil {
		return c.JSON(fiber.Map{"error": rs.Error.Error()})
	}
	if rs.RowsAffected == 1 {
		return c.JSON(fiber.Map{"success": "true"})
	}
	return nil
}
func ShowUsers(c *fiber.Ctx) error {

	//var data map[string]string
	//
	//if err := c.BodyParser(&data); err != nil {
	//	return err
	//}

	var records = []user.Account{}

	rs := mysql.DB.Find(&records)
	if rs.Error != nil {
		return c.JSON(fiber.Map{"error": rs.Error.Error()})
	}
	if rs.RowsAffected > 0 {
		return c.JSON(records)
	}
	return nil
}

func AddUser(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	var account = user.Account{
		Username: data["username"],
		Password: string(password),
		Email:    data["email"],
	}

	usernameExist, _ := repo.MatchRecord("username", data["username"], &user.Account{})
	if usernameExist == true {
		return c.JSON(fiber.Map{"error": "Username already exists"})
	}
	emailExist, _ := repo.MatchRecord("email", data["email"], &user.Account{})
	if emailExist == true {
		return c.JSON(fiber.Map{"error": "Email already exists"})
	}

	result := mysql.DB.Create(&account)

	if result.Error != nil {
		return c.JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.JSON(account)
}
