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
	Bins []*Bin	`json:"bins"`
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
