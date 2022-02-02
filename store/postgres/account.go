package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/snowzach/queryp"
	"github.com/snowzach/queryp/qppg"

	gtsasamiserver "github.com/Tim-vo/gt-sasami-server/gt-sasami-server"
	"github.com/Tim-vo/gt-sasami-server/store"
)

const (
	AccountSchema = ``
	AccountTable  = `account`
	AccountJoins  = ``
	AccountFields = /* sql */ `COALESCE(account.id) as "account.id",
	COALESCE(account.email, '') as "account.email",
	COALESCE(account.username, '') as "account.username",
	COALESCE(account.passowrd, '') as "account.password",
	COALESCE(account.created, '0001-01-01 00:00:00 UTC') as "account.created",
	COALESCE(account.updated, '0001-01-01 00:00:00 UTC') as "account.updated",
	`
)

var (
	AccountSelect = /* sql */ `SELECT` + strings.Join([]string{
		"accpunt.*",
	}, ",")
)

func (client *Client) AccountSave(ctx context.Context, record *gtsasamiserver.Account) error {

	if record.ID == "" {
		record.ID = client.newID()
	}

	fields, values, updates, args := composeUpsert([]field{
		{name: "id", insert: "$#", update: "", arg: record.ID},
		{name: "email", insert: "$#", update: "$#", arg: record.Email},
		{name: "username", insert: "$#", update: "$#", arg: record.Username},
		{name: "password", insert: "$#", update: "$#", arg: record.Password},
		{name: "created", insert: "NOW()", update: ""},
		{name: "updated", insert: "", update: "NOW()"},
	})

	err := client.db.GetContext(ctx, record, `
	WITH `+AccountTable+` AS (
        INSERT INTO `+AccountSchema+AccountTable+` (`+fields+`)
        VALUES(`+values+`) ON CONFLICT (id) DO UPDATE
        SET `+updates+` RETURNING *
	) `+AccountSelect+" FROM "+AccountTable+AccountJoins, args...)

	if err != nil {
		return wrapError(err)
	}

	return nil
}

func (client *Client) AccountGetByID(ctx context.Context, id string) (*gtsasamiserver.Account, error) {

	account := new(gtsasamiserver.Account)
	err := client.db.GetContext(ctx, account, AccountSelect+` FROM `+AccountSchema+AccountTable+AccountJoins+` WHERE `+AccountTable+`.id = $1`, id)

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, wrapError(err)
	}
	return account, nil
}

func (client *Client) AccountDeleteByID(ctx context.Context, id string) error {

	_, err := client.db.ExecContext(ctx, `DELETE FROM `+AccountSchema+AccountTable+` WHERE `+AccountTable+`.id = $1`, id)
	if err != nil {
		return wrapError(err)
	}
	return nil
}

func (client *Client) AccountsFind(ctx context.Context, qp *queryp.QueryParameters) ([]*gtsasamiserver.Account, int64, error) {

	var queryClause strings.Builder
	var queryParams = []interface{}{}

	filterFields := queryp.FilterFieldTypes{
		"account.id":       queryp.FilterTypeSimple,
		"account.email":    queryp.FilterTypeString,
		"account.username": queryp.FilterTypeString,
		"account.password": queryp.FilterTypeString,
	}

	sortFields := queryp.SortFields{
		"account.id":       "",
		"account.created":  "",
		"account.updated":  "",
		"account.username": "",
		"account.email":    "",
	}
	// Default sort
	if len(qp.Sort) == 0 {
		qp.Sort.Append("account.id", false)
	}

	if len(qp.Filter) > 0 {
		queryClause.WriteString(" WHERE ")
	}

	if err := qppg.FilterQuery(filterFields, qp.Filter, &queryClause, &queryParams); err != nil {
		return nil, 0, &store.Error{Type: store.ErrorTypeQuery, Err: err}
	}
	var count int64
	if err := client.db.GetContext(ctx, &count, `SELECT COUNT(*) AS count FROM `+AccountSchema+AccountTable+AccountJoins+queryClause.String(), queryParams...); err != nil {
		return nil, 0, wrapError(err)
	}
	if err := qppg.SortQuery(sortFields, qp.Sort, &queryClause, &queryParams); err != nil {
		return nil, 0, &store.Error{Type: store.ErrorTypeQuery, Err: err}
	}
	if qp.Limit > 0 {
		queryClause.WriteString(" LIMIT " + strconv.FormatInt(qp.Limit, 10))
	}
	if qp.Offset > 0 {
		queryClause.WriteString(" OFFSET " + strconv.FormatInt(qp.Offset, 10))
	}

	var records = make([]*gtsasamiserver.Account, 0)
	err := client.db.SelectContext(ctx, &records, AccountSelect+` FROM `+AccountSchema+AccountTable+AccountJoins+queryClause.String(), queryParams...)
	if err != nil {
		return records, 0, wrapError(err)
	}

	return records, count, nil
}
