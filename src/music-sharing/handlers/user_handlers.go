package handlers

import (
	"fmt"
	"html/template"
	"log"
	"music-sharing/models/db"
	"music-sharing/models/track"
	"music-sharing/utilities"
	_ "music-sharing/utilities"
	"net/http"
	"time"

	_ "github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go/request"
)

type User struct {
	Username        string
	FirstName       string
	LastName        string
	Email           string
	Password        string
	SearchValue     string
	Success         bool
	DisplayTracks   []track.Track
	DisplayPlaylist []track.Playlist
	NotAvailable    bool
	PlaylistName    string
	TokenString     string
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

var database = db.InitDatabase()

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := VerifyKey
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("token")
	u := User{}
	println(user)
	claimms, isClaim := extractClaims(user)
	println(isClaim, claimms)
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/index.html")

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = tmpl.Execute(w, u) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {          // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}

func Handler(w http.ResponseWriter, r *http.Request) {
	var user track.User
	database.Where("username = ?", "jinzhu").First(&user)
	fmt.Fprintf(w, user.FirstName)
}

func GetUsersPageHandler(w http.ResponseWriter, r *http.Request) {
	u := User{}
	t, err := template.ParseFiles(utilities.TemplatePath + "/assets/index.html") //parse the template file held in the templates folder

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, u) //execute the template and pass in the variables to fill the gaps

	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}
func GetLoginFormHandler(w http.ResponseWriter, r *http.Request) {
	u := User{}
	t, err := template.ParseFiles(utilities.TemplatePath + "/assets/login.html") //parse the template file held in the templates folder

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, u) //execute the template and pass in the variables to fill the gaps

	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func AuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	user := track.User{}
	session := User{Success: true}
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/login.html")
	println("I am in login handler")

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	details := User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	println(details.Email)

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	database.Where(&User{Email: details.Email, Password: details.Password}).First(&user)
	//database.Where("email = ? AND password = ?", details.Email, details.Password).Find(&user)
	if user.Password == details.Password {
		//create a rsa 256 signer
		signer := jwt.New(jwt.GetSigningMethod("RS256"))

		//set claims
		claims := make(jwt.MapClaims)
		claims["iss"] = user.Username
		claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
		claims["CustomUserInfo"] = struct {
			Name string
			Role string
		}{user.Username, "User"}

		signer.Claims = claims

		tokenString, err := signer.SignedString(SignKey)

		println(tokenString, err)
		session.Email = user.Email
		session.FirstName = user.FirstName
		session.LastName = user.LastName
		session.Password = user.Password
		session.Username = user.Username
		session.Success = false
		session.TokenString = tokenString
		println(user.Username)
		tmpl, err = template.ParseFiles(utilities.TemplatePath + "/assets/index.html")
		tmpl.Execute(w, session)
		r.Header.Add("token", tokenString)
		return
	}
	if user.Password != details.Password {
		println(user.Username)
		//execute the template and pass it the session struct to fill in the gaps
		session.Success = true
		if err != nil { // if there is an error
			log.Print("template executing error: ", err) //log it
		}
		tmpl, err = template.ParseFiles(utilities.TemplatePath + "/assets/login.html")
		tmpl.Execute(w, session)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	user := track.User{}
	session := User{Success: true}
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/login.html")
	println("I am in login handler")

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	details := User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	println(details.Email)

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	database.Where(&User{Email: details.Email, Password: details.Password}).First(&user)
	//database.Where("email = ? AND password = ?", details.Email, details.Password).Find(&user)
	if user.Password == details.Password {
		println("I'm in")
		session.Email = user.Email
		session.FirstName = user.FirstName
		session.LastName = user.LastName
		session.Password = user.Password
		session.Username = user.Username
		session.Success = false
		println(user.Username)
		tmpl, err = template.ParseFiles(utilities.TemplatePath + "/assets/index.html")
		tmpl.Execute(w, session)
		return
	}
	if user.Password != details.Password {
		println(user.Username)
		//execute the template and pass it the session struct to fill in the gaps
		session.Success = true
		if err != nil { // if there is an error
			log.Print("template executing error: ", err) //log it
		}
		tmpl, err = template.ParseFiles(utilities.TemplatePath + "/assets/login.html")
		tmpl.Execute(w, session)
		return
	}

}

func GetRegisterHandler(w http.ResponseWriter, r *http.Request) {
	u := User{}
	t, err := template.ParseFiles(utilities.TemplatePath + "/assets/register.html") //parse the template file held in the templates folder

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, u) //execute the template and pass in the variables to fill the gaps

	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/register.html")
	println("I am in register handler")

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	details := track.User{
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("firstname"),
		LastName:  r.FormValue("lastname"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
	}
	println(details.FirstName)

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	database.Create(&details)
	tmpl.Execute(w, struct{ Success bool }{true})
	//execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}
