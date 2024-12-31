package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("Sorry route: " + r.URL.Path + " does not exist. check your url"))
}

// to help validate register request body
type RegisterUserBody struct {
	User_name string `json:"user_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int32  `json:"age"`
}

func (rg RegisterUserBody) ValidateUserName() bool {
	return len(strings.TrimSpace(rg.User_name)) >= 2
}

func ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func ValidateStrongPassword(pwd string) bool {
	passwordRegex := regexp.MustCompile(`^[a-zA-Z\d!@#$%^&*(),.?":{}|<>]+$`)
	return passwordRegex.MatchString(pwd) && len(pwd) >= 6
}

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

func (rg RegisterUserBody) HashUserPassword() string {

	generatedHash, cost := bcrypt.GenerateFromPassword([]byte(rg.Password), MinCost)
	fmt.Println(cost)

	hashedPassword := string(generatedHash)
	return hashedPassword

}

type Response struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
	Success    bool        `json:"success"`
	Error      string      `json:"error,omitempty"`
}

func CustomResponseInJson(w http.ResponseWriter, statusCode int, data interface{}, err error) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Create response object
	response := Response{
		StatusCode: statusCode,
		Success:    statusCode >= 200 && statusCode < 300,
	}

	if data != nil {
		response.Data = data
	}

	// Add error message if present
	if err != nil {
		response.Error = err.Error()
	}

	// Encode and send the response
	json.NewEncoder(w).Encode(response)

}

func GenerateJwtToken(claims jwt.MapClaims) string {
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	JWT_KEY := os.Getenv("JWT_KEY")
	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = []byte(JWT_KEY)

	t = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	createdToken, err := t.SignedString(key)

	if err != nil {
		return fmt.Sprintf("Error occurred %v", err)

	}

	s = createdToken

	return s
}

func ValidateJwtToken(token string) (interface{}, error) {

	JWT_KEY := os.Getenv("JWT_KEY")
	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		fmt.Println()
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		secretKey := []byte(JWT_KEY)

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return tk.Claims, nil

}
