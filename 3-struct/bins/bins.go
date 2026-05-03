package bins

import (
	"errors"
	"time"
)

type Bin struct {
	id        string
	private   bool
	createdAt time.Time
	name      string
}

func NewBin(id string, private bool, name string, createdAt time.Time) (*Bin, error) {
	if id == "" {
		return nil, errors.New("id не может быть пустым")
	}
	if name == "" {
		return nil, errors.New("name не может быть пустым")
	}
	return &Bin{
		id:        id,
		private:   private,
		createdAt: createdAt,
		name:      name,
	}, nil
}

type BinList struct {
	bins []*Bin
}

func NewBinList(bins []*Bin) *BinList {
	return &BinList{
		bins: bins,
	}
}

func (b *BinList) AddBin(bin *Bin) {
	b.bins = append(b.bins, bin)
}
