package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool — интерфейс пула соединений с базой данных.
// Конкретная реализация — в пакете pgx (internal/core/repository/postgres/pool/pgx).
//
// OpTimeout() возвращает максимальное время выполнения одного запроса к БД.
type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()

	OpTimeout() time.Duration
}

// ConnPool — конкретная реализация интерфейса core_postgres_pool.ConnPool
// на базе pgxpool.Pool (пул соединений pgx).
//
// Встраивание *pgxpool.Pool даёт доступ ко всем методам pgx,
// но мы переопределяем Query/QueryRow/Exec, чтобы:
//  1. Оборачивать результаты в наши интерфейсы (Rows, Row, CommandTag)
//  2. Скрывать зависимость от pgx от слоя репозиториев
type PgxConnPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

// NewPool создаёт и проверяет пул соединений с PostgreSQL.
// Ping() при инициализации гарантирует, что БД доступна до начала работы сервера.
func NewPgxPool(
	ctx context.Context,
	config Config,
) (*ConnPool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	// Проверяем доступность БД сразу при старте.
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &PgxConnPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *PgxConnPool) OpTimeout() time.Duration {
	return p.opTimeout
}
