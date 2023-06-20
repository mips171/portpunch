package main

import (
	"unicode"

	"fyne.io/fyne/v2/widget"
)

// Custom entry for validating an IPv4 address
// TODO: Add validation for IPv6 addresses
type ipEntry struct {
	widget.Entry
}

func newIPEntry() *ipEntry {
	entry := &ipEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (i *ipEntry) TypedRune(r rune) {
	if unicode.IsDigit(r) || r == '.' {
		i.Entry.TypedRune(r)
	}
}
