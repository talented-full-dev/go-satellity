package models

import (
	"context"
	"database/sql"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

// Category is used to categorize topics.
type Category struct {
	CategoryID  string
	Name        string
	Alias       string
	Description string
	TopicsCount int64
	LastTopicID sql.NullString
	Position    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var categoryColumns = []string{"category_id", "name", "alias", "description", "topics_count", "last_topic_id", "position", "created_at", "updated_at"}

func (c *Category) values() []interface{} {
	return []interface{}{c.CategoryID, c.Name, c.Alias, c.Description, c.TopicsCount, c.LastTopicID, c.Position, c.CreatedAt, c.UpdatedAt}
}

func categoryFromRows(row durable.Row) (*Category, error) {
	var c Category
	err := row.Scan(&c.CategoryID, &c.Name, &c.Alias, &c.Description, &c.TopicsCount, &c.LastTopicID, &c.Position, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

// CreateCategory create a new category with none blank name and alias, and optional description.
// alias use for human-readable, position for ordering categories
func CreateCategory(ctx context.Context, name, alias, description string, position int64) (*Category, error) {
	alias = strings.TrimSpace(alias)
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return nil, session.BadDataError(ctx)
	}
	if alias == "" {
		alias = name
	}

	t := time.Now()
	category := &Category{
		CategoryID:  uuid.Must(uuid.NewV4()).String(),
		Name:        name,
		Alias:       alias,
		Description: description,
		TopicsCount: 0,
		Position:    position,
		CreatedAt:   t,
		UpdatedAt:   t,
	}

	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		if position == 0 {
			last, err := lastCategoryPosition(ctx, tx)
			if err != nil {
				return err
			}
			category.Position = last + 1
		}

		rows := [][]interface{}{
			category.values(),
		}
		_, err := tx.CopyFrom(ctx, pgx.Identifier{"categories"}, categoryColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// UpdateCategory update fields of a category
func (category *Category) Update(ctx context.Context, name, alias, description string, position int64) error {
	alias, name = strings.TrimSpace(alias), strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		if name != "" {
			category.Name = name
		}
		if alias != "" {
			category.Alias = alias
		}
		category.Description = description
		category.Position = position
		category.UpdatedAt = time.Now()
		cols, posits := durable.PrepareColumnsAndExpressions([]string{"name", "alias", "description", "position", "updated_at"}, 1)
		values := []interface{}{category.CategoryID, category.Name, category.Alias, category.Description, category.Position, category.UpdatedAt}
		_, err := tx.Exec(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id=$1", cols, posits), values...)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

// ReadCategory read a category by ID
func ReadCategory(ctx context.Context, id string) (*Category, error) {
	var category *Category
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// ReadCategoryByIDOrName read a category by id or name
func ReadCategoryByIDOrName(ctx context.Context, identity string) (*Category, error) {
	query := fmt.Sprintf("SELECT %s FROM categories WHERE category_id=$1 OR name=$1", strings.Join(categoryColumns, ","))
	row := session.Database(ctx).QueryRow(ctx, query, identity)
	c, err := categoryFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return c, nil
}

// ReadAllCategories read categories order by position
func ReadAllCategories(ctx context.Context) ([]*Category, error) {
	var categories []*Category
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		categories, err = readCategories(ctx, tx)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return categories, nil
}

func readCategories(ctx context.Context, tx pgx.Tx) ([]*Category, error) {
	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT %s FROM categories ORDER BY position LIMIT 500", strings.Join(categoryColumns, ",")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category, err := categoryFromRows(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func readCategorySet(ctx context.Context, tx pgx.Tx, ids []string) (map[string]*Category, error) {
	categories, err := readCategoriesByIds(ctx, tx, ids)
	if err != nil {
		return nil, err
	}
	set := make(map[string]*Category, 0)
	for _, c := range categories {
		set[c.CategoryID] = c
	}
	return set, nil
}

func readCategoriesByIds(ctx context.Context, tx pgx.Tx, ids []string) ([]*Category, error) {
	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT %s FROM categories WHERE category_id IN ('%s') LIMIT 100", strings.Join(categoryColumns, ","), strings.Join(ids, "','")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category, err := categoryFromRows(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

// EmitToCategory update category's info, e.g.: LastTopicID, TopicsCount
func EmitToCategory(ctx context.Context, id string) (*Category, error) {
	var category *Category
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
		if err != nil || category == nil {
			return err
		}
		topic, err := category.latestTopic(ctx, tx)
		if err != nil {
			return err
		}
		lastTopicID := sql.NullString{String: "", Valid: false}
		if topic != nil {
			lastTopicID = sql.NullString{String: topic.TopicID, Valid: true}
		}
		if category.LastTopicID.String != lastTopicID.String {
			category.LastTopicID = lastTopicID
		}
		category.TopicsCount = 0
		if category.LastTopicID.Valid {
			count, err := fetchTopicsCount(ctx, tx, category.CategoryID)
			if err != nil {
				return err
			}
			category.TopicsCount = count
		}
		category.UpdatedAt = time.Now()
		cols, posits := durable.PrepareColumnsAndExpressions([]string{"last_topic_id", "topics_count", "updated_at"}, 1)
		values := []interface{}{category.CategoryID, category.LastTopicID, category.TopicsCount, category.UpdatedAt}
		_, err = tx.Exec(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id=$1", cols, posits), values...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

func findCategory(ctx context.Context, tx pgx.Tx, id string) (*Category, error) {
	if uuid.FromStringOrNil(id).String() != id {
		return nil, nil
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM categories WHERE category_id=$1", strings.Join(categoryColumns, ",")), id)
	c, err := categoryFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func lastCategoryPosition(ctx context.Context, tx pgx.Tx) (int64, error) {
	var position int64
	row := tx.QueryRow(ctx, "SELECT position FROM categories ORDER BY position DESC LIMIT 1")
	err := row.Scan(&position)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return position, err
}
