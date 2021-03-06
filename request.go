package filter

import (
	. "github.com/miraclesu/keywords-filter/keyword"
)

type Request struct {
	filter  *Filter
	content []rune
	Content string

	*Response
}

func (this *Request) Init(filter *Filter) {
	this.filter, this.content = filter, []rune(this.Content)
	this.Response = &Response{
		Threshold: filter.Threshold,
	}
}

func (this *Request) Scan() *Response {
	for i := 0; i < len(this.content); i++ {
		if node := this.trigger(this.content[i]); node != nil {
			if this.search(node, i) {
				break
			}
		}
	}
	return this.Response
}

func (this *Request) trigger(data rune) *Word {
	return this.filter.word.search(data)
}

func (this *Request) search(node *Word, index int) bool {
	//only one word
	if this.check(node, index, 0) {
		return true
	}
	for i := 1; node != nil && i < len(this.content)-index; i++ {
		c := this.content[index+i]
		tmpNode := node.search(c)
		if tmpNode == nil {
			//filter special characters
			if b, _ := this.filter.symb.search(c); b {
				continue
			}
			//is not keyword, break
			break
		}

		tmpNode.lk.RLock()
		ok := this.check(tmpNode, index, i)
		tmpNode.lk.RUnlock()
		if ok {
			return true
		}
		node = tmpNode
	}
	return false
}

func (this *Request) check(node *Word, index, i int) bool {
	if node == nil || !node.isLeaf {
		return false
	}

	this.Rate += node.rate
	this.Keywords = append(this.Keywords, Keyword{
		Rate:  node.rate,
		Index: index,
		Kind:  node.kind,
		Word:  string(this.content[index : index+i+1]),
	})
	return this.Rate >= this.filter.Threshold
}
