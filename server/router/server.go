package router

import (
	"awesomeProject/server/db"
	"awesomeProject/server/models"
	"awesomeProject/server/random"
	"awesomeProject/server/sendmail"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
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
		RedirectURL:  "https://habit-makers.herokuapp.com/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ready!")
	})

	r.Get("/loginGoogle", func(w http.ResponseWriter, r *http.Request) {
		url := googleOauthConfig.AuthCodeURL(random.GenerateRandomString())
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		token, err := googleOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Unable to exchange code for token", http.StatusInternalServerError)
			return
		}
		// Получаем информацию о пользователе из Google
		oauth2Service := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

		userInfo, err := oauth2Service.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			http.Error(w, "Unable to get user info from Google", http.StatusInternalServerError)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(userInfo.Body)

		gg := new(models.Gmail)
		err = json.NewDecoder(userInfo.Body).Decode(gg)
		if err != nil {
			http.Error(w, "Unable to parse user info", http.StatusInternalServerError)
			return
		}
		if len(gg.Name) == 0 {
			http.Error(w, "Unable to get user name from Google", http.StatusInternalServerError)
			return
		}
		// Проверяем, существует ли пользователь в базе данных
		existingUser, err := s.database.User().FindByEmail(r.Context(), gg.Email)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Unable to check user existence in database", http.StatusInternalServerError)
			return
		}

		var user *models.User
		if existingUser != nil {
			// Пользователь существует, выполняем вход
			user = existingUser
		} else {
			user = &models.User{
				Nickname:  gg.Name,
				Email:     gg.Email,
				Confirmed: true,
			}
			err = s.database.User().CreateGoogle(r.Context(), user)
			if err != nil {
				http.Error(w, "Unable to create user in database", http.StatusInternalServerError)
				return
			}
		}

	})

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		var request models.User
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		request.ConfirmToken = random.GenerateRandomString()
		s.database.User().Create(r.Context(), &request)

		// Извлекаем данные из объекта RegistrationRequest
		email := request.Email
		//password := request.Password

		err = sendmail.SendConfirmationEmail(email, request.ConfirmToken)
		if err != nil {
			return
		}

		log.Println("Email sent to:", email)
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/confirm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		token := r.URL.Query().Get("token")
		user, err := s.database.User().FindByConfirmToken(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !user.Confirmed {
			err = s.database.User().ConfirmRegistration(r.Context(), token)
			if err != nil {
				http.Error(w, "Unable to confirm registration", http.StatusBadRequest)
				return
			}

			user.Confirmed = true
			user.UpdatedAt = time.Now()

			err = s.database.User().Update(r.Context(), user)
			if err != nil {
				http.Error(w, "Unable to update user", http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, "Your registration has been confirmed. You can now log in.")
		} else {
			fmt.Fprint(w, "Your registration has already been confirmed. You can now log in.")
		}
	})

	r.Post("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		var request models.ResetPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if request.Email == "" || request.Password == "" {
			http.Error(w, "Email and new password are required", http.StatusBadRequest)
			return
		}

		resetToken := random.GenerateRandomString()
		err = s.database.User().SetPasswordResetToken(r.Context(), request.Email, resetToken)
		if err != nil {
			http.Error(w, "Unable to set password reset token", http.StatusInternalServerError)
			return
		}

		resetLink := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", resetToken)
		err = sendmail.SendPasswordResetEmail(request.Email, resetLink)
		if err != nil {
			http.Error(w, "Unable to send password reset email", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "A password reset link has been sent to your email address.")
	})

	r.Post("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		var request models.ResetPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token := r.URL.Query().Get("token")
		user, err := s.database.User().FindByPasswordResetToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid password", http.StatusBadRequest)
			return
		}

		if user == nil {
			http.Error(w, "Invalid password", http.StatusBadRequest)
			return
		}

		if request.Password == "" {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}

		err = s.database.User().UpdatePassword(r.Context(), user.Email, request.Password)
		if err != nil {
			http.Error(w, "Unable to update password", http.StatusInternalServerError)
			return
		}

		err = s.database.User().DeletePasswordResetToken(r.Context(), user.Email)
		if err != nil {
			http.Error(w, "Unable to delete password reset password", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "Your password has been successfully reset.")
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
