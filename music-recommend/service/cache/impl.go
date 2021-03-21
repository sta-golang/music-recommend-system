package cache

import (
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	ca "github.com/sta-golang/go-lib-utils/cache"
	"github.com/sta-golang/go-lib-utils/cache/memory"
	"github.com/sta-golang/go-lib-utils/log"
	systeminfo "github.com/sta-golang/go-lib-utils/os/system_info"
	"github.com/sta-golang/music-recommend/config"
	"math"
	"runtime"
	"sync"
	"time"
)

const (
	maxMemoryScore = 0.8
	monitorScore   = 0.9
)

type cacheService struct {
	memory        ca.Cache
	priorityQueue []*set.StringSet
	setCh         chan *keyAndPriority
	pool          sync.Pool
}

type keyAndPriority struct {
	priority Priority
	key      *string
}

func InitCache() {
	PubCacheService = &cacheService{
		memory:        memory.New(config.GlobalConfig().MemoryConfig),
		priorityQueue: make([]*set.StringSet, Ten+One),
		setCh:         make(chan *keyAndPriority, 512),
		pool: sync.Pool{
			New: func() interface{} {
				return &keyAndPriority{}
			},
		},
	}
	tempService := PubCacheService.(*cacheService)
	go tempService.setAndGCRoutine()
}

var PubCacheService Cache

func (cs *cacheService) Set(key string, val interface{}, expire int, priority Priority) {
	if priority <= zero || priority > Ten {
		return
	}
	if expire == NoExpire {
		cs.memory.Set(key, val)
		chData := cs.pool.Get().(*keyAndPriority)
		chData.key = &key
		chData.priority = priority
		cs.setCh <- chData
		return
	}
	cs.memory.SetWithRemove(key, val, expire)

}

func (cs *cacheService) Get(key string) (interface{}, bool) {
	return cs.memory.Get(key)
}

func (cs *cacheService) Delete(key string) {
	var val interface{}
	var ok bool
	if val, ok = cs.memory.Get(key); !ok {
		return
	}
	cs.memory.SetWithRemove(key, val, zeroExpire)
	chData := cs.pool.Get().(*keyAndPriority)
	chData.key = &key
	chData.priority = zero
	cs.setCh <- chData
}

func (cs *cacheService) setAndGCRoutine() {
	maxCleanPriority := One
	idleTime := time.Second * time.Duration(config.GlobalConfig().MemoryConfig.GCInterval)
	ticker := time.NewTimer(idleTime)
	for {
		select {
		case priorityData := <-cs.setCh:
			cs.processPriorityData(priorityData)
		case <-ticker.C:
			if cs.cleanMemory(maxCleanPriority) {
				maxCleanPriority += One
				ticker.Reset(idleTime << 1)
				continue
			}
			maxCleanPriority = Priority(math.Max(float64(One), float64(maxCleanPriority-One)))
			ticker.Reset(idleTime)
		}
	}
}

func (cs *cacheService) cleanMemory(priority Priority) bool {
	usage, _, total := systeminfo.MemoryUsage()
	if float64(usage)/float64(total) < maxMemoryScore {
		return false
	}
	if float64(usage)/float64(total) > monitorScore {
		log.Warnf("memory 使用量报警!\n %v", systeminfo.GetSystemInfo())
		// todo 发送报警信息 如果想加可以使用邮件 短信 或者报警平台
	}
	priority = priority + One
	cnt := One
	for i := One; i <= Ten && cnt < priority; i++ {
		if cs.priorityQueue[i] == nil {
			continue
		}
		iterator := cs.priorityQueue[i].Iterator()
		cs.priorityQueue[i].Clear()
		for k := range iterator {
			cs.Delete(iterator[k])
		}
		iterator = nil
		cnt++
	}
	runtime.GC()
	return true
}

func (cs *cacheService) processPriorityData(data *keyAndPriority) {
	if data.priority == zero {
		for i := One; i <= Ten; i++ {
			if cs.priorityQueue[i] == nil {
				continue
			}
			if cs.priorityQueue[i].Contains(*data.key) {
				cs.priorityQueue[i].Remove(*data.key)
				break
			}
		}
		cs.pool.Put(data)
		return
	}
	if cs.priorityQueue[data.priority] == nil {
		cs.priorityQueue[data.priority] = set.NewStringSet()
	}
	cs.priorityQueue[data.priority].Add(*data.key)
	cs.pool.Put(data)
}
