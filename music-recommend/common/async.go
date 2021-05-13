package common

import (
	"github.com/sta-golang/go-lib-utils/async/dag"
	"github.com/sta-golang/go-lib-utils/async/group"
)

var SingleRunGroup = group.SingleExecGroup{}

func InitDag() {
	dag.Config().SetPool()
}
