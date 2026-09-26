package siblings

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
)

func run(ctx context.Context, db *sql.DB) {
	_, _ = db.Query("select 1")                      // want `\(\*sql.DB\).Query drops the context; call QueryContext`
	_, _ = db.Exec("select 1")                       // want `\(\*sql.DB\).Exec drops the context; call ExecContext`
	_ = exec.Command("true")                         // want `exec.Command drops the context; call CommandContext`
	_, _ = http.NewRequest(http.MethodGet, "/", nil) // want `http.NewRequest drops the context; call NewRequestWithContext`

	_, _ = db.QueryContext(ctx, "select 1")
	_ = exec.CommandContext(ctx, "true")
	_, _ = http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
}
