package cache

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	ca "github.com/sta-golang/go-lib-utils/cache"
	"github.com/sta-golang/go-lib-utils/cache/memory"
	"github.com/sta-golang/go-lib-utils/log"
	systeminfo "github.com/sta-golang/go-lib-utils/os/system_info"
	"github.com/sta-golang/music-recommend/config"
	"github.com/sta-golang/music-recommend/service/email"
	"math"
	"runtime"
	"sync"
	"time"
)

const (
	maxMemoryScore = 0.8
	monitorScore   = 0.9

	memoryWarnCacheKey = "memory_warn"
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
	if priority != Forever && (priority <= zero || priority > Ten) {
		return
	}
	if priority != Forever {
		chData := cs.pool.Get().(*keyAndPriority)
		chData.key = &key
		chData.priority = priority
		cs.setCh <- chData
	}
	if expire <= NoExpire {
		cs.memory.Set(key, val)
		return
	}
	cs.memory.SetWithRemove(key, val, expire)
}

func (cs *cacheService) Get(key string) (interface{}, bool) {
	return cs.memory.Get(key)
}

func (cs *cacheService) Delete(key string) {
	cs.doDelete(key, true)
}

func (cs *cacheService) doDelete(key string, isRef bool) {
	var val interface{}
	var ok bool
	if val, ok = cs.memory.Get(key); !ok {
		return
	}
	cs.memory.SetWithRemove(key, val, zeroExpire)
	if !isRef {
		return
	}
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
	usage, _, total := systeminfo.MemoryUsage(false)
	if float64(usage)/float64(total) < maxMemoryScore {
		return false
	}
	if float64(usage)/float64(total) > monitorScore {
		info := systeminfo.GetSystemInfo()
		log.Warnf("memory 使用量报警!\n %v", info)
		warnNum := 1
		if val, ok := cs.Get(memoryWarnCacheKey); ok {
			warnNum = val.(int) + 1
		}
		cs.Set(memoryWarnCacheKey, warnNum, Hour, Ten)
		if warnNum == 1 || warnNum == 2 || warnNum == 4 || warnNum == 8 || warnNum == 16 || warnNum == 32 || warnNum == 64 || warnNum == 128 || warnNum == 256 {

			go func() {
				err := email.PubEmailService.SendEmail("内存使用告警",
					fmt.Sprintf("内存使用告警 次数 : %d \n%v", warnNum, info),
					config.GlobalConfig().EmailConfig.Email)
				if err != nil {
					log.Error(err)
				}
			}()
		}
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
			if iterator[k] == memoryWarnCacheKey {
				continue
			}
			cs.doDelete(iterator[k], false)
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
