package common

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

var (
	Base16     = huh.ThemeBase16()
	Blue       = lipgloss.Color("#5e81ac")
	Red        = lipgloss.Color("#bf616a")
	Green      = lipgloss.Color("#a3be8c")
	Yellow     = lipgloss.Color("#ebcb8b")
	Orange     = lipgloss.Color("#d08770")
	Purple     = lipgloss.Color("#b48ead")
	Cyan       = lipgloss.Color("#88c0d0")
	White      = lipgloss.Color("#e5e9f0")
	LightGray  = lipgloss.Color("#d8dee9")
	Gray       = lipgloss.Color("#7B869B")
	DarkGray   = lipgloss.Color("#3b4252")
	DarkerGray = lipgloss.Color("#2e3440")

	SuccessText = style(Green, true)
	InfoText    = style(Blue, true)
	ErrorText   = style(Red, true)
	WarningText = style(Orange, true)
	LogText1    = style(Cyan, false)
	LogText2    = style(Gray, false)
	LogText3    = style(Purple, false)
)

func style(c lipgloss.TerminalColor, bold bool) lipgloss.Style {
	return lipgloss.NewStyle().Bold(bold).Foreground(c)
}

func TeaEnumerator(_ list.Items, i int) string {
	if i < 9 {
		return fmt.Sprintf("0%d. ", i+1)
	} else {
		return fmt.Sprintf("%d. ", i+1)
	}
}

func TeaList(items ...any) *list.List {
	return list.New(items).
		Enumerator(TeaEnumerator).
		EnumeratorStyle(WarningText)
}
