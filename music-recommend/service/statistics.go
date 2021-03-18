package service

import (
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure"
	tm "github.com/sta-golang/go-lib-utils/time"
	"github.com/valyala/bytebufferpool"
	"sync"
	"sync/atomic"
	"time"
)

const (
	statisticsStatusProcess = 1
	statisticsStatusReady = 0
)

type StatisticsFunc struct {
	ParseFunc func(bytes []byte)
	ProcessFunc func()
}

type statisticsService struct {
	funcArr []*StatisticsFunc
	registerMap map[string]uint8
	status []int32
	queue *data_structure.Queue
	statisticsChan chan *bytebufferpool.ByteBuffer
}

var PubStatisticsService = &statisticsService{
	funcArr:        make([]*StatisticsFunc, 0),
	registerMap: map[string]uint8{},
	statisticsChan: make(chan *bytebufferpool.ByteBuffer, 10240),
	queue: data_structure.NewQueue(),
}
var onceStatisticsService sync.Once
var lock = sync.Mutex{}

var signalStatistics = &bytebufferpool.ByteBuffer{}

// Register 注册函数 应该在一开始就注册好 不关注性能
func (ss *statisticsService) Register(name string, fn *StatisticsFunc) bool {
	lock.Lock()
	defer lock.Unlock()
	length := len(ss.funcArr)
	if length >= 0xff {
		return false
	}
	ss.registerMap[name] = uint8(length)
	ss.funcArr = append(ss.funcArr, fn)
	return true
}

func (ss *statisticsService) Statistics(name string, buff *bytebufferpool.ByteBuffer, canDiscard bool) {
	buff.B = append(buff.B, ss.registerMap[name])
	select {
	case ss.statisticsChan <- buff:
	default:
		if !canDiscard {
			go func() {
				ss.statisticsChan <- buff
			}()
		}
	}
}

func (ss *statisticsService) Run() {
	onceStatisticsService.Do(func() {
		if len(ss.funcArr) <= 0 {
			return
		}
		idleTime := time.Second * 10
		cnt := 0
		ticker := time.NewTimer(idleTime)
		fn := func(info *bytebufferpool.ByteBuffer, index int) {
			ss.funcArr[index].ParseFunc(info.Bytes()[:info.Len() -1])
			bytebufferpool.Put(info)
		}
		processChan := make(chan struct{}, 1)
		for  {
			select {
			case info := <- ss.statisticsChan:
				if info == signalStatistics {
					for i := 0; i < ss.queue.Size(); i++ {
						queueInfo := ss.queue.Pop().(*bytebufferpool.ByteBuffer)
						queueIndex := queueInfo.B[queueInfo.Len() - 1]
						if atomic.LoadInt32(&ss.status[queueIndex]) == statisticsStatusProcess {
							ss.queue.Push(queueInfo)
							continue
						}
						fn(queueInfo, int(queueIndex))
					}
					continue
				}
				cnt++
				index := info.B[info.Len()-1]
				if atomic.LoadInt32(&ss.status[index]) == statisticsStatusProcess {
					ss.queue.Push(info)
					continue
				}
				fn(info, int(index))
			case <- ticker.C:
				if cnt <= 0 {
					continue
				}
				maxTimeOutTime := idleTime >> 1
				if cnt > len(ss.statisticsChan) / 2 {
					maxTimeOutTime = time.Duration(float64(maxTimeOutTime) * (float64(len(ss.statisticsChan))/float64(cnt)))
					if maxTimeOutTime < idleTime >> 2 {
						maxTimeOutTime = idleTime >> 2
					}
					if maxTimeOutTime > (idleTime << 1) - (idleTime >> 1) {
						maxTimeOutTime = (idleTime << 1) - (idleTime >> 1)
					}
				}
				onceTimeOut := time.Duration(float64(maxTimeOutTime)/float64(len(ss.funcArr)))
				var runTm time.Duration
				for i := 0; i < len(ss.funcArr);i++ {
					ss.status[i] = statisticsStatusProcess
					timer := time.NewTimer(onceTimeOut)
					go func(tempIndex int) {
						timing := tm.FuncTiming(func() {
							ss.funcArr[tempIndex].ProcessFunc()
							atomic.StoreInt32(&ss.status[tempIndex], statisticsStatusReady)
						})
						if time.Duration(float64(timing) * 1.2 ) < onceTimeOut {
							onceTimeOut += onceTimeOut - timing
						}
						runTm += timing
						processChan <- struct{}{}
					}(i)
					select {
					case <- processChan:
					case <- timer.C:
					}
				}
				cnt = 0
				ticker.Reset(idleTime)
			}
		}
	})
}
