package helpers

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// Similar to Laravel's dd, help with debugging
func Dd(myVar ...any) {
	fmt.Printf("%v\n", myVar)
	os.Exit(0)
}

func ParseFloatFromMap(m map[string]any, key string) float64 {
	strValue := fmt.Sprintf("%v", m[key])
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		fmt.Printf("helpers.ParseFloatFromMap: error parsing to float: %v\n", err)
		return float64(0)
	}
	return value
}

// This function for convert days to months
func ParseDaysToMonths(day int) int {
	var daysPerMonth = 30.44
	months := float64(day) / daysPerMonth
	return int(math.Ceil(months))
}

// This function for calculate installment
func CalculateInstallment(receivableTotal float64, months int) float64 {
	if months <= 0 {
		return receivableTotal
	}
	return receivableTotal / float64(months)
}

// This function fro parse string to int
func ParseStringToInt(str string) int {
	strValue := fmt.Sprintf("%v", str)
	value, valueErr := strconv.Atoi(strValue)
	if valueErr != nil {
		fmt.Printf("helpers.ParseIntFromMap: error parsing to int: %v\n", valueErr)
		return int(0)
	}
	return value
}

func CheckHashedPasswordMatches(HashedPassword, InputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPassword), []byte(InputPassword))
	if err != nil {
		return err
	}
	return nil
}

func GeneratePassword(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	password := make([]byte, length)
	for i := range password {
		password[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(password)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ConvertStatusToString(status int) string {
	var convertedStatus string
	switch status {
	case 1:
		convertedStatus = "Aktif"
	case 0:
		convertedStatus = "Non - Aktif"
	}
	return convertedStatus
}
