package server

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"final/pkg/api"
	"final/pkg/api/services"
	"final/pkg/db"

	"github.com/golang-jwt/jwt/v5"
	_ "modernc.org/sqlite"
)

type Settings struct {
	Port   string
	PathDB string
	Logger *log.Logger
}

var setting *Settings

func initRouter(router *http.ServeMux, db *sql.DB) {
	router.HandleFunc("/api/signin", auth(func(w http.ResponseWriter, r *http.Request) {
		// Логика для обработки запроса на вход в систему
	}))
	router.Handle("/login.html", http.FileServer(http.Dir("./web")))

	router.Handle("/", http.FileServer(http.Dir("./web")))

	router.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		api.NextDayHandler(w, r, db)
	})

	router.HandleFunc("/api/task", auth(func(w http.ResponseWriter, r *http.Request) {
		api.TaskHandler(w, r, db)
	}))
	router.HandleFunc("/api/tasks", auth(func(w http.ResponseWriter, r *http.Request) {
		api.TasksHandler(w, r, db)
	}))
	router.HandleFunc("/api/task/done", auth(func(w http.ResponseWriter, r *http.Request) {
		api.DoneHandler(w, r, db)
	}))
}

func StartServer(s *Settings) *http.Server {
	setting = s
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		setting.Port = envPort
	}
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		setting.PathDB = path
	}

	err := db.Init(setting.PathDB)
	if err != nil {
		setting.Logger.Println("Ошибка при открытии базы данных:", err)
		return nil
	}

	db, err := sql.Open("sqlite", "pkg/db/scheduler.db")

	if err != nil {
		setting.Logger.Println("Ошибка при открытии базы данных:", err)
		return nil
	}

	router := http.NewServeMux()

	initRouter(router, db)
	server := http.Server{
		Addr:    ":" + setting.Port,
		Handler: router,
	}

	return &server
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := ""
		/*if password = os.Getenv("TODO_PASSWORD"); password == "" {
			password = "1"
		}*/

		if password != "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				services.Er(w, errors.New(`{'"error":"reading request body"}`), http.StatusBadRequest)
				return
			}

			i, err := tokenValid(password, r)
			if err == nil && i == http.StatusOK {
				r.Body = io.NopCloser(bytes.NewReader(body))
				next(w, r)
				return
			}
			//проверяем валидность токена
			jwtToken, i, err := passwordValid(password, body)
			if err != nil {
				services.Er(w, err, i)
				return
			} else {
				services.WriteJson(w, map[string]string{"token": jwtToken}, http.StatusOK)
				return
			}
		}
		next(w, r)
	})
}

func passwordValid(password string, body []byte) (string, int, error) {
	var passwordToken struct {
		Password string `json:"password"`
	}

	//проверям корректность данных пароля
	err := json.Unmarshal(body, &passwordToken)
	if err != nil {
		return "", http.StatusUnauthorized, err

	}
	if password != passwordToken.Password {
		return "", http.StatusUnauthorized, errors.New(`{"error":"Неверный пароль"}`)

	}
	jwtToken, err := tokenGenerate(password)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return jwtToken, http.StatusOK, nil
}

func tokenGenerate(password string) (string, error) {
	h := sha256.Sum256([]byte(password))
	hesh := hex.EncodeToString(h[:]) //hash
	var data = jwt.MapClaims{
		"token": hesh,
		"exp":   time.Now().Add(time.Minute * 5).Unix(),
	}

	j := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	jwt, err := j.SignedString([]byte(password))
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func tokenValid(password string, r *http.Request) (int, error) {
	if len(password) > 0 {
		var token string // JWT-токен из куки
		// получаем куку
		cookie, err := r.Cookie("token")
		if err == nil {
			token = cookie.Value
		} else {
			return http.StatusUnauthorized, err
		}

		// здесь код для валидации и проверки JWT-токена
		t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Здесь можно добавить проверку подписи JWT
			return []byte(password), nil
		})
		if err != nil {
			return http.StatusUnauthorized, err
		}

		if !t.Valid {
			return http.StatusUnauthorized, errors.New(`{"error":"Неверный токен"}`)
		}
	}
	return http.StatusOK, nil
}
