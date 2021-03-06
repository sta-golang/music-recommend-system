package utils

import (
	"strconv"

	"github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"
	"github.com/sta-golang/music-recommend/model"
)

type Filter func(*model.Item) bool

type filterChain struct {
	filters []Filter
}

func NewFilterChain(filters ...Filter) *filterChain {
	ret := &filterChain{}
	if len(filters) <= 0 {
		return ret
	}
	ret.filters = append(ret.filters, filters...)
	return ret
}

func NewExistFilterChain(existSet, userRead *set.StringSet) *filterChain {
	return NewFilterChain(ItemBaseFilter(), ItemNotExist(existSet), ItemMusicStatusHasMusicUrl(),
		ItemUserRead(userRead))
}

func (fc *filterChain) DoFilter(item *model.Item) bool {
	if len(fc.filters) <= 0 {
		return true
	}
	for _, filter := range fc.filters {
		if !filter(item) {
			return false
		}
	}
	return true
}

func ItemUserRead(userRead *set.StringSet) Filter {
	return func(item *model.Item) bool {
		if item == nil {
			return false
		}
		if userRead == nil {
			return true
		}
		if userRead.Contains(strconv.Itoa(item.Music.ID)) {
			return false
		}
		return true
	}
}

func ItemMusicStatusHasMusicUrl() Filter {
	return func(item *model.Item) bool {
		if item == nil {
			return false
		}
		if item.Music.Status != model.MusicHasMusicUrlStatus {
			return false
		}
		return true
	}
}

func ItemNotExist(existSet *set.StringSet) Filter {
	return func(item *model.Item) bool {
		if item == nil {
			return false
		}
		if item.Music.ID <= 0 {
			return false
		}
		if existSet.Contains(strconv.Itoa(item.Music.ID)) {
			return false
		}
		return true
	}
}

func ItemBaseFilter() Filter {
	return func(item *model.Item) bool {
		if item == nil {
			return false
		}
		if item.Music.ID <= 0 {
			return false
		}
		return true
	}
}
