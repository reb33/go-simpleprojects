package bins

import (
	"errors"
	"time"
)

type Bin struct {
	Id        string    `json:"id"`
	Private   bool      `json:"private"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
}

func NewBin(id string, private bool, name string, createdAt time.Time) (*Bin, error) {
	if id == "" {
		return nil, errors.New("id не может быть пустым")
	}
	if name == "" {
		return nil, errors.New("name не может быть пустым")
	}
	return &Bin{
		Id:        id,
		Private:   private,
		CreatedAt: createdAt,
		Name:      name,
	}, nil
}

type BinList struct {
	Bins []*Bin `json:"bins"`
}

func NewBinList(bins []*Bin) *BinList {
	if bins == nil {
		return &BinList{
			Bins: []*Bin{},
		}
	}
	return &BinList{
		Bins: bins,
	}
}

func (b *BinList) AddBin(bin *Bin) {
	b.Bins = append(b.Bins, bin)
}

func (b *BinList) DelBin(bin *Bin) {
	bins := []*Bin{}
	for _, b := range b.Bins {
		if b.Id != bin.Id {
			bins = append(bins, b)
		}
	}
	b.Bins = bins
}

func (b *BinList) DelBinById(id string) {
	bins := []*Bin{}
	for _, b := range b.Bins {
		if b.Id != id {
			bins = append(bins, b)
		}
	}
	b.Bins = bins
}
