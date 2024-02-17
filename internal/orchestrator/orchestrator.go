package orchestrator

//Здесь лежит и агент и оркестратор

import (
	"disCom/internal/agent"
	"disCom/internal/expression"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	Waiting []expression.Expression //Здесь лежат выражения, которым не хватило воркеров
	Agent   agent.Agent
)

func CreateTask(expr expression.Expression) { // Создаёт задание

	var mu sync.Mutex
	mu.Lock()
	err := Agent.AddTask(expr)
	if err != nil {
		AddtoWaiting(expr)
	}

	defer mu.Unlock()

}

//Далее синхронизированные методы для добавление/удаления из слайса

func AddtoWaiting(expr expression.Expression) {
	var mu sync.Mutex
	mu.Lock()

	Waiting = append(Waiting, expr)

	defer mu.Unlock()
}

func GetFromWaiting() expression.Expression {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	expr := Waiting[0]
	Waiting = Waiting[1:]
	return expr

}

func check() { // проверяет, свободны ли воркеры, чтобы их заставить делать таски из Waiting
	go func() {
		for {
			if len(Waiting) != 0 {
				time.Sleep(100 * time.Millisecond) // Я хз, без этого не работает
				expr := GetFromWaiting()
				err := Agent.AddTask(expr)
				//fmt.Println(err)
				if err != nil {
					AddtoWaiting(expr)
				}
			}
		}
	}()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Работает успешно\nДлина Waiting:", len(Waiting))
		// http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusSeeOther) ахахаха, очень хотелось
	} else if r.Method == http.MethodPost {
		data, _ := io.ReadAll(r.Body)

		expr := expression.NewExpression("")
		json.Unmarshal(data, &expr)
		CreateTask(*expr)

	}
}

func GetInfo() []string {
	return Agent.GetAll()
}

func StartServer() { //запускает горутину с оркестратором
	Waiting = make([]expression.Expression, 0)
	check()
	go func() {
		mux1 := http.NewServeMux()
		mux1.HandleFunc("/", mainHandler)
		http.ListenAndServe(":8081", mux1)
	}()
	Agent = *agent.NewAgent()
	fmt.Println("Orchestrator is running on http://localhost:8081")

}
