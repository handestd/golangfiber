package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	errorEntity "newanysock/internal/entity/error"
	"newanysock/internal/entity/user"
	"newanysock/internal/usecase/repo"
	"newanysock/pkg"
	"newanysock/pkg/database/mysql"
	"time"
)

func Profile(c *fiber.Ctx) error {
	data := c.Locals("userClaim")
	if data == nil {
		return c.JSON(fiber.Map{"error": "we cant detect your infomation"})
	}
	if pkg.CompareType(data, user.UserClaim{}) == false {
		return c.JSON(fiber.Map{"error": "we cant detect your infomation"})
	}
	newdata := data.(user.UserClaim)
	usernameExist, rawAccount := repo.MatchRecord("username", newdata.Username, &user.Account{})
	if usernameExist == false {
		return c.JSON(fiber.Map{"error": "Username does not exists"})
	}
	if rawAccount == nil {
		return c.JSON(fiber.Map{"error": "can not found your user"})
	}

	account := rawAccount.(*user.Account)

	return c.JSON(account)
}

func ChangePassword(c *fiber.Ctx) error {
	data := c.Locals("userClaim")
	if data != nil {
		if pkg.CompareType(data, user.UserClaim{}) {
			yourData := data.(user.UserClaim)

			usernameExist, rawAccount := repo.MatchRecord("username", yourData.Username, &user.Account{})
			if usernameExist == false {
				return c.JSON(fiber.Map{"error": "Username does not exists"})
			}

			if rawAccount == nil {
				return c.JSON(fiber.Map{"error": "can not found your user"})
			}

			account := rawAccount.(*user.Account)

			var dataPayload map[string]string

			if err := c.BodyParser(&dataPayload); err != nil {
				return err
			}

			originPassword := account.Password
			oldPassword := dataPayload["oldPassword"]
			newPassword := dataPayload["newPassword"]

			if oldPassword == newPassword {
				return c.JSON(fiber.Map{"error": "old password is same with new password"})
			}

			err := bcrypt.CompareHashAndPassword([]byte(originPassword), []byte(oldPassword))
			if err != nil {
				return c.JSON(fiber.Map{"error": "your old password isn't correct"})
			}
			password, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 14)

			rs := mysql.DB.Model(&user.Account{}).Where("id = ?", account.Id).Update("password", password)

			if rs.Error != nil {
				return c.JSON(fiber.Map{"error": rs.Error.Error()})
			}
			if rs.RowsAffected == 1 {
				return c.JSON(fiber.Map{"success": "true"})
			}

			return nil
		}
	}
	return c.JSON(data)
}

func Login(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	usernameExist, rawAccount := repo.MatchRecord("username", data["username"], &user.Account{})
	if usernameExist == false {
		return c.JSON(fiber.Map{"error": "Username does not exists"})
	}

	if rawAccount == nil {
		return c.JSON(fiber.Map{"error": "can not found your user"})
	}
	//password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	account := rawAccount.(*user.Account)

	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(data["password"]))
	if err != nil {
		return c.JSON(fiber.Map{"error": errorEntity.UserAccountWrongPassword.Error.Error()})
	}

	//jwt
	jwtToken, err := CreateJWTToken(account.Id, account.Email, account.Username, account.Role)

	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	var userClaim user.UserClaim
	err = ParseJWTToken(jwtToken, &userClaim)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		SameSite: "none",
		Secure:   true,
	}

	c.Cookie(&cookie)

	return c.JSON(account)
}

func Register(c *fiber.Ctx) error {
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

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

const key = "301f06b88dde8ca2926bb82544dab952"

func CreateJWTToken(id int, email string, name string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user.UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		ID:               id,
		Email:            email,
		Username:         name,
		Role:             role,
	})

	signedString, err := token.SignedString([]byte(key))

	if err != nil {
		return "", fmt.Errorf("error creating signed string: %v", err)
	}

	return signedString, nil
}
func ParseJWTToken(jwtToken string, userClaim *user.UserClaim) error {
	token, err := jwt.ParseWithClaims(jwtToken, userClaim, func(token *jwt.Token) (interface{}, error) {
		// returning the secret key
		return []byte(key), nil
	})
	if err != nil {
		return err
	}

	// check token validity, for example token might have been expired
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
