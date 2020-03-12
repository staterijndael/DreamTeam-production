package stores

import (
	"dt/config"
	"dt/logwrap"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

const (
	zdbIdxTypeSuffix = "zdb_idx_type"
	zdbIdxSuffix     = "zdb_idx"
	schema           = "public"
	zdbIdxColumnName = "ctid"
)

type ZDBIndexer interface {
	TableName() string
	ZDBIdxTypeDefinition() string
	ZDBRowBuilder() string
	ZDBIdxIDType() string
}

func zdbInit(db *gorm.DB, c *config.Config, tables ...ZDBIndexer) error {
	if err := db.Exec("create extension if not exists zombodb").Error; err != nil {
		return err
	}

	for _, table := range tables {
		oldIdxDef, isExists, err := getOldZDBIdxDefIfExists(db, table)
		if err != nil {
			logwrap.Debug("[zdb init]: %s", err.Error())
			return err
		}

		if isExists {
			if isSameZDBIdxDef(c.ElasticSearchUrl, oldIdxDef, table) {
				continue
			}
		}

		if err := dropZDBIdx(db, table); err != nil {
			logwrap.Debug("[zdb init]: %s", err.Error())
			return err
		}

		if err := zdbInitTable(db, c.ElasticSearchUrl, table); err != nil {
			logwrap.Debug("[zdb init]: %s", err.Error())
			return err
		}
	}

	return nil
}

func dropZDBIdx(db *gorm.DB, indexer ZDBIndexer) error {
	if err := db.Exec("drop index if exists " + indexer.TableName() + "_" + zdbIdxSuffix).Error; err != nil {
		return err
	}

	return db.Exec("drop type if exists " + indexer.TableName() + "_" + zdbIdxTypeSuffix).Error
}

func getOldZDBIdxDefIfExists(db *gorm.DB, indexer ZDBIndexer) (old string, exists bool, err error) {
	old, err = getZDBIdxDef(db, indexer)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return "", false, nil
		}

		return "", false, err
	}

	return old, true, nil
}

func isSameZDBIdxDef(elasticUrl, oldIdx string, indexer ZDBIndexer) bool {
	idxExpr := createZDBIdxExpr(indexer, elasticUrl)
	return strings.Contains(strings.ToLower(oldIdx), strings.ToLower(idxExpr)[:len(idxExpr)-1])
}

func zdbInitTable(db *gorm.DB, elasticUrl string, table ZDBIndexer) error {
	err := db.Exec(createZDBIdxTypeExpr(table)).Error
	if err != nil {
		return err
	}

	return db.Exec(createZDBIdxExpr(table, elasticUrl)).Error
}

func createZDBIdxTypeExpr(indexer ZDBIndexer) string {
	return fmt.Sprintf(
		`create type %s_%s as (%s)`,
		indexer.TableName(),
		zdbIdxTypeSuffix,
		indexer.ZDBIdxTypeDefinition(),
	)
}

func createZDBIdxExpr(indexer ZDBIndexer, elasticUrl string) string {
	return fmt.Sprintf(
		"create index %s_%s on %s.%s using zombodb (%s, (row(%s)::%s_%s)) with (url='%s')",
		indexer.TableName(),
		zdbIdxSuffix,
		schema,
		indexer.TableName(),
		zdbIdxColumnName,
		indexer.ZDBRowBuilder(),
		indexer.TableName(),
		zdbIdxTypeSuffix,
		elasticUrl,
	)
}

func getZDBIdxDef(db *gorm.DB, indexer ZDBIndexer) (string, error) {
	var it IndexTable
	_db := db.Raw("select indexdef from pg_indexes where tablename = ? and schemaname = ? and indexname = ?",
		indexer.TableName(), schema, indexer.TableName()+"_"+zdbIdxSuffix)
	if err := _db.Error; err != nil {
		return "", err
	}

	if err := _db.Scan(&it).Error; err != nil {
		return "", err
	}

	return it.IndexDef, nil
}

type IndexTable struct {
	IndexDef string `gorm:"column:indexdef"`
}
