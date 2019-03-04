package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/kabukky/httpscerts"
	"github.com/kaleempeeroo/wheresmydosh/server/cmd/com/wheresmydosh/db"
	"github.com/xlzd/gotp"
)

const (
	privKeyPath = "key.pem"
	certPath    = "cert.pem"
)

var mySigningKey []byte
var validationJWTS JWTS
var testToken Token
var user db.User

type JWT struct {
	UUIDKey uuid.UUID
	OTP     string
	User    db.User
	Verfied bool
}

func (jwts *JWTS) AddItem(item JWT) []JWT {
	jwts.Jwts = append(jwts.Jwts, item)
	return jwts.Jwts
}

type JWTS struct {
	Jwts []JWT
}

type Token struct {
	AccessToken string `json:"access_token"`
}

// Jwks stores a slice of JSON Web Keys
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	fmt.Println(q.FormattedQuery())
}

func main() {
	dbConnection := connectToDB()
	createDBTables(dbConnection)
	generateKeys()
	login()
	startServer()
	SMSValidationCode("", "+23057352919")

}

// start server with SSL connection endpoint
func startServer() {
	r := gin.Default()
	r.POST("/register", authMiddleware(), registerClient)
	r.POST("/verify_phone", authMiddleware(), verifyPhone)
	err := r.RunTLS(":8080", certPath, privKeyPath)
	fmt.Println(err)

}

func registerClient(c *gin.Context) {

	// get user registration data from POST request body
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&user)
	fmt.Printf("error decoding %v", err)
	fmt.Printf("user %v %v %v %v %v", user.FirstName, user.MiddleName, user.LastName, user.Email, user.PhoneNumber)
	if err != nil {
		log.Println(err)
	}

	// set response return type
	c.Header("Content-Type", "application/json")

	// generate UUID
	UUIDKey, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v UUID: %v", err, UUIDKey)

	//generate OTP
	otp := generateValidationCode()

	//send OTP via cellular network
	SMSValidationCode(otp, "+23057352919")

	// return UUID + OTP in response.  OTP needs to be removed, only testing. SMS OTP
	c.SetCookie("UUID", UUIDKey.String(), 3600, "/", "localhost", true, true)
	c.SetCookie("OTP", otp, 3600, "/", "localhost", true, true)
	c.Writer.WriteHeader(http.StatusOK)

	// store UUID and OTP for client in JWT struct. Add JWT to JWTS
	validationJWTS.AddItem(JWT{UUIDKey, otp, user, false})

}

func verifyPhone(c *gin.Context) {
	// req = token + OTP + UUID + phoneNumber
	requestUUID, err := c.Cookie("UUID")
	log.Println(err)
	otp, err := c.Cookie("OTP")
	log.Println(err)
	phoneNumber, err := c.Cookie("phone_number")
	log.Println(err)

	fmt.Println(err)
	fmt.Println(requestUUID)
	fmt.Println(phoneNumber)

	for _, jwt := range validationJWTS.Jwts {
		if jwt.UUIDKey.String() == requestUUID {
			if jwt.OTP == otp {
				// Hint Android API can allow user to choose multiple numbers. Check here.
				if jwt.User.PhoneNumber == phoneNumber {
					jwt.Verfied = true
					fmt.Println("verified")
				}
			}
		}
	}
	
}

// Test Auth0 with Test Token
func login() Token {
	url := "https://dev-wheresmydosh.auth0.com/oauth/token"

	payload := strings.NewReader("{\"client_id\":\"9fLvd37b2s6Fi8XzBF1hFFCcXlorMtke\",\"client_secret\":\"maV_NO4ktxR27v9XmcCDt_9K1WOpgmo7GrJ46bbHOATRZvCfKk3v27bAXCUgfKm7\",\"audience\":\"wheresmydosh\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	err := json.Unmarshal(body, &testToken)
	fmt.Println(err)
	fmt.Println(testToken.AccessToken)
	return testToken
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		aud := "wheresmydosh"
		checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAudience {
			return token, errors.New("Invalid audience.")
		}
		// verify iss claim
		iss := "https://dev-wheresmydosh.auth0.com/"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}

		cert, err := getPemCert(token)
		if err != nil {
			log.Fatalf("could not get cert: %+v", err)
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
})

// authMiddleware intercepts the requests, and check for a valid jwt token
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client secret key
		err := jwtMiddleware.CheckJWT(c.Writer, c.Request)
		if err != nil {
			// Token not found
			fmt.Println(err)
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}
	}
}

// Using AWS SNS service to sms OTP

func SMSValidationCode(otp string, phone string) {
	// need to set MessageType to Transactional
	mySession := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	snsClient := sns.New(mySession)
	message := "Testing SMS from AWS ... ismail" + otp
	fmt.Println(phone)
	pi := &sns.PublishInput{Message: &message, PhoneNumber: &phone}
	output, err := snsClient.Publish(pi)
	fmt.Println(output)
	fmt.Println(err)
}

func connectToDB() *pg.DB {
	dba := pg.Connect(&pg.Options{
		User:     "root",
		Password: "password",
		Database: "postgres",
		Addr:     "wheresmydosh-cluster.cluster-cnuxmvkomgbc.us-east-2.rds.amazonaws.com:5432",
	})
	dba.AddQueryHook(dbLogger{})
	var n int
	_, err := dba.QueryOne(pg.Scan(&n), "SELECT 1")
	fmt.Println(err)
	fmt.Println(n)
	return dba
}

func createDBTables(database *pg.DB) {
	db.CreateDBTables(database)
}

// internal function used by server to generate OTP
func generateValidationCode() string {
	secretLength := 16
	secret := gotp.RandomSecret(secretLength)
	fmt.Println(secret)
	hotp := gotp.NewDefaultHOTP(secret)
	hotp1 := hotp.At(0)
	fmt.Println(hotp1)
	fmt.Println(hotp.Verify(hotp1, 0))
	return hotp1
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://dev-wheresmydosh.auth0.com/" + ".well-known/jwks.json")
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	x5c := jwks.Keys[0].X5c
	for k, v := range x5c {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + v + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		return cert, errors.New("unable to find appropriate key.")
	}

	return cert, nil
}

// internal function used by server to generate public/private keys and self-signed certificate
func generateKeys() {

	// Check if the cert files are available.
	err := httpscerts.Check("cert.pem", "key.pem")

	// If they are not available, generate new ones.
	if err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8081")
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}
}

func initializeKeys() {
	var err error
	mySigningKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}
}

func getJWT(UUIDKey uuid.UUID) string {
	// get a JWT token containing UUID for client

	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["uuid"] = UUIDKey
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	return tokenString
}

