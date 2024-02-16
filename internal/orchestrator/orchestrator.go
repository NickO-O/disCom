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
)

var (
	Waiting []expression.Expression
	Agent   agent.Agent
)

func CreateTask(expr expression.Expression) {
	var mu sync.Mutex
	mu.Lock()
	err := Agent.AddTask(expr)
	if err != nil {
		Waiting = append(Waiting, expr)
	}

	defer mu.Unlock()

}

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
				Agent.AddTask(GetFromWaiting())
			}
		}
	}()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusSeeOther)
	} else if r.Method == http.MethodPost {
		data, _ := io.ReadAll(r.Body)

		expr := expression.NewExpression("")
		json.Unmarshal(data, &expr)
		CreateTask(*expr)

	}
}

func StartServer() {
	check()
	go func() {
		mux1 := http.NewServeMux()
		mux1.HandleFunc("/", mainHandler)
		fmt.Println("Orchestrator is running on http://localhost:8081")
		http.ListenAndServe(":8081", mux1)
	}()
	Agent = *agent.NewAgent()
	Waiting = make([]expression.Expression, 0)
}
