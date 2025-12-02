// Licensed under Elastic License 2.0
// See LICENSE.txt for details

package log

import (
	"log/slog"
)

var (
	MinLevelDebug = MinLevel{slog.LevelDebug}
	MinLevelInfo  = MinLevel{slog.LevelInfo}
	MinLevelWarn  = MinLevel{slog.LevelWarn}
	MinLevelError = MinLevel{slog.LevelError}
	MinLevelFatal = MinLevel{LevelFatal}
)

type MinLevel struct {
	level slog.Level
}

func (ml MinLevel) Level() slog.Level {
	return ml.level
}

const LevelFatal = slog.Level(12)

func resolveLogLevel(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		if level == LevelFatal {
			a.Value = slog.StringValue("FATAL")
		}
	}

	return a
}
