package postgres

import (
	"context"
	"route256/loms/internal/converter"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type kafkaOutboxRepository struct {
	queryEngineProvider QueryEngineProvider
}

const (
	outboxRecordTable = "outbox_record"

	outboxRecordIdColumn        = "outbox_record_id"
	outboxRecordTopicColumn     = "topic"
	outboxRecordKeyColumn       = "key"
	outboxRecordMessageColumn   = "message"
	outboxRecordStateColumn     = "state"
	outboxRecordCreatedAtColumn = "created_on"
)

func NewKafkaOutboxRepository(queryEngineProvider QueryEngineProvider) *kafkaOutboxRepository {
	return &kafkaOutboxRepository{
		queryEngineProvider: queryEngineProvider,
	}
}

func (r *kafkaOutboxRepository) GetUnprocessedRecords(ctx context.Context) ([]*model.KafkaRecord, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Select(outboxRecordIdColumn, outboxRecordTopicColumn, outboxRecordKeyColumn, outboxRecordMessageColumn, outboxRecordStateColumn, outboxRecordCreatedAtColumn).
		From(outboxRecordTable).
		Where(sq.Eq{outboxRecordStateColumn: schema.PendingDelivery}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var result []*schema.KafkaRecord
	if err := pgxscan.Select(ctx, db, &result, query, args...); err != nil {
		return nil, err
	}

	records := converter.SchemaToKafkaRecordListModel(result)
	return records, nil
}

func (r *kafkaOutboxRepository) UpdateRecordByID(ctx context.Context, message *model.KafkaRecord) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Update(outboxRecordTable).
		Set(outboxRecordStateColumn, message.State).
		Where(sq.Eq{outboxRecordIdColumn: message.ID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *kafkaOutboxRepository) RemoveRecordsBeforeDatetime(ctx context.Context, expireTime time.Time) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Delete(outboxRecordTable).
		Where(sq.LtOrEq{outboxRecordCreatedAtColumn: expireTime}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *kafkaOutboxRepository) CreateKafkaRecord(ctx context.Context, record *model.KafkaRecord) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Insert(outboxRecordTable).
		Columns(outboxRecordTopicColumn, outboxRecordKeyColumn, outboxRecordMessageColumn, outboxRecordCreatedAtColumn, outboxRecordStateColumn).
		Values(record.Topic, record.PartitionKey, record.Message, record.CreatedOn, schema.PendingDelivery).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
