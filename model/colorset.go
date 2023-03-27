package model

import (
	"github.com/lucasb-eyer/go-colorful"
)

type ColorSet struct {
	colormap  map[string]*colorful.Color
	colorlist []colorful.Color
}

func NewColorSet(num int) *ColorSet {
	set := &ColorSet{
		colormap:  make(map[string]*colorful.Color),
		colorlist: colorful.FastHappyPalette(num),
	}

	return set
}

func (set *ColorSet) TryGet(name string) *colorful.Color {
	if _, ok := set.colormap[name]; !ok {
		return nil
	}
	return set.colormap[name]
}

func (set *ColorSet) Get(name string) *colorful.Color {
	if _, ok := set.colormap[name]; !ok {
		if len(set.colorlist) == 0 {
			return nil
		}
		// pop color
		color := set.colorlist[len(set.colorlist)-1]
		set.colorlist = set.colorlist[:len(set.colorlist)-1]

		set.colormap[name] = &color
		return set.colormap[name]
	}
	return set.colormap[name]
}
