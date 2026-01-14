package storage

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.ExecContext(ctx, `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT uuid_v7(),
   first_name VARCHAR(255) NOT NULL,
   second_name VARCHAR(255) NOT NULL,
   email VARCHAR(255) NOT NULL UNIQUE,
   phone_number VARCHAR(255) NOT NULL UNIQUE,
   status INTEGER NOT NULL DEFAULT 0,
   password VARCHAR(255) NOT NULL,
   roles TEXT[] NOT NULL DEFAULT '{}',
   address JSONB,
   confirmed_at TIMESTAMP,
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_status ON users(status);
CREATE UNIQUE INDEX idx_users_email ON users(email) NULLS DISTINCT;
CREATE UNIQUE INDEX idx_users_phone_number ON users(phone_number) NULLS DISTINCT;
CREATE INDEX idx_users_roles ON users USING GIN(roles);
CREATE INDEX idx_users_created_at ON users(created_at);
`)
			return err
		},
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.ExecContext(ctx, `DROP TABLE users`)
			return err
		},
	)
}
