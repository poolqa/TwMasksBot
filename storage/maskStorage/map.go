package maskStorage

import (
	"../../entity/pharmacy"
	"sync"
	"time"
)

type MaskMap struct {
	m map[string]pharmacy.Pharmacy
	sync.RWMutex
}

func NewMap() *MaskMap {
	return &MaskMap{
		m: map[string]pharmacy.Pharmacy{},
	}
}

func (s *MaskMap) Add(code string, data pharmacy.Pharmacy) {
	s.Lock()
	defer s.Unlock()
	s.m[code] = data
}

func (s *MaskMap) Upd(code string, adultCount int64, childCount int64, updTime *time.Time) bool {
	s.Lock()
	defer s.Unlock()
	data, ok := s.m[code]
	if !ok {
		return false
	}
	data.AdultCount = adultCount
	data.ChildCount = childCount
	data.UpdTime = updTime
	s.m[code] = data
	return true
}

func (s *MaskMap) UpdSellRule(code string, sellRule string) bool {
	s.Lock()
	defer s.Unlock()
	data, ok := s.m[code]
	if !ok {
		return false
	}
	data.SellRule = sellRule
	s.m[code] = data
	return true
}

func (s *MaskMap) UpdSoldOut(code string, soldOutDate *time.Time) bool {
	s.Lock()
	defer s.Unlock()
	data, ok := s.m[code]
	if !ok {
		return false
	}
	data.SoldOut = 1
	data.SoldOutDate = soldOutDate
	s.m[code] = data
	return true
}

func (s *MaskMap) Get(code string) *pharmacy.Pharmacy {
	s.RLock()
	defer s.RUnlock()
	data, ok := s.m[code]
	if !ok {
		return nil
	}
	return &data
}

func (s *MaskMap) Remove(code string) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, code)
}

func (s *MaskMap) Has(code string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[code]
	return ok
}

func (s *MaskMap) Len() int {
	return len(s.m)
}

func (s *MaskMap) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[string]pharmacy.Pharmacy{}
}

func (s *MaskMap) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *MaskMap) Copy() *MaskMap {
	s.RLock()
	defer s.RUnlock()
	newSet := NewMap()
	for k, v := range s.m {
		newSet.Add(k, v)
	}
	return newSet
}

func (s *MaskMap) KeyList() []string {
	s.RLock()
	defer s.RUnlock()
	list := []string{}
	for k := range s.m {
		list = append(list, k)
	}
	return list
}

func (s *MaskMap) ValList() []pharmacy.Pharmacy {
	s.RLock()
	defer s.RUnlock()
	list := []pharmacy.Pharmacy{}
	for _, v := range s.m {
		list = append(list, v)
	}
	return list
}

//func (s *MaskMap) ListString() []string {
//	s.RLock()
//	defer s.RUnlock()
//	list := []string{}
//	for item := range s.m {
//		list = append(list, item.(string))
//	}
//	return list
//}
