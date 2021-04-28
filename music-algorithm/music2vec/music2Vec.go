package music2vec

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-algorithm/common"
	"github.com/ynqa/wego/pkg/model"
	"github.com/ynqa/wego/pkg/model/modelutil/matrix"
	"github.com/ynqa/wego/pkg/model/modelutil/vector"
	"github.com/ynqa/wego/pkg/model/word2vec"
)

type Model struct {
	idTable      map[int]int
	disFunc      common.Distance
	model        model.Model
	dataStr      *string
	loadFlag     bool
	similarMusic []*similar
	vector       *matrix.Matrix
}

type similar struct {
	edges []edge
}

type edge struct {
	distance float64
	musicID  int
}

func NewDefaultModel() (*Model, error) {
	return newModel(common.CosSimi,
		word2vec.Window(7),
		word2vec.Model(word2vec.Cbow),
		word2vec.Optimizer(word2vec.NegativeSampling),
		word2vec.NegativeSampleSize(3),
		word2vec.Verbose(),
		word2vec.Dim(150),
		word2vec.Iter(200),
	)
}

func NewModel(opts ...word2vec.ModelOption) (*Model, error) {
	return newModel(common.CosSimi, opts...)
}

func NewModelWithDisFunc(disFunc common.Distance, opts ...word2vec.ModelOption) (*Model, error) {
	return newModel(disFunc, opts...)
}

func newModel(disFunc common.Distance, opts ...word2vec.ModelOption) (*Model, error) {
	modelTmp, err := word2vec.New(opts...)
	if err != nil {
		return nil, err
	}
	return &Model{
		disFunc: disFunc,
		model:   modelTmp,
	}, nil
}

func (m *Model) LoadData(bys []byte) error {
	dataStr := str.BytesToString(bys)
	dataLine := strings.Split(dataStr, "\n")
	if len(dataLine) <= 0 {
		log.ConsoleLogger.Warn("loadData len is 0")
		return nil
	}
	if m.idTable == nil {
		m.idTable = make(map[int]int)
	}
	cnt := 0
	for i := range dataLine {
		musicIDs := strings.Split(dataLine[i], " ")
		for j := range musicIDs {
			if len(musicIDs[j]) <= 0 {
				continue
			}
			musicID, err := strconv.Atoi(musicIDs[j])
			if err != nil {
				log.ConsoleLogger.Error(err)
				continue
			}
			if _, ok := m.idTable[musicID]; !ok {
				m.idTable[musicID] = cnt
				cnt++
			}
		}
	}
	m.dataStr = &dataStr
	m.loadFlag = true
	return nil
}

func (m *Model) Train() error {
	if !m.loadFlag {
		return fmt.Errorf("Please LoadData")
	}
	if m.dataStr == nil {
		log.ConsoleLogger.Warn("use empty data")
		return nil
	}
	reader := strings.NewReader(*m.dataStr)
	if err := m.model.Train(reader); err != nil {
		return err
	}
	m.vector = m.model.WordVector(vector.Single)
	return nil
}

func (m *Model) SimilarMusic(musicID int) *similar {
	index := -1
	if val, ok := m.idTable[musicID]; ok {
		index = val
	}
	if index == -1 {
		return nil
	}
	if m.similarMusic == nil {
		m.similarMusic = make([]*similar, len(m.idTable))
	}
	if m.similarMusic[index] != nil {
		return m.similarMusic[index]
	}
	m.similarMusic[index] = &similar{
		edges: m.computeEdges(index),
	}
	return m.similarMusic[index]
}

func (m *Model) computeEdges(index int) []edge {
	ret := make([]edge, 0, len(m.idTable)-1)
	vec := m.vector.Slice(index)
	for key, val := range m.idTable {
		if val == index {
			continue
		}
		curVec := m.vector.Slice(val)
		ret = append(ret, edge{
			musicID:  key,
			distance: m.disFunc(vec, curVec),
		})
	}
	sort.Slice(ret, func(i, j int) bool {
		if ret[i].distance > ret[j].distance {
			return true
		}
		return false
	})
	return ret
}

func (s *similar) TopN(n int) []int {
	length := common.MinInt(n, len(s.edges))
	ret := make([]int, length)
	for i := 0; i < length; i++ {
		ret[i] = s.edges[i].musicID
	}
	return ret
}
