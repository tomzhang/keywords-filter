package filter

import (
	"errors"

	. "github.com/miraclesu/keywords-filter/keyword"
	"github.com/miraclesu/keywords-filter/loader"
)

var (
	InvalidFilter = errors.New("Invalid filter")
)

type Filter struct {
	word      *Word    // keyword
	symb      *symbols // special symbols
	Threshold int
}

func New(threshold int, load loader.Loader) (f *Filter, err error) {
	kws, sbs, err := load.Load()
	if err != nil {
		return
	}

	f = &Filter{
		word:      new(Word),
		symb:      new(symbols),
		Threshold: threshold,
	}

	f.AddWords(kws)
	f.AddSymbs(sbs)
	return
}

func (this *Filter) Filter(content string) (b bool, err error) {
	if this.word == nil || this.symb == nil {
		err = InvalidFilter
		return
	}
	return NewRequest(this, content).scan(), nil
}

func (this *Filter) AddWord(w *Keyword) {
	if this.word == nil {
		this.word = new(Word)
	}
	this.word.addWord(w)
}

func (this *Filter) RemoveWord(w *Keyword) {
	if this.word == nil {
		return
	}
	this.word.removeWord(w)
}

func (this *Filter) AddWords(kws []*Keyword) {
	for i, count := 0, len(kws); i < count; i++ {
		this.AddWord(kws[i])
	}
}

func (this *Filter) AddSymb(w *Keyword) {
	if this.symb == nil {
		this.symb = new(symbols)
	}
	for _, v := range w.Word {
		this.symb.add(v)
	}
}

func (this *Filter) RemoveSymb(w *Keyword) {
	if this.word == nil {
		return
	}
	for _, v := range w.Word {
		this.symb.remove(v)
	}
}

func (this *Filter) AddSymbs(sbs []*Keyword) {
	for i, count := 0, len(sbs); i < count; i++ {
		this.AddSymb(sbs[i])
	}
}
