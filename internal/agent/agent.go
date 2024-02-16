package agent

import (
	"disCom/internal/database"
	"disCom/internal/env"
	"disCom/internal/expression"
	"disCom/internal/worker"
	"errors"
)

type Agent struct {
	ind    int
	Tasks  []chan expression.Expression
	IsFree []bool
}

func NewAgent() *Agent {
	tasks := make([]chan expression.Expression, 0)
	IsFree := make([]bool, 0)
	for i := 0; i < env.Workers; i++ {
		tasks = append(tasks, make(chan expression.Expression))

		worker.StartWorker(tasks[i])
		IsFree = append(IsFree, true)

	}
	return &Agent{Tasks: tasks, IsFree: IsFree}
}

func (ag *Agent) AddTask(expr expression.Expression) error {
	for ind, i := range ag.IsFree {
		if i {
			ag.Tasks[ind] <- expr
			ag.IsFree[ind] = false
			go func() {
				newexp := <-ag.Tasks[ind]
				newexp.Status = 0
				database.UpdateExpr(newexp)
				ag.IsFree[ind] = true
			}()
			return nil
		}
	}
	return errors.New("Все работают")
}
