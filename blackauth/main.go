package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const cookieName = "isuports_session"

var key rsa.PrivateKey

func getTenantName(domain string) string {
	return strings.Split(domain, ".")[0]
}

func getNameParam(r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}
	name := r.Form.Get("name")
	if name == "" {
		return "", fmt.Errorf("name is not found")
	}
	return name, nil
}

func loginPlayerHandler(w http.ResponseWriter, r *http.Request) {
	name, err := getNameParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tenant := getTenantName(r.Host)

	token := jwt.New()
	token.Set(jwt.IssuerKey, "isuports")
	token.Set(jwt.SubjectKey, name)
	token.Set(jwt.AudienceKey, tenant)
	token.Set("role", "player")
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour).Unix())

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privateKey))
	if err != nil {
		fmt.Println("error jwt.Sign: %w", err)
		return
	}

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: fmt.Sprintf("%s", signed),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func loginOrganizerHandler(w http.ResponseWriter, r *http.Request) {
	tenant := getTenantName(r.Host)

	token := jwt.New()
	token.Set(jwt.IssuerKey, "isuports")
	token.Set(jwt.SubjectKey, "organizer")
	token.Set(jwt.AudienceKey, tenant)
	token.Set("role", "organizer")
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour).Unix())

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privateKey))
	if err != nil {
		fmt.Println("error jwt.Sign: %w", err)
		return
	}

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: fmt.Sprintf("%s", signed),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func loginAdminHandler(w http.ResponseWriter, r *http.Request) {
	token := jwt.New()
	token.Set(jwt.IssuerKey, "isuports")
	token.Set(jwt.SubjectKey, "admin")
	token.Set(jwt.AudienceKey, "admin")
	token.Set("role", "admin")
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour).Unix())

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privateKey))
	if err != nil {
		fmt.Println("error jwt.Sign: %w", err)
		return
	}

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: fmt.Sprintf("%s", signed),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

var privateKey *rsa.PrivateKey

func init() {
	// load private key
	//pemFilePath := os.Getenv("ISUCON_PEM_PATH")
	pemFilePath := "isuports.pem"
	f, err := os.Open(pemFilePath)
	if err != nil {
		log.Fatalf("failed to open pem file: %w", err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read file: %w", err)
	}
	block, _ := pem.Decode(buf)
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %w", err)
	}
}

func main() {

	// setup handler
	http.HandleFunc("/api/player/login", loginPlayerHandler)
	http.HandleFunc("/api/organizer/login", loginOrganizerHandler)
	http.HandleFunc("/api/admin/login", loginAdminHandler)
	http.ListenAndServe(":3001", nil)
}