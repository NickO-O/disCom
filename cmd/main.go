package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"disCom/internal/database"
	"disCom/internal/env"
	constants "disCom/internal/env"
	"disCom/internal/expression"
	"disCom/internal/logger"
	"disCom/internal/orchestrator"
	"disCom/internal/parser"
)

type expTempl struct {
	Items []string
}

type jsonSet struct {
	Plus  string `json:"plus"`
	Minus string `json:"minus"`
	Mul   string `json:"mul"`
	Div   string `json:"div"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		f, err := os.Open("frontend/main.html")
		if err != nil {
			logger.Log.Println("cannor open file main.html file: main.go func: calculateHandler")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		body, err := io.ReadAll(f)
		fmt.Fprintf(w, string(body))
	}
	if r.Method == http.MethodPost {
		body, _ := io.ReadAll(r.Body)
		expr := expression.NewExpression(string(body))

		database.WriteExpression(*expr)
		node, err := parser.ParseExpr(expr.Name)
		if err != nil {
			expr.Status = 3
			database.UpdateExpr(*expr)
			return
		}
		expr.Node = *node
		b, _ := json.Marshal(expr)
		rb := bytes.NewReader(b)
		http.Post("http://localhost:8081", "application/json", rb)
		//req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))\
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("frontend/expressions.html"))
		exprs, err := database.GetAll()
		if err != nil {
			logger.Log.Error(err.Error())
		}
		line := expTempl{}

		arr := make([]string, 0)
		for _, i := range exprs {
			arr = append(arr, i.ForTemplate())
		}
		line.Items = arr
		tmpl.Execute(w, line)

	}
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		f, err := os.Open("frontend/settings.html")
		if err != nil {
			logger.Log.Println("cannor open file main.html file: main.go func: settingsHandler")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		body, err := io.ReadAll(f)
		fmt.Fprintf(w, string(body))
	} else if r.Method == http.MethodPost {
		var set jsonSet
		var plus, minus, mul, div int
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Println("Ошибка чтения файл: main.go func: settingsHandler")
		}
		err = json.Unmarshal(body, &set)
		if err != nil {
			logger.Log.Println("Ошибка при чтении json: main.go func:settingsHandler")
		}
		plus, err = strconv.Atoi(set.Plus)
		if err == nil {
			constants.Plus = plus
		}
		minus, err = strconv.Atoi(set.Minus)
		if err == nil {
			constants.Minus = minus
		}
		mul, err = strconv.Atoi(set.Mul)
		if err == nil {
			constants.Mul = mul
		}
		div, err = strconv.Atoi(set.Div)
		if err == nil {
			constants.Div = div
		}
		env.Save()

	}
}

func computersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("frontend/computers.html"))

		line := expTempl{}

		arr := orchestrator.GetInfo()
		line.Items = arr
		tmpl.Execute(w, line)

	}
}

func main() {

	logger.Init()
	database.DeleteAll()
	orchestrator.StartServer()
	env.Init()
	filepath.Abs("/")
	mux := http.NewServeMux()
	mux.HandleFunc("/computers", computersHandler)
	mux.HandleFunc("/settings", settingsHandler)
	mux.HandleFunc("/expressions", resultHandler)
	mux.HandleFunc("/", calculateHandler)
	defer logger.End()
	fmt.Println("Server is running on http://localhost:8080	")
	http.ListenAndServe(":8080", mux)
}
