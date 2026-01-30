package migrations

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func CreateUserNotificationTable(m *migrate.Migrations) {
	m.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE user_notifications (
						id UUID PRIMARY KEY DEFAULT uuidv7(),
						user_id UUID not null,
						data JSONB,
						status INTEGER NOT NULL DEFAULT 1,
						updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
						created_at TIMESTAMP NOT NULL DEFAULT NOW()
					);
					CREATE INDEX idx_u_n_user_id ON user_notifications USING BTREE(user_id);
					CREATE INDEX idx_u_n_created_at ON user_notifications USING BTREE(created_at);
					CREATE INDEX idx_u_n_status ON user_notifications USING HASH(status);
			`)
				return err
			})
		},
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS user_notifications`)
			return err
		},
	)
}
