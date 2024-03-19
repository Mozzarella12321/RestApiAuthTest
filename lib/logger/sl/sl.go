package sl

import (
	"log/slog"

	"github.com/google/uuid"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Inf(msg string) slog.Attr {
	return slog.Attr{
		Key:   "info",
		Value: slog.StringValue(msg),
	}
}

func Token(token uuid.UUID) slog.Attr {
	return slog.Attr{
		Key:   "token",
		Value: slog.AnyValue(token),
	}
}
