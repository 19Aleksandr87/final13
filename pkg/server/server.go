package server

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

// порт по умочанию
var port = "7540"
var pathDB = "./pkg/db/scheduler.db"
var TODO_PASSWORD = "" // разный стиль написания переменных начал использовать camelCase используй везде

func Init(router *http.ServeMux, db *sql.DB) {
	router.HandleFunc("/api/signin", Auth(func(w http.ResponseWriter, r *http.Request) {
		// Логика для обработки запроса на вход в систему
	}))
	router.Handle("/login.html", http.FileServer(http.Dir("./web")))

	router.Handle("/", http.FileServer(http.Dir("./web")))

	router.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		api.NextDayHandler(w, r, db)
	})

	router.HandleFunc("/api/task", Auth(func(w http.ResponseWriter, r *http.Request) {
		api.TaskHandler(w, r, db)
	}))
	router.HandleFunc("/api/tasks", Auth(func(w http.ResponseWriter, r *http.Request) {
		api.TasksHandler(w, r, db)
	}))
	router.HandleFunc("/api/task/done", Auth(func(w http.ResponseWriter, r *http.Request) {
		api.DoneHandler(w, r, db)
	}))
}

func StartServer(logger *log.Logger) *http.Server {

	// Переменные убери окружения в main формируй структуру конфига инициируй его ими
	//и передавай конфиг в функции
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		pathDB = path
	}

	err := db.Init(pathDB)
	// возвращай ошибку не фаталь тут
	if err != nil {
		log.Fatal(err)
	}
	// в переменную окружения
	db, err := sql.Open("sqlite", "pkg/db/scheduler.db")
	// возвращай ошибку
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	// инит используется только в 1 месте зачем она публичная, спрячь ее
	Init(router, db)
	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return &server
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := TODO_PASSWORD
		if pass != "" {
			body, err := io.ReadAll(r.Body)

			if err != nil {
				services.Er(w, errors.New(`{'"error":"reading request body"}`), http.StatusBadRequest)
				return
			}
			var password struct {
				Password string `json:"password"`
			}

			json.Unmarshal(body, &password) // не обработал ошибку Unmarshal
			fmt.Println(r.Cookie("token")) // пиши в лог не в консоль
			c, err := r.Cookie("token")
			fmt.Println(c) ////////////// пиши в лог не в консоль
			fmt.Println(err) // пиши в лог не в консоль

			// отдельно отработай ошибку err

			// тут отработай кейс с верным паролем  и иначе ветку с не верным так очень тяжело читать

			if TODO_PASSWORD != password.Password && err != nil {
				services.Er(w, errors.New(`{"error":"Неверный пароль"}`), http.StatusUnauthorized)
				return
			} else if TODO_PASSWORD == password.Password {
				// вынеси вычисление jwt в отдельную функцию. ты в нее пароль она тебе jwt и ошибку
				h := sha256.Sum256([]byte(pass))
				hesh := hex.EncodeToString(h[:]) //hash
				var data = jwt.MapClaims{
					"token": hesh,
					"exp":   time.Now().Add(time.Minute * 5).Unix(),
				}

				j := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
				jwt, err := j.SignedString([]byte(pass))
				if err != nil {
					services.Er(w, err, http.StatusInternalServerError)
					return
				}
				services.WriteJson(w, map[string]string{"token": jwt}, http.StatusOK)
				fmt.Printf("token: %s\n", jwt) ////////////// пиши в лог не в консоль

				return

			}
			// Всю эту логику в отдельную функцию
			// смотрим наличие пароля ты ей pass и cookie она тебе jwt
			if len(pass) > 0 {
				var token string // JWT-токен из куки
				// получаем куку
				cookie, err := r.Cookie("token")
				if err == nil {
					token = cookie.Value
				} else {
					services.Er(w, err, http.StatusInternalServerError)
					return
				}

				// здесь код для валидации и проверки JWT-токена
				t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					// Здесь можно добавить проверку подписи JWT
					return []byte(pass), nil
				})
				if err != nil {
					services.Er(w, err, http.StatusBadRequest)
					return
				}

				if !t.Valid {
					// возвращаем ошибку авторизации 401
					services.Er(w, errors.New(`{"error":"Неверный пароль"}`), http.StatusUnauthorized)
					return
				}

			}
			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		next(w, r)
	})
}
