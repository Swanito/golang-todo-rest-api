package service

import (
	"errors"
	"go-playground/model"
	"go-playground/templates"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(user *model.User) (string, error) {
	createAt := time.Now().Unix()
	expireAt := time.Now().AddDate(0, 0, 10).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  user.Id,
		"created": createAt,
		"exp":     expireAt,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", errors.New(err.Error())
	}
	return tokenString, nil
}

func CheckUsername(username string, repository model.UserRepository) *model.CustomError {
	storedUsername, errUsername := repository.GetUserByUsername(username)
	if errUsername != nil {
		return &model.CustomError{Error: errors.New(errUsername.Error()), StatusCode: http.StatusServiceUnavailable}
	}
	if storedUsername.Username != "" {
		return &model.CustomError{Error: errors.New("Username already exists"), StatusCode: http.StatusConflict}
	}
	return nil
}

func CheckEmail(email string, repository model.UserRepository) *model.CustomError {
	storedEmail, errEmail := repository.GetUserByEmail(email)
	if errEmail != nil {
		return &model.CustomError{Error: errors.New(errEmail.Error()), StatusCode: http.StatusServiceUnavailable}
	}
	if storedEmail.Email != "" {
		return &model.CustomError{Error: errors.New("Email already exists"), StatusCode: http.StatusConflict}
	}
	return nil
}

func LoginService(email string, password string, repository model.UserRepository) (*model.LoginResponse, *model.CustomError) {
	maxLoginAttempts, _ := strconv.Atoi(os.Getenv("MAX_LOGIN_ATTEMPTS"))
	user, userErr := repository.GetUserByEmail(email)
	if user.Email == "" {
		return nil, &model.CustomError{Error: errors.New("User not found."), StatusCode: http.StatusNotFound}
	}
	if userErr != nil {
		return nil, &model.CustomError{Error: errors.New(userErr.Error()), StatusCode: http.StatusInternalServerError}
	}
	loginAttempts, _ := repository.GetLoginAttempts(email)
	if loginAttempts == maxLoginAttempts {
		return nil, &model.CustomError{Error: errors.New("Max login attempt reached. Account Blocked."), StatusCode: http.StatusForbidden}
	}
	if !checkPasswordHash(password, user.Password) {
		err := repository.AddInvalidLoginAttempt(email)
		if err != nil {
			return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusForbidden}
		}
		error := errors.New("Invalid credentials. " + strconv.FormatInt(int64(maxLoginAttempts-loginAttempts), 10) + " out of " + strconv.FormatInt(int64(maxLoginAttempts), 10) + "attempts left.")
		return nil, &model.CustomError{Error: error, StatusCode: http.StatusUnauthorized}
	}
	if !user.Activated {
		SendActivationEmail(user.Username, user.Email)
		return nil, &model.CustomError{Error: errors.New("User not activated. Please check your email to activate the user."), StatusCode: http.StatusUnauthorized}
	}
	errSession := repository.AddLoginSession(user.Id, time.Now())
	if errSession != nil {
		return nil, &model.CustomError{Error: errors.New(errSession.Error()), StatusCode: http.StatusInternalServerError}
	}
	token, err := generateToken(user)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(errSession.Error()), StatusCode: http.StatusInternalServerError}
	}
	response := model.LoginResponse{
		Username:    user.Username,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		AccessToken: token,
	}
	return &response, nil
}

func RegisterService(username, password, confirmPassword, email string, ur model.UserRepository) (int, *model.CustomError) {
	if password != confirmPassword {
		return 0, &model.CustomError{Error: errors.New("Passwords does not match"), StatusCode: http.StatusBadRequest}
	}
	passwordHash, _ := hashPassword(password)
	createdAt := time.Now()
	userId, err := ur.Register(username, passwordHash, email, createdAt)
	if err != nil {
		return 0, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return userId, nil
}

func SendActivationEmail(username, address string) *model.CustomError {
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"

	from := mail.NewEmail("TODO Team", "scrapertestingebay@gmail.com")
	to := mail.NewEmail(username, address)
	subject := "TODO: Activate your account"
	content := mail.NewContent("text/html", templates.ActivationEmailTemplate(username, address))
	m := mail.NewV3MailInit(from, subject, to, content)
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		return &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusServiceUnavailable}
	}
	return nil
}

func ActivateUser(userEmail string, repository model.UserRepository) *model.CustomError {
	err := repository.ActivateUser(userEmail)
	if err != nil {
		return &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return nil
}
