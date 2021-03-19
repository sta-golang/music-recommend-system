package service

import (
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
	"github.com/valyala/bytebufferpool"
	"sync"
	"sync/atomic"
	"time"
)

const (
	statisticsStatusProcess = 1 //状态为处理中
	statisticsStatusReady = 0 // 状态为就绪状态
)

type StatisticsFunc struct {
	ParseFunc func(bytes []byte) // parse函数
	ProcessFunc func() // 定时处理任务
}

type statisticsService struct {
	funcArr []*StatisticsFunc
	registerMap map[string]uint8 //下标索引map
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
	ss.status = append(ss.status, statisticsStatusReady)
	return true
}

func (ss *statisticsService) Statistics(name string, buff *bytebufferpool.ByteBuffer, canDiscard bool) {
	// 最后一位为下标索引
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
				// 如果状态为信号 代表 之前要超时任务 现在完成了。检测一下队列
				if info == signalStatistics {
					for i := 0; i < ss.queue.Size(); i++ {
						queueInfo := ss.queue.Pop().(*bytebufferpool.ByteBuffer)
						queueIndex := uint8(queueInfo.B[queueInfo.Len() - 1])
						if atomic.LoadInt32(&ss.status[queueIndex]) == statisticsStatusProcess {
							ss.queue.Push(queueInfo)
							continue
						}
						fn(queueInfo, int(queueIndex))
					}
					continue
				}
				cnt++
				// 如果对应索引位置状态还是在处理 就先加入到队列里
				index := uint8(info.B[info.Len()-1])
				if atomic.LoadInt32(&ss.status[index]) == statisticsStatusProcess {
					ss.queue.Push(info)
					continue
				}
				// 处理自己的统计任务
				fn(info, int(index))
			case <- ticker.C:
				if cnt <= 0 {
					continue
				}
				// 计算最大超时时间
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
				// 计算单个的超时时间
				onceTimeOut := time.Duration(float64(maxTimeOutTime)/float64(len(ss.funcArr)))
				var runTm time.Duration
				timeOutFlag := false
				//遍历所有注册了的
				for i := 0; i < len(ss.funcArr);i++ {
					if atomic.LoadInt32(&ss.status[i]) == statisticsStatusProcess {
						continue
					}
					// 修改状态 由于目前现在是单线程 所以不会出现问题
					ss.status[i] = statisticsStatusProcess
					timer := time.NewTimer(onceTimeOut)
					var timing time.Duration
					go func(tempIndex int) {
						timing = tm.FuncTiming(func() {
							ss.funcArr[tempIndex].ProcessFunc()
						})
						runTm += timing
						processChan <- struct{}{}
					}(i)
					select {
					case <- processChan:
						onceTimeOut += onceTimeOut - timing
						ss.status[i] = statisticsStatusReady
					case <- timer.C:
						timeOutFlag = true
						timeOutChan := processChan
						go func(tempIndex int) {
							<- timeOutChan
							atomic.StoreInt32(&ss.status[tempIndex], statisticsStatusReady)
							ss.statisticsChan <- signalStatistics
						}(i)
					}
				}
				cnt = 0
				ticker.Reset(idleTime)
				if timeOutFlag {
					log.Warnf("statisticsService is process timeout maxTimeOut : %v ms",maxTimeOutTime.Milliseconds())
					continue
				}
				log.Infof("statisticsService is process run time : %v ms", runTm.Milliseconds())
			}
		}
	})
}
