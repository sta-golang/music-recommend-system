package task

import (
	"context"
	"time"

	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
)

const (
	maxRetry = 3
	waitTime = time.Second * 2
)

type TaskHandler interface {
	Handler(ctx context.Context) error
	RunNow() bool
	Name() string
	TaskTimeInterval() time.Duration
}

func InitTask() {
	RegisterTask(NewCreatorHotScoreTask())
	Handler()
}

var tasks []TaskHandler

func RegisterTask(task TaskHandler) {
	tasks = append(tasks, task)
}

func TimeHandler(ctx context.Context, name string, handler func(ctx context.Context) error) error {
	var err error
	timing := tm.FuncTiming(func() {
		err = handler(ctx)
	})
	log.Debugf("timing task name : %s run time %v 毫秒", name, timing.Milliseconds())
	return err
}

func Handler() {
	ctx := context.Background()
	for _, task := range tasks {
		go func(curTask TaskHandler) {
			var taskErr error
			if curTask.RunNow() {
				if taskErr = TimeHandler(ctx, curTask.Name(), curTask.Handler); taskErr != nil {
					log.Error(taskErr)
				}
				for {
					time.Sleep(curTask.TaskTimeInterval())
					if taskErr = TimeHandler(ctx, curTask.Name(), curTask.Handler); taskErr != nil {
						log.Error(taskErr)
					}
				}
			}
		}(task)
	}
}
