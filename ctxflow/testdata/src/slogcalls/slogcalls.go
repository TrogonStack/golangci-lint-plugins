package slogcalls

import (
	"context"
	"log/slog"
)

func run(ctx context.Context, logger *slog.Logger) {
	slog.Info("started")                 // want `slog.Info drops the context; call InfoContext with the caller's ctx instead`
	slog.Debug("started")                // want `slog.Debug drops the context; call DebugContext`
	logger.Warn("slow")                  // want `\(\*slog.Logger\).Warn drops the context; call WarnContext`
	slog.Default().Error("failed")       // want `\(\*slog.Logger\).Error drops the context; call ErrorContext`
	logger.With("k", "v").Info("scoped") // want `\(\*slog.Logger\).Info drops the context`

	slog.InfoContext(ctx, "started")
	logger.WarnContext(ctx, "slow")
	logger.Log(ctx, slog.LevelInfo, "logged")
	logger.LogAttrs(ctx, slog.LevelInfo, "logged")
	_ = logger.With("k", "v")
	_ = logger.Enabled(ctx, slog.LevelInfo)
}
