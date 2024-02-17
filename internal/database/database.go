package database

import (
	"database/sql"
	"disCom/internal/expression"
	"disCom/internal/logger"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var mu sync.Mutex

func WriteExpression(exp expression.Expression) { // запысывает выражение в бд
	mu.Lock()
	defer mu.Unlock()
	db, err := sql.Open("sqlite3", "internal/database/data.sql")
	if err != nil {
		logger.Log.Error("Cannot open file data.sql file: database.go func:WriteExpression")
	}
	db.Exec("INSERT INTO Expressions (Id, Name, Status, Result) values ($1, $2, $3, $4)", exp.Id, exp.Name, exp.Status, exp.Result)

}

func ReadExpression(id int) *expression.Expression { // считывает инфу о выражении с бд
	mu.Lock()
	defer mu.Unlock()
	db, err := sql.Open("sqlite3", "internal/database/data.sql")
	if err != nil {
		logger.Log.Error("Cannot open file data.sql file: database.go func:ReadExpression")
	}
	defer db.Close()
	row := db.QueryRow("select * from Expressions where Id = $1", id)
	exp := expression.NewExpression("")
	err = row.Scan(&exp.Id, &exp.Name, &exp.Status, &exp.Result)
	if err != nil {
		return nil
	}
	return exp
}

func c() {
	db, _ := sql.Open("sqlite3", "internal/database/data.sql")
	db.Exec("CREATE TABLE Expressions(Id INTEGER PRIMARY KEY , Name TEXT, Status INTEGER, Result FLOAT);")

	defer db.Close()
}

func DeleteAll() {
	db, _ := sql.Open("sqlite3", "internal/database/data.sql")
	db.Exec("DELETE From Expressions")
	defer db.Close()
}

func UpdateExpr(expr expression.Expression) { // обновляет инфу о выражении в базе данных
	mu.Lock()
	defer mu.Unlock()
	db, _ := sql.Open("sqlite3", "internal/database/data.sql")
	if expr.Status == 0 {
		db.Exec("update Expressions set status = $1 where id = $2", expr.Status, expr.Id)
		db.Exec("update Expressions set result = $1 where id = $2", expr.Result, expr.Id)
	} else {
		db.Exec("update Expressions set status = $1 where id = $2", expr.Status, expr.Id)
	}

	defer db.Close()

}

func GetAll() ([]*expression.Expression, error) { // даёт все выражениея
	all := make([]*expression.Expression, 0)
	db, err := sql.Open("sqlite3", "internal/database/data.sql")
	if err != nil {
		logger.Log.Errorf("%s file: database.go func: GetAll", err.Error())
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Expressions")
	if err != nil {
		logger.Log.Errorf("%s file: database.go func: GetAll", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		expr := expression.NewExpression("")
		err = rows.Scan(&expr.Id, &expr.Name, &expr.Status, &expr.Result)
		if err != nil {
			logger.Log.Errorf("%s file: database.go func: GetAll", err.Error())
			return nil, err
		}
		if expr.Name != "added" {
			all = append(all, expr)
		}
	}
	return all, nil
}
