package common

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
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

	SuccessText = lipgloss.NewStyle().
			Bold(true).
			Foreground(Green)

	InfoText = lipgloss.NewStyle().
			Bold(true).
			Foreground(Blue)

	ErrorText = lipgloss.NewStyle().
			Bold(true).
			Foreground(Red)

	WarningText = lipgloss.NewStyle().
			Bold(true).
			Foreground(Orange)

	LogText1 = InfoText.Foreground(Cyan)
	LogText2 = InfoText.Foreground(Gray)
	LogText3 = InfoText.Foreground(Purple)
)
