package converter

import (
	"route256/loms/internal/model"
	"route256/loms/internal/repository/schema"
)

func SchemaToKafkaRecordListModel(kafkaRecords []*schema.KafkaRecord) []*model.KafkaRecord {
	if kafkaRecords == nil {
		return nil
	}

	result := make([]*model.KafkaRecord, 0, len(kafkaRecords))
	for _, kafkaRecord := range kafkaRecords {
		result = append(result, SchemaToKafkaRecordModel(kafkaRecord))
	}
	return result
}

func SchemaToKafkaRecordModel(record *schema.KafkaRecord) *model.KafkaRecord {
	if record == nil {
		return nil
	}

	return &model.KafkaRecord{
		ID:           record.ID,
		Topic:        record.Topic,
		PartitionKey: record.PartitionKey,
		Message:      record.Message,
		State:        model.RecordState(record.State),
		CreatedOn:    record.CreatedOn,
	}
}
