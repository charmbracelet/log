package log

import "github.com/charmbracelet/lipgloss"

// Styles is the styles for the logger.
type Styles struct {
	Timestamp  lipgloss.Style
	Caller     lipgloss.Style
	Prefix     lipgloss.Style
	Message    lipgloss.Style
	Key        lipgloss.Style
	Value      lipgloss.Style
	Separetor  lipgloss.Style
	DebugLevel lipgloss.Style
	InfoLevel  lipgloss.Style
	WarnLevel  lipgloss.Style
	ErrorLevel lipgloss.Style
}

// DefaultStyles returns the default styles for the logger.
func DefaultStyles() Styles {
	s := Styles{}

	// TimestampStyle is the style for timestamps.
	s.Timestamp = lipgloss.NewStyle() //.Faint(true)

	// CallerStyle is the style for caller.
	s.Caller = lipgloss.NewStyle().Faint(true)

	// PrefixStyle is the style for prefix.
	s.Prefix = lipgloss.NewStyle().Bold(true).Faint(true)

	// MessageStyle is the style for messages.
	s.Message = lipgloss.NewStyle()

	// KeyStyle is the style for keys.
	s.Key = lipgloss.NewStyle().Faint(true)

	// ValueStyle is the style for values.
	s.Value = lipgloss.NewStyle()

	// SeparetorStyle is the style for separetors.
	s.Separetor = lipgloss.NewStyle().Faint(true)

	// DebugLevel is the style for debug level.
	s.DebugLevel = lipgloss.NewStyle().
		SetString("DEBUG").
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.AdaptiveColor{
			Light: "62",
			Dark:  "62",
		})

	// InfoLevel is the style for info level.
	s.InfoLevel = lipgloss.NewStyle().
		SetString("INFO").
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.AdaptiveColor{
			Light: "39",
			Dark:  "86",
		})

	// WarnLevel is the style for warn level.
	s.WarnLevel = lipgloss.NewStyle().
		SetString("WARN").
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.AdaptiveColor{
			Light: "208",
			Dark:  "192",
		})

	// ErrorLevel is the style for error level.
	s.ErrorLevel = lipgloss.NewStyle().
		SetString("ERROR").
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.AdaptiveColor{
			Light: "203",
			Dark:  "204",
		})

	return s
}

// Level returns the style for the level.
func (s Styles) Level(level Level) lipgloss.Style {
	switch level {
	case DebugLevel:
		return s.DebugLevel
	case InfoLevel:
		return s.InfoLevel
	case WarnLevel:
		return s.WarnLevel
	case ErrorLevel:
		return s.ErrorLevel
	default:
		return lipgloss.NewStyle()
	}
}
