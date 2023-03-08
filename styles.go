package log

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// TimestampStyle is the style for timestamps.
	TimestampStyle = lipgloss.NewStyle()

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

	// SeparatorStyle is the style for separators.
	SeparatorStyle = lipgloss.NewStyle().Faint(true)

	// DebugLevel is the style for debug level.
	DebugLevelStyle = lipgloss.NewStyle().
			SetString(strings.ToUpper(DebugLevel.String())).
			Bold(true).
			MaxWidth(4).
			Foreground(lipgloss.AdaptiveColor{
			Light: "63",
			Dark:  "63",
		})

	// InfoLevel is the style for info level.
	InfoLevelStyle = lipgloss.NewStyle().
			SetString(strings.ToUpper(InfoLevel.String())).
			Bold(true).
			MaxWidth(4).
			Foreground(lipgloss.AdaptiveColor{
			Light: "39",
			Dark:  "86",
		})

	// WarnLevel is the style for warn level.
	WarnLevelStyle = lipgloss.NewStyle().
			SetString(strings.ToUpper(WarnLevel.String())).
			Bold(true).
			MaxWidth(4).
			Foreground(lipgloss.AdaptiveColor{
			Light: "208",
			Dark:  "192",
		})

	// ErrorLevel is the style for error level.
	ErrorLevelStyle = lipgloss.NewStyle().
			SetString(strings.ToUpper(ErrorLevel.String())).
			Bold(true).
			MaxWidth(4).
			Foreground(lipgloss.AdaptiveColor{
			Light: "203",
			Dark:  "204",
		})

	// FatalLevel is the style for error level.
	FatalLevelStyle = lipgloss.NewStyle().
			SetString(strings.ToUpper(FatalLevel.String())).
			Bold(true).
			MaxWidth(4).
			Foreground(lipgloss.AdaptiveColor{
			Light: "133",
			Dark:  "134",
		})

	// KeyStyles overrides styles for specific keys.
	KeyStyles = map[string]lipgloss.Style{}

	// ValueStyles overrides value styles for specific keys.
	ValueStyles = map[string]lipgloss.Style{}
)

// levelStyle is a helper function to get the style for a level.
func levelStyle(level Level) lipgloss.Style {
	switch level {
	case DebugLevel:
		return DebugLevelStyle
	case InfoLevel:
		return InfoLevelStyle
	case WarnLevel:
		return WarnLevelStyle
	case ErrorLevel:
		return ErrorLevelStyle
	case FatalLevel:
		return FatalLevelStyle
	default:
		return lipgloss.NewStyle()
	}
}
