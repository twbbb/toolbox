package cachecount

import (
	"context"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"gorm.io/gorm"
)

// CountInsertDb ..
type CountInsertDb struct {
	table     string
	sql       string
	db        *gorm.DB
	usageData map[int]int64
}

// BatchUpdate 批量更新使用统计（支持外部传入 db）
// usageData: key为f_id，value为使用次数增量
// 每1000条提交一次事务，避免长事务锁定过多资源
func (impl *CountInsertDb) BatchUpdate(
	ctx context.Context,
) (effectRow int64, err error) {
	if len(impl.usageData) == 0 {
		return 0, nil
	}
	if impl.sql == "" && impl.table == "" {
		log.ErrorContextf(ctx, "BatchUpdateReportTranslationUsage: both sql and table are empty")
		return 0, gorm.ErrInvalidData
	}
	if impl.db == nil {
		log.ErrorContextf(ctx, "BatchUpdateReportTranslationUsage: db is nil")
		return 0, gorm.ErrInvalidData
	}

	// 1. 拆分map成多个小map
	batches := impl.splitToBatches(impl.usageData)
	log.InfoContextf(
		ctx, "BatchUpdateReportTranslationUsage split into %d batches, total items: %d", len(batches),
		len(impl.usageData),
	)

	var totalAffectedRows int64
	// 2. 逐个处理每个批次
	for i, batch := range batches {
		batchNum := i + 1
		log.InfoContextf(ctx, "Processing batch %d/%d, items: %d", batchNum, len(batches), len(batch))

		batchAffectedRows, err := impl.updateInTx(ctx, batch)
		if err != nil {
			log.ErrorContextf(ctx, "BatchUpdateReportTranslationUsage batch %d failed, error: %v", batchNum, err)
			return totalAffectedRows, err
		}

		totalAffectedRows += batchAffectedRows
		log.InfoContextf(
			ctx, "BatchUpdateReportTranslationUsage batch %d success, batch affected rows: %d, total affected rows: %d",
			batchNum, batchAffectedRows, totalAffectedRows,
		)
	}

	log.InfoContextf(
		ctx, "BatchUpdateReportTranslationUsage all batches completed, total affected rows: %d, total items: %d",
		totalAffectedRows, len(impl.usageData),
	)
	return totalAffectedRows, nil
}

// split 将大map拆分成多个小map，每个map最多包含batchSize个元素
func (impl *CountInsertDb) splitToBatches(usageData map[int]int64) []map[int]int64 {
	if len(usageData) == 0 {
		return nil
	}
	batchSize := 1000
	if len(usageData) < batchSize {
		return []map[int]int64{usageData}
	}

	var batches []map[int]int64
	currentBatch := make(map[int]int64)
	count := 0

	for fid, usage := range usageData {
		currentBatch[fid] = usage
		count++

		if count >= batchSize {
			batches = append(batches, currentBatch)
			currentBatch = make(map[int]int64)
			count = 0
		}
	}

	// 添加最后一个批次（如果有剩余数据）
	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	return batches
}

// updateInTx 在事务中更新使用统计
func (impl *CountInsertDb) updateInTx(
	ctx context.Context,
	usageData map[int]int64,
) (effectRow int64, err error) {
	if len(usageData) == 0 {
		return 0, nil
	}

	// 开启事务
	tx := impl.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.ErrorContextf(ctx, "updateUsageInTransaction begin transaction failed, error: %v", tx.Error)
		return 0, tx.Error
	}

	// 确保事务回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.ErrorContextf(ctx, "updateUsageInTransaction panic recovered: %v", r)
		}
	}()

	var totalAffectedRows int64
	currentTime := time.Now()

	// 逐条更新
	for fid, count := range usageData {
		var result *gorm.DB

		if impl.sql != "" {
			// 使用自定义 SQL
			result = tx.Exec(impl.sql, count, currentTime, fid)
		} else if impl.table != "" {
			// 使用表名动态构建 SQL（注意：表名不能用占位符，需要拼接）
			sql := "UPDATE ? SET f_usage_count = f_usage_count + ?, f_usage_time = ?, " +
				"f_update_time = ? WHERE f_id = ?"
			result = tx.Exec(sql, impl.table, count, currentTime, currentTime, fid)
		} else {
			// 如果既没有 sql 也没有 table，回滚事务并返回错误
			tx.Rollback()
			log.ErrorContextf(ctx, "updateUsageInTransaction: both sql and table are empty")
			return 0, gorm.ErrInvalidData
		}

		if result.Error != nil {
			tx.Rollback()
			log.ErrorContextf(ctx, "updateUsageInTransaction update failed, fid: %d, error: %v", fid, result.Error)
			return 0, result.Error
		}
		totalAffectedRows += result.RowsAffected
		log.DebugContextf(
			ctx, "Updated report translation usage, fid: %d, count: %d, affected rows: %d", fid, count,
			result.RowsAffected,
		)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.ErrorContextf(ctx, "updateUsageInTransaction commit transaction failed, error: %v", err)
		return 0, err
	}

	log.InfoContextf(
		ctx, "updateUsageInTransaction success, total affected rows: %d, updated count: %d", totalAffectedRows,
		len(usageData),
	)
	return totalAffectedRows, nil
}
