package postgres

import sq "github.com/Masterminds/squirrel"

func queryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
