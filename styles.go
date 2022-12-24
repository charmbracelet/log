package log

import "github.com/charmbracelet/lipgloss"

var (
	// TimestampStyle is the style for timestamps.
	TimestampSytle = lipgloss.NewStyle().Faint(true)

	// CallerStyle is the style for caller.
	CallerStyle = lipgloss.NewStyle().Faint(true)

	// PrefixStyle is the style for prefix.
	PrefixStyle = lipgloss.NewStyle().Bold(true).Faint(true)

	// MessageStyle is the style for messages.
	MessageStyle = lipgloss.NewStyle()

	// KeyStyle is the style for keys.
	KeyStyle = lipgloss.NewStyle().Faint(true)

	// ValueStyle is the style for values.
	ValueStyle = lipgloss.NewStyle()

	// SeparetorStyle is the style for separetors.
	SeparetorStyle = lipgloss.NewStyle().Faint(true)
)

// LevelString is a map of level to string.
var LevelString = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
}

// LevelStyle is a map of level to style.
var LevelStyle = map[Level]lipgloss.Style{
	DebugLevel: lipgloss.NewStyle().
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.AdaptiveColor{
			Light: "62",
			Dark:  "62",
		}),
	InfoLevel: lipgloss.NewStyle().
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.AdaptiveColor{
			Light: "39",
			Dark:  "86",
		}),
	WarnLevel: lipgloss.NewStyle().
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.AdaptiveColor{
			Light: "208",
			Dark:  "192",
		}),
	ErrorLevel: lipgloss.NewStyle().
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.AdaptiveColor{
			Light: "203",
			Dark:  "204",
		}),
}
