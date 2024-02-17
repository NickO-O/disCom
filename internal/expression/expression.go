package expression

import (
	"disCom/internal/parser"
	"fmt"

	"github.com/google/uuid"
)

type Expression struct {
	Name   string // Изначальное значение выражения
	Status int    // Статус выражения: 0 если посчиталось, 1 если считается, 2 если ждёт вычисления, 3 если выражение невалидно
	Id     int
	Result float64 // результат выражения, если посчиталось
	Node   parser.Node
}

func NewExpression(Name string) *Expression {

	return &Expression{Name: Name, Status: 2, Id: int(uuid.New().ID())}
}
func (exp *Expression) ForTemplate() string {
	var stat string
	if exp.Status == 0 {
		stat = "Выражение посчиталось, результат:"
	} else if exp.Status == 1 {
		stat = "Выражение считается"
	} else if exp.Status == 2 {
		stat = "Выражение ожидает рассчёта"
	} else if exp.Status == 3 {
		stat = "Выражение невалидно"
	}
	if exp.Status == 0 {
		return fmt.Sprintf("id: %d, %s %s %.4f", exp.Id, exp.Name, stat, exp.Result)
	} else {
		return fmt.Sprintf("id: %d, %s %s", exp.Id, exp.Name, stat)
	}

}
