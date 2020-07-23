package config

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// general function to validate all kind of api request payload / body
func ValidateRequestPayload(c echo.Context, rules govalidator.MapData, data interface{}) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Data:    data,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	mappedError := v.ValidateJSON()
	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

func ValidateRequestFormData(c echo.Context, rules govalidator.MapData) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Rules:   rules,
	}

	v := govalidator.New(opts)
	mappedError := v.Validate()
	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

// general function to validate all kind of api request url query
func ValidateRequestQuery(c echo.Context, rules govalidator.MapData) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Rules:   rules,
	}

	v := govalidator.New(opts)

	mappedError := v.Validate()

	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

func ReturnInvalidResponse(httpcode int, details interface{}, message string) error {
	responseBody := map[string]interface{}{
		"message": message,
		"details": details,
	}

	return echo.NewHTTPError(httpcode, responseBody)
}

//encrypt data with AES 256 CFB
func Encrypt(text string) string {
	var sha = sha1.New()
	sha.Write([]byte(text))
	var encrypted = sha.Sum(nil)
	return fmt.Sprintf("%x", encrypted)
}

func Upload(c echo.Context, prefix string, userID string) (string, error) {

	file, err := c.FormFile(prefix)
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	filename := fmt.Sprintf("%s-%s%s", prefix, userID, filepath.Ext(file.Filename))

	fileLocation := filepath.Join(App.Config.GetString("dir"), filename)

	// Destination
	dst, err := os.Create(fileLocation)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filename, err
}

// RandString random string alphanumeric. parameter length
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func DiffDate(a time.Time, b time.Time) int {

	diff := a.Sub(b)

	// convert diff to days
	return int(diff.Hours())
}

func DiffTime(a time.Time, b time.Time) int {

	diff := a.Sub(b)

	// convert diff to days
	return int(diff.Seconds())
}
