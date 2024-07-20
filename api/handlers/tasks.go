package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/RedrikShuhartRed/finalTODO/models"
	"github.com/RedrikShuhartRed/finalTODO/task_transfer"
)

const (
	dateTimeFormat = "20060102"
)

var (
	errEmptyTitle = errors.New("error Decode request body, Task title is empty")
)

func jsonError(w http.ResponseWriter, message string) {
	resp := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetNextDate(w http.ResponseWriter, r *http.Request) {

	date := r.URL.Query().Get("date")
	now := r.URL.Query().Get("now")
	repeat := r.URL.Query().Get("repeat")

	nowTime, err := time.Parse(dateTimeFormat, now)
	if err != nil {
		log.Printf("error time.Parse now %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	result, err := task_transfer.NextDate(nowTime, date, repeat)
	if err != nil {
		log.Printf("error time.Parse now %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(result)); err != nil {
		log.Printf("error w.Write %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func AddNewTask(w http.ResponseWriter, r *http.Request) {
	var task *models.Task
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		log.Printf("error Decode request body, %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		log.Printf("error %v", errEmptyTitle)
		jsonError(w, errEmptyTitle.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task.Date, err = models.CheckDate(task)
	if err != nil {
		log.Printf("error %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lastId, err := db.AddNewTask(task)
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]int64{
		"id": lastId,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error Encode response, %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return

	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")

	tasks, err := db.GetAllTasks(search)
	if err != nil {
		log.Printf("error Scan data in Task: %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	response := map[string][]models.Task{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

	}
}

func GetTasksById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("error get id task: id == \"\"")
		jsonError(w, "Не указан индентификатор")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("error get id task, id not int: %v", err)
		jsonError(w, "Идентификатор должен быть числовым значением")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := db.GetTasksById(id)
	if err != nil {
		log.Printf("error get task from DB: %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task *models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("error Decode request body, %v", err)
		jsonError(w, "error Decode request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(task.ID) == 0 {
		log.Printf("error %v", err)
		jsonError(w, "поле id не должно быть пустым")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = strconv.Atoi(task.ID)
	if err != nil {
		log.Printf("error get id task, id not int: %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		log.Printf("error %v", errEmptyTitle)
		jsonError(w, errEmptyTitle.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task.Date, err = models.CheckDate(task)
	if err != nil {
		log.Printf("error %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rowsAffected, err := db.UpdateTask(task)
	if err != nil {
		log.Printf("error getting rows affected: %v", err)
		jsonError(w, "Ошибка сервера при обновлении задачи")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		log.Printf("no rows affected: %v", err)
		jsonError(w, "no rows affected:")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		log.Printf("Ошибка при отправке пустого JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func DoneTask(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("error get id task: id == \"\"")
		jsonError(w, "Не указан индентификатор")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("error get id task, id not int: %v", err)
		jsonError(w, "Идентификатор должен быть числовым значением")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := db.GetTasksById(id)
	if err != nil {
		log.Printf("error get task from DB: %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(task.Repeat) == 0 {

		DeleteTask(w, r)
		return
	}

	task.Date, err = models.CheckDoneDate(task)
	if err != nil {
		log.Printf("error %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rowsAffected, err := db.UpdateTask(task)
	if err != nil {
		log.Printf("error getting rows affected: %v", err)
		jsonError(w, "Ошибка сервера при обновлении задачи")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		log.Printf("no rows affected: %v", err)
		jsonError(w, "no rows affected:")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		log.Printf("Ошибка при отправке пустого JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("error get id task: id == \"\"")
		jsonError(w, "Не указан индентификатор")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("error get id task, id not int: %v", err)
		jsonError(w, "Идентификатор должен быть числовым значением")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rowsAffected, err := db.DeleteTask(id)
	if err != nil {
		log.Printf("error getting rows affected: %v", err)
		jsonError(w, "Ошибка сервера при обновлении задачи")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		log.Printf("no rows affected: %v", err)
		jsonError(w, "no rows affected:")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		log.Printf("Ошибка при отправке пустого JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// func AuthorizationGetToken(w http.ResponseWriter, r *http.Request) {
// 	password := map[string]string{}

// 	err := json.NewDecoder(r.Body).Decode(&password)
// 	if err != nil {
// 		log.Printf("error Decode request body, %v", err)
// 		jsonError(w, "error Decode request body")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	if password["password"] != environment.LoadEnvPassword() {
// 		log.Printf("error wrong password")
// 		jsonError(w, "error wrong password")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}
// 	passwordBytes := []byte(environment.LoadEnvPassword())
// 	passwordSaltBytes := []byte(environment.LoadEnvPasswordSalt())
// 	passwordBytes = append(passwordBytes, passwordSaltBytes...)
// 	hashedPasswordBytes := sha256.Sum256(passwordBytes)
// 	claims := jwt.MapClaims{
// 		"hashPass": hex.EncodeToString(hashedPasswordBytes[:]),
// 	}
// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenSaltBytes := []byte(environment.LoadEnvTokenSalt())

// 	signedToken, err := jwtToken.SignedString(tokenSaltBytes)
// 	if err != nil {
// 		log.Printf("error , %v", err)
// 		jsonError(w, err.Error())
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	res, _ := json.Marshal(&map[string]string{"token": signedToken})
// 	_, err = w.Write(res)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "error during writing data to response writer %s", err.Error())
// 		return
// 	}
// }

// func Auth(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if len(environment.LoadEnvPassword()) > 0 {
// 			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

// 			var token string
// 			cookie, err := r.Cookie("token")
// 			if err == nil {
// 				token = cookie.Value
// 			}

// 			jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
// 				return []byte(environment.LoadEnvTokenSalt()), nil
// 			})

// 			if err != nil {
// 				log.Printf("error , %v", err)
// 				jsonError(w, err.Error())
// 				w.WriteHeader(http.StatusUnauthorized)
// 				return
// 			}
// 			if !jwtToken.Valid {
// 				log.Printf("error jwt token isn't valid")
// 				jsonError(w, "jwt token isn't valid")
// 				w.WriteHeader(http.StatusUnauthorized)
// 				return
// 			}

// 			res, ok := jwtToken.Claims.(jwt.MapClaims)
// 			if !ok {
// 				log.Printf("failed to typecast to jwt.MapCalims, %v", err)
// 				jsonError(w, "failed to typecast to jwt.MapCalims")
// 				w.WriteHeader(http.StatusUnauthorized)

// 				return
// 			}

// 			hashPassRaw := res["hashPass"]
// 			hashPass, ok := hashPassRaw.(string)
// 			if !ok {
// 				log.Printf("failed to typecase password hash to string, %v", err)
// 				jsonError(w, "failed to typecase password hash to string")
// 				w.WriteHeader(http.StatusUnauthorized)

// 				return
// 			}

// 			passwordBytes := []byte(environment.LoadEnvPassword())
// 			passwordSaltBytes := []byte(environment.LoadEnvPasswordSalt())
// 			passwordBytes = append(passwordBytes, passwordSaltBytes...)
// 			hashedPasswordBytes := sha256.Sum256(passwordBytes)
// 			if hashPass != hex.EncodeToString(hashedPasswordBytes[:]) {
// 				log.Printf("token password hash doesn't match, %v", err)
// 				jsonError(w, "token password hash doesn't match")
// 				w.WriteHeader(http.StatusUnauthorized)

// 				return
// 			}
// 		}
// 		next(w, r)
// 	})
// }
