-- Runs once on first ClickHouse startup (mounted at
-- /docker-entrypoint-initdb.d/). Schema matches the fields the app's
-- structured slog output actually emits (internal/platform/middleware/logging.go):
-- request_id/method/path/status/duration_ms for "request completed" lines,
-- plus a few fields that only appear on other log lines (panic, startup
-- messages) -- `raw` keeps the untouched JSON line so nothing is lost even
-- for log shapes this schema doesn't have a dedicated column for.
CREATE DATABASE IF NOT EXISTS logs;

CREATE TABLE IF NOT EXISTS logs.app_logs
(
    ts          DateTime64(3) DEFAULT now64(3),
    level       LowCardinality(String),
    msg         String,
    request_id  String,
    method      LowCardinality(String),
    path        String,
    status      Int32,
    duration_ms Int64,
    raw         String
)
ENGINE = MergeTree()
ORDER BY ts
TTL toDateTime(ts) + INTERVAL 7 DAY;
