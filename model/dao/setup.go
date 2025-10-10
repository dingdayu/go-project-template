// Package dao 所有定义的模型操作方法，均在此包里。
package dao

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"net/url"
	"strings"

	"github.com/dingdayu/singbox-adapter/pkg/logger"
	pkgOtel "github.com/dingdayu/singbox-adapter/pkg/otel"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	once sync.Once
	db   *gorm.DB

	tracer = otel.Tracer("github.com/dingdayu/singbox-adapter/model/dao")
)

func Init() {
	var err error
	once.Do(func() {
		dsn := viper.GetString("db.dsn")
		dbCfg := &gorm.Config{
			Logger: logger.NewGormLogger(logger.WithNamespace("gorm")),
		}

		var dialector gorm.Dialector

		// try parse dsn as URL to detect scheme first
		if u, perr := url.Parse(dsn); perr == nil && u.Scheme != "" {
			switch strings.ToLower(u.Scheme) {
			case "postgres", "postgresql":
				dialector = postgres.Open(dsn)
			case "mysql":
				dialector = mysql.Open(dsn)
			case "sqlite", "file":
				// sqlite dsn might be like file:test.db?_foreign_keys=1 or just path
				// pass raw dsn to sqlite driver
				dialector = sqlite.Open(dsn)
			}
		}

		// fallback to simple prefix/content checks when URL parsing didn't help
		if dialector == nil {
			ld := strings.ToLower(dsn)
			switch {
			case strings.HasPrefix(ld, "postgres"), strings.HasPrefix(ld, "postgresql"), strings.Contains(ld, "host="):
				dialector = postgres.Open(dsn)
			case strings.HasPrefix(ld, "mysql"), strings.Contains(ld, "@tcp("), strings.Contains(ld, "charset="):
				dialector = mysql.Open(dsn)
			case strings.HasSuffix(ld, ".db"), strings.HasSuffix(ld, ".sqlite"), strings.HasPrefix(ld, "file:"):
				dialector = sqlite.Open(dsn)
			}
		}

		if dialector == nil {
			fmt.Printf("\033[1;30;41m[error]\033[0m db [master] connect error: unsupported dsn or driver for '%s'\n", dsn)
			os.Exit(1)
		}

		db, err = gorm.Open(dialector, dbCfg)
		if err == nil {
			var databaseName string
			// 通过 session 方法设置 logger 级别为 Silent，避免打印 SQL 语句
			pkgOtel.NeverLogSessionWithGORM(db).Raw("SELECT current_database()").Scan(&databaseName)
			fmt.Printf("\033[1;30;42m[info]\033[0m db [master:%s] connect success\n", databaseName)

			// 启用 gorm otel 追踪
			if err := db.Use(tracing.NewPlugin()); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("\033[1;30;41m[error]\033[0m db [master] connect error: %s", err.Error())
			os.Exit(1)
		}
	})
}

func GetDB() *gorm.DB {
	return db
}

func GetContextDB(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

type JSON json.RawMessage

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// MarshalJSON implements json.Marshaler interface
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSON: UnmarshalJSON on nil pointer")
	}
	*j = JSON(data)
	return nil
}

func (JSON) GormDataType() string {
	return "jsonb"
}

func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// use field.Tag, field.TagSettings gets field's tags
	// checkout https://github.com/go-gorm/gorm/blob/master/schema/field.go for all options

	// returns different database type based on driver name
	switch db.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

type NullableField[T any] struct {
	Set   bool
	Value *T
}

func (f *NullableField[T]) UnmarshalJSON(b []byte) error {
	f.Set = true
	if string(b) == "null" {
		f.Value = nil
		return nil
	}
	var val T
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}
	f.Value = &val
	return nil
}
