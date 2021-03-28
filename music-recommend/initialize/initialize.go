package initialize

import "github.com/sta-golang/go-lib-utils/algorithm/data_structure"

type initChain struct {
	initFuncQueue *data_structure.PriorityQueue
}

var InitializationChain = initChain{initFuncQueue: data_structure.NewPriorityQueue()}

func (ic *initChain) RegisterInitialization(priority int, fn func() error) {
	ic.initFuncQueue.Push(&data_structure.Element{
		Value:    fn,
		Priority: priority,
		Index:    0,
	})
}

func (ic *initChain) Initialization() error {
	for !ic.initFuncQueue.Empty() {
		ele := ic.initFuncQueue.Pop()
		fn := ele.Value.(func() error)
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}
