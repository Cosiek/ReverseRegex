package reverse_regex

type ReverseRegex struct {
	mainGroup *runeGroup
}

type indexedElement interface {
	getIndex() int16
}

type runeGroup struct {
	strings   []*stringWrapper
	subGroups []*runeGroup
	isOpen    bool
	idx       int16
}

func (rg *runeGroup) getIndex() int16 { return rg.idx }

type stringWrapper struct {
	str string
	idx int16
}

func (sw *stringWrapper) getIndex() int16 { return sw.idx }

func (rg *runeGroup) addRune(rune_ string) {
	if len(rg.subGroups) > 0 {
		lastSG := rg.subGroups[len(rg.subGroups)-1]
		if lastSG.isOpen {
			lastSG.addRune(rune_)
		} else {
			rg.addRuneToString(rune_)
		}
	} else {
		rg.addRuneToString(rune_)
	}
}

func (rg *runeGroup) addRuneToString(rune_ string) {
	// get or create last string wrapper
	var lastString *stringWrapper
	var idx int16
	if len(rg.strings) > 0 {
		lastString = rg.strings[len(rg.strings)-1]
		if len(rg.subGroups) > 0 {
			lastSG := rg.subGroups[len(rg.subGroups)-1]
			if lastString.idx < lastSG.idx {
				rg.strings = append(rg.strings, newStringWrapper(lastSG.idx+1))
				lastString = rg.strings[len(rg.strings)-1]
			}
		}
	} else {
		idx = rg.getNextIdx()
		rg.strings = append(rg.strings, newStringWrapper(idx))
		lastString = rg.strings[len(rg.strings)-1]
	}
	// add rune
	lastString.str += rune_
}

func (rg *runeGroup) addSubGroup(initialRune string) {
	// get parent
	last, _, isSubGroup := rg.getLastElement()
	if isSubGroup {
		parent := last.(*runeGroup)
		if parent.isOpen {
			parent.addSubGroup(initialRune)
			return
		}
	}
	// create subgroup
	idx := rg.getNextIdx()
	subGroup := newRuneGroup(idx)
	subGroup.addRuneToString(initialRune)
	rg.subGroups = append(rg.subGroups, subGroup)
}

func (rg *runeGroup) getLastElement() (indexedElement, int16, bool) {
	var last indexedElement
	var idx int16 = 0
	var isSubGroup bool
	if len(rg.strings) > 0 {
		isSubGroup = false
		last = rg.strings[len(rg.strings)-1]
		idx = last.getIndex()
	}
	if len(rg.subGroups) > 0 {
		if rg.subGroups[len(rg.subGroups)-1].getIndex() >= idx {
			isSubGroup = true
			last = rg.subGroups[len(rg.subGroups)-1]
			idx = last.getIndex()
		}
	}
	return last, idx, isSubGroup
}

func (rg *runeGroup) getNextIdx() int16 {
	_, idx, _ := rg.getLastElement()
	return idx + 1
}

// closes the last open subGroup if there is one that is still open,
// or this one if there aren't any
func (rg *runeGroup) closeGroup(remainingRune string) {
	if len(rg.subGroups) > 0 {
		lastSG := rg.subGroups[len(rg.subGroups)-1]
		if lastSG.isOpen {
			lastSG.closeGroup(remainingRune)
			return
		}
	}
	// if we're here, then there ware no open (or any) subGroups
	// wer'e closing this one
	rg.addRuneToString(remainingRune)
	rg.isOpen = false
}

func (rRx *ReverseRegex) GetReversedString(inputs ...string) string {
	return rRx.mainGroup.getReversedString(inputs...)
}

func (rg *runeGroup) getReversedString(inputs ...string) string {
	rS := ""
	var idx int16 = 1
	var stringsIdx int16 = 0
	var subGroupsIdx int16 = 0
	stringsLen := int16(len(rg.strings))
	subGroupsLen := int16(len(rg.subGroups))
	for {
		// gather strings
		for sIdx := stringsIdx; sIdx < stringsLen; sIdx++ {
			if idx == rg.strings[sIdx].idx {
				rS += rg.strings[sIdx].str
				stringsIdx++
				idx++
			} else if idx < rg.strings[sIdx].idx {
				break
			}
		}
		// apply subGroups
		for sgIdx := subGroupsIdx; sgIdx < subGroupsLen; sgIdx++ {
			if idx == rg.subGroups[sgIdx].idx {
				rS += inputs[subGroupsIdx]
				subGroupsIdx++
				idx++
			} else if idx < rg.subGroups[sgIdx].idx {
				break
			}
		}
		// is this the end?
		if idx > stringsLen+subGroupsLen {
			break
		}
	}
	return rS
}

func newRuneGroup(idx int16) *runeGroup {
	return &runeGroup{
		strings:   make([]*stringWrapper, 0),
		subGroups: make([]*runeGroup, 0),
		isOpen:    true,
		idx:       idx,
	}
}

func newStringWrapper(idx int16) *stringWrapper {
	return &stringWrapper{
		str: "",
		idx: idx,
	}
}

func NewReverseRegex(pattern string) *ReverseRegex {
	mainGroup := newRuneGroup(0)
	for idx := 0; idx < len(pattern); idx++ {
		letter := string(pattern[idx])
		if letter == `\` {
			// an escape char
			mainGroup.addRune(`\`)
			idx++
			if idx < len(pattern) {
				mainGroup.addRune(string(pattern[idx]))
			}
		} else if letter == "(" {
			// open new group
			mainGroup.addSubGroup("(")
		} else if letter == ")" {
			// close last group
			mainGroup.closeGroup(")")
		} else {
			// something else - simply add to group
			mainGroup.addRune(string(pattern[idx]))
		}
	}
	mainGroup.closeGroup("")
	return &ReverseRegex{mainGroup}
}
