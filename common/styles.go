package common

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	Base16 = huh.ThemeBase16()
	Green  = lipgloss.Color("#a3be8c")
	Blue   = lipgloss.Color("#5e81ac")
	Red    = lipgloss.Color("#bf616a")

	SuccessText = lipgloss.NewStyle().
			Bold(true).
			Foreground(Green)

	InfoText = lipgloss.NewStyle().
			Bold(true).
			Foreground(Blue)

	ErrorText = lipgloss.NewStyle().
			Bold(true).
			Foreground(Red)
)
