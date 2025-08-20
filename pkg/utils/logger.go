package utils

import (
	"context"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"runtime/debug"
	"sync/atomic"
	"time"
)

type Logger struct {
	*slog.Logger
}

type stackHandler struct {
	*slog.TextHandler

	disable atomic.Bool
}

func (h *stackHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.disable.Load() {
		return nil
	}

	if r.Level >= slog.LevelError {
		r.AddAttrs(slog.String("stack", string(debug.Stack())))
	}

	return h.TextHandler.Handle(ctx, r)
}

func NewLogger(level slog.Level) *Logger {
	filename := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02 15"))

	rotator := &lumberjack.Logger{
		Filename:   filename, // имя базового файла
		MaxSize:    10,       // мегабайты до ротации
		MaxBackups: 5,        // сколько старых файлов хранить
		MaxAge:     7,        // сколько дней хранить логи
		Compress:   true,     // сжимать ли старые логи в .gz
	}

	// Можно писать логи и в файл и в stdout одновременно
	multi := slog.NewTextHandler(rotator, &slog.HandlerOptions{Level: level}) // io.MultiWriter(os.Stdout, rotator)
	return &Logger{Logger: slog.New(&stackHandler{TextHandler: multi})}
}

func (l *Logger) Disable() {
	if h, ok := l.Logger.Handler().(*stackHandler); ok {
		h.disable.Store(true)
	}
}
