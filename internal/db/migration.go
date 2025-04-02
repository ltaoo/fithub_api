package db

import (
	"database/sql"
	"fmt"
	"myapi/config"
	"myapi/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrator 处理数据库迁移
type Migrator struct {
	config *config.Config
	logger *logger.Logger
}

// NewMigrator 创建新的迁移器
func NewMigrator(cfg *config.Config, logger *logger.Logger) *Migrator {
	return &Migrator{
		config: cfg,
		logger: logger,
	}
}

// MigrateUp 运行所有向上迁移
func (m *Migrator) MigrateUp() error {
	migrator, err := m.createMigrator()
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, dbErr := migrator.Close()
		if sourceErr != nil {
			m.logger.Error("Error closing migration source", sourceErr)
		}
		if dbErr != nil {
			m.logger.Error("Error closing migration database", dbErr)
		}
	}()

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.logger.Info("Database migrations completed successfully")
	return nil
}

// MigrateDown 回滚所有迁移
func (m *Migrator) MigrateDown() error {
	migrator, err := m.createMigrator()
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, dbErr := migrator.Close()
		if sourceErr != nil {
			m.logger.Error("Error closing migration source", sourceErr)
		}
		if dbErr != nil {
			m.logger.Error("Error closing migration database", dbErr)
		}
	}()

	if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	m.logger.Info("Database rollback completed successfully")
	return nil
}

// MigrateTo 迁移到特定版本
func (m *Migrator) MigrateTo(version uint) error {
	migrator, err := m.createMigrator()
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, dbErr := migrator.Close()
		if sourceErr != nil {
			m.logger.Error("Error closing migration source", sourceErr)
		}
		if dbErr != nil {
			m.logger.Error("Error closing migration database", dbErr)
		}
	}()

	if err := migrator.Migrate(version); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	m.logger.Info(fmt.Sprintf("Database migrated to version %d successfully", version))
	return nil
}

// createMigrator 创建迁移实例
func (m *Migrator) createMigrator() (*migrate.Migrate, error) {
	var db *sql.DB
	var driver string
	var instance database.Driver
	var err error

	// 根据数据库类型创建连接
	switch m.config.DBType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true",
			m.config.DBUser, m.config.DBPassword, m.config.DBHost, m.config.DBPort, m.config.DBName)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open mysql connection: %w", err)
		}
		instance, err = mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create mysql instance: %w", err)
		}
		driver = "mysql"
	case "postgres":
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			m.config.DBUser, m.config.DBPassword, m.config.DBHost, m.config.DBPort, m.config.DBName)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open postgres connection: %w", err)
		}
		instance, err = postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create postgres instance: %w", err)
		}
		driver = "postgres"
	case "sqlite":
		db, err = sql.Open("sqlite3", m.config.DBPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite connection: %w", err)
		}
		instance, err = sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create sqlite instance: %w", err)
		}
		driver = "sqlite3"
	default:
		return nil, fmt.Errorf("unsupported database type: %s", m.config.DBType)
	}

	// 创建迁移实例
	migrator, err := migrate.NewWithDatabaseInstance(
		m.config.MigrationsPath,
		driver,
		instance,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return migrator, nil
}
