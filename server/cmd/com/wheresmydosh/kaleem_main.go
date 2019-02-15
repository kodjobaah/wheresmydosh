package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	fmt.Println(q.FormattedQuery())
}

func main() {
	//dbConnection := connectToDB()
	//createDBTables(dbConnection)
	generateKeys()
	startServer()
}

// start server with SSL connection endpoint
func startServer() {
	r := gin.Default()
	r.POST("/register", authMiddleware(), registerClient)
	err := r.RunTLS(":8080", certPath, privKeyPath)
	fmt.Println(err)
}

func registerClient(c *gin.Context) {
	Firstname := c.PostForm("first_name")
	Middlename := c.PostForm("middle_name")
	Lastname := c.PostForm("last_name")
	PhoneNumber := c.PostForm("phone_number")
	Email := c.PostForm("email")

	c.Header("Content-Type", "application/json")

	// generate UUID
	UUIDKey, err := uuid.NewRandom()
	fmt.Println("%v UUID: %v", err, UUIDKey)

	//generate OTP
	generateValidationCode()

	//send OTP via cellular network

	// send JWT with UUID
	c.JSON(http.StatusOK, gin.H{
		"token":        []byte(getJWT()),
		"First Name":   Firstname,
		"Middle Name":  Middlename,
		"Last Name":    Lastname,
		"Phone Number": PhoneNumber,
		"Email":        Email,
	})
	// store UUID and JWT in JWTKs struct
}

func getJWT() string {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["admin"] = true
	claims["name"] = "Ado Kukic"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	return tokenString
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
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

func connectToDB() *pg.DB {
	dba := pg.Connect(&pg.Options{
		User:     "root",
		Password: "password",
		Database: "postgres",
		Addr:     "wheresmydosh-cluster.cluster-cnuxmvkomgbc.us-east-2.rds.amazonaws.com:5432",
	})
	dba.AddQueryHook(dbLogger{})
	defer dba.Close()
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
