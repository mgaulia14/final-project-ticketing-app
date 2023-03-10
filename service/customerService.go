package service

import (
	"errors"
	"final-project-ticketing-api/database"
	"final-project-ticketing-api/middleware"
	"final-project-ticketing-api/repository"
	"final-project-ticketing-api/structs"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/mail"
	"regexp"
	"time"
)

const nameRegex = `^[A-Za-z]+(?:[ _-][A-Za-z]+)*$`
const phoneNumberRegex = `^(^\+62|62|^08)(\d{3,4}-?){2}\d{3,4}$`

func GetAllCustomer() (structs.Customer, error) {
	var result structs.Customer
	err, result := repository.GetAllCustomer(database.DBConnection)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetCustomerById(customerId int) (structs.Customer, error) {
	var result structs.Customer
	err, result := repository.GetByCustomerId(database.DBConnection, customerId)
	if err != nil {
		return result, err
	}
	return result, nil
}

func CreateCustomer(request structs.CustomerRequest) (structs.Customer, []error) {
	customer, err := prepareRequestCustomer(request)
	if err != nil {
		return customer, err
	}
	customer, err = repository.InsertCustomer(database.DBConnection, customer)
	err1, cust := repository.GetCustomerByEmail(database.DBConnection, request.Email)
	if err1 != nil {
		err = append(err, err1)
		return cust, err
	}
	// create wallet account
	var wallet structs.Wallet
	wallet.AccountNumber = generateNumber(10000000, 99999999)
	wallet.AccountName = request.FullName
	wallet.Balance = 0
	wallet.CustomerId = cust.ID
	wallet, err = repository.InsertWallet(database.DBConnection, wallet)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func UpdateCustomer(request structs.CustomerRequest, customerId int) (structs.Customer, []error) {
	var result structs.Customer
	var err []error
	err1, _ := repository.GetByCustomerId(database.DBConnection, customerId)
	if err1 != nil {
		err = append(err, err1)
		return result, err
	}
	customer, err := prepareRequestCustomer(request)
	customer.ID = customerId
	if err != nil {
		return customer, err
	}
	customer, err = repository.UpdateCustomer(database.DBConnection, customer)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func Login(request structs.CustLogin) (structs.Customer, []error) {
	var err []error
	// check username and password correct
	err1, cust := repository.GetCustomerByEmail(database.DBConnection, request.Email)
	if err1 != nil {
		err = append(err, err1)
		return cust, err
	}
	match := CheckPasswordHash(request.Password, cust.Password)
	if match == false {
		err = append(err, errors.New("incorrect password"))
		return cust, err
	}
	var role string
	if cust.IsAdmin == true {
		role = "admin"
	} else {
		role = "user"
	}
	token, err2 := middleware.GenerateJWT(request.Email, request.Password, role)
	if err2 != nil {
		err = append(err, err2)
		return cust, err
	} else {
		cust.Token = token
		return cust, err
	}
}

func prepareRequestCustomer(request structs.CustomerRequest) (structs.Customer, []error) {
	var customer structs.Customer
	request, err, dateCustomer := validateRequestCustomer(request)
	if err != nil {
		return customer, err
	}
	hashPass, _ := HashPassword(request.Password)
	customer.FullName = request.FullName
	customer.BirthDate = dateCustomer
	customer.Address = request.Address
	customer.PhoneNumber = request.PhoneNumber
	customer.Email = request.Email
	customer.Password = hashPass
	customer.IsAdmin = request.IsAdmin
	return customer, nil
}

func validateRequestCustomer(request structs.CustomerRequest) (structs.CustomerRequest, []error, time.Time) {
	var dateInt []int
	var err []error
	dateRequest := request.BirthDate
	err0, _ := repository.GetCustomerByEmail(database.DBConnection, request.Email)
	if err0 == nil {
		err = append(err, errors.New("email already registered"))
	}
	if isValidName(request.FullName) == false {
		err = append(err, errors.New("parameter 'full_name' must be in alphabet only"))
	}
	regex, _ := regexp.Compile(formatDate)
	if !regex.MatchString(dateRequest) {
		err = append(err, errors.New("parameter 'birth_date' must be in format yyyy-MM-dd"))
		panic(err)
	}
	dateCustomer, err1 := GetDate(dateRequest, dateInt)
	if err1 != nil {
		err = append(err, errors.New("parameter 'birth_date' must be in format yyyy-MM-dd"))
	}
	if isValidPhoneNumber(request.PhoneNumber) == false {
		err = append(err, errors.New("prefix 'phone_number' must be '08' or '62' or '+62' and max 13 digit"))
	}
	if isValidEmail(request.Email) == false {
		err = append(err, errors.New("invalid 'email' format"))
	}
	if len(err) > 0 {
		return request, err, dateCustomer
	}
	return request, nil, dateCustomer
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPhoneNumber(phone string) bool {
	regex, _ := regexp.Compile(phoneNumberRegex)
	if !regex.MatchString(phone) {
		return false
	}
	return true
}

func isValidName(name string) bool {
	regex, _ := regexp.Compile(nameRegex)
	if !regex.MatchString(name) {
		return false
	}
	return true
}

func generateNumber(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
