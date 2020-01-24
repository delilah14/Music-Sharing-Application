package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	_ "github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
)

const (
	privKeyPath = "path/to/keys/app.rsa"
	pubKeyPath  = "path/to/keys/app.rsa.pub"
)

var VerifyKey, SignKey []byte

func NewRouter() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/hello", Handler).Methods("GET")
	//r.HandleFunc("/users", GetUsersHandler).Methods("GET")
	r.HandleFunc("/register", GetRegisterHandler).Methods("GET")
	r.HandleFunc("/registeruser", RegisterUserHandler).Methods("POST")
	r.HandleFunc("/login-form", GetLoginFormHandler).Methods("GET")
	r.HandleFunc("/login", AuthLoginHandler).Methods("POST")
	r.HandleFunc("/success", GetUsersHandler).Methods("POST")
	r.HandleFunc("/allmusic", GetAllMusicHandler).Methods("GET")
	r.HandleFunc("/addmusic-page", GetAddMusicPage).Methods("GET")
	r.HandleFunc("/addmusic", AddMusicHandler).Methods("POST")
	r.HandleFunc("/searchmusic", SearchMusicHandler).Methods("POST")
	r.HandleFunc("/allplaylist", GetAllPlaylistHandler).Methods("GET")
	r.HandleFunc("/addplaylist", AddPlaylistHandler).Methods("POST")
	r.HandleFunc("/viewsingleplaylist", ViewSinglePlaylistHandler).Methods("POST")
	r.HandleFunc("/addmsc2plylst", AddMusicToPlaylistHandler).Methods("POST")
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/")))
	r.PathPrefix("/static/").Handler(s)

	r.Handle("/users/{token}", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(GetUsersHandler)),
	))

	return r
}
func initKeys() {
	var err error

	SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}

}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//validate token
	println(r.Header.Get("token"))
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return VerifyKey, nil
	})

	println(token)

	if err == nil {

		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
	}

}
