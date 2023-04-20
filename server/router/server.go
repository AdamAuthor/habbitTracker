package router

import (
	"awesomeProject/server/db"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

type Server struct {
	ctx        context.Context
	idleConsCh chan struct{}
	Address    string
	database   db.Database
}

func NewServer(ctx context.Context, address string, database db.Database) *Server {
	return &Server{
		ctx:        ctx,
		idleConsCh: make(chan struct{}),
		database:   database,
		Address:    address,
	}
}

func (s *Server) basicHandler() chi.Router {

	googleOauthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint}
	//TODO: addRandomizer
	randomState := "random"

	r := chi.NewRouter()
	r.Get("/main", func(w http.ResponseWriter, r *http.Request) {
		var html = `<html>
			<body>
				<form action="/register" method="post">
					<label for="email">Email:</label>
					<input type="email" id="email" name="email" required>
					<label for="password">Password:</label>
					<input type="password" id="password" name="password" required>
					<input type="submit" value="Register">
				</form>
				<a href="/loginGoogle">Войти через Google</a>
			</body>
		</html>`
		fmt.Fprintf(w, html)

	})

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		// Получаем email и password из параметров POST запроса
		email := r.FormValue("email")
		content := "Thank you for registering with our service! Please click the link below to confirm your registration:\n\n" +
			"http://example.com/confirm?email=" + email
		log.Println(email)

		// Создаем сообщение для отправки
		msg := []byte("To: " + email + "\r\n" +
			"Subject: Confirm registration\r\n" +
			"\r\n" +
			content)

		// pszzdoosimhclkcv
		// Отправляем сообщение на SMTP серверВ
		err := smtp.SendMail("smtp.gmail.com:587",
			smtp.PlainAuth("", "ahmediarolzasov@gmail.com", "pszzdoosimhclkcv", "smtp.gmail.com"),
			"ahmediarolzasov@gmail.com", // От кого
			[]string{email},             // Кому
			msg)
		if err != nil {
			log.Println("Error sending email:", err)
			return
		}

		log.Println("Email sent to:", email)
	})

	r.Get("/loginGoogle", func(w http.ResponseWriter, r *http.Request) {
		url := googleOauthConfig.AuthCodeURL(randomState)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != randomState {
			fmt.Println("stage is not valid")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
		if err != nil {
			fmt.Printf("couldn`t get token: %s\n\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			fmt.Printf("couldn`t create get request: %s\n\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("couldn`t parse response: %s\n\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		fmt.Fprintf(w, "Response: %s", content)
	})

	return r
}

// Run is the function for running server
func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println("[HTTP] Got err while shutting down:", err)
	}

	log.Println("[HTTP] Processed all idle connections")
	close(s.idleConsCh)
}

func (s *Server) WaitForGT() {
	<-s.idleConsCh
}
