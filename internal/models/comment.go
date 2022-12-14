package models

import (
	"context"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	commentBodySizeLimit = 6
)

// Comment is struct for comment of topic
type Comment struct {
	CommentID string
	Body      string
	TopicID   string
	UserID    string
	Score     int
	CreatedAt time.Time
	UpdatedAt time.Time

	User *User
}

var commentColumns = []string{"comment_id", "body", "topic_id", "user_id", "score", "created_at", "updated_at"}

func (c *Comment) values() []interface{} {
	return []interface{}{c.CommentID, c.Body, c.TopicID, c.UserID, c.Score, c.CreatedAt, c.UpdatedAt}
}

func commentFromRows(row durable.Row) (*Comment, error) {
	var c Comment
	err := row.Scan(&c.CommentID, &c.Body, &c.TopicID, &c.UserID, &c.Score, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

// CreateComment create a new comment
func (user *User) CreateComment(ctx context.Context, body string, topic *Topic) (*Comment, error) {
	body = strings.TrimSpace(body)
	if len(body) < commentBodySizeLimit {
		return nil, session.BadDataError(ctx)
	}
	t := time.Now()
	c := &Comment{
		CommentID: uuid.Must(uuid.NewV4()).String(),
		Body:      body,
		UserID:    user.UserID,
		CreatedAt: t,
		UpdatedAt: t,
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		count, err := fetchCommentsCount(ctx, tx, topic.TopicID)
		if err != nil {
			return err
		}
		topic.CommentsCount = count + 1
		if topic.UpdatedAt.Add(time.Hour * 24 * 30).After(time.Now()) {
			topic.UpdatedAt = t
		}
		cols, posits := durable.PrepareColumnsAndExpressions([]string{"comments_count", "updated_at"}, 1)
		query := fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id=$1", cols, posits)
		_, err = tx.Exec(ctx, query, topic.TopicID, topic.CommentsCount, topic.UpdatedAt)
		if err != nil {
			return err
		}
		c.TopicID = topic.TopicID
		rows := [][]interface{}{c.values()}
		_, err = tx.CopyFrom(ctx, pgx.Identifier{"comments"}, commentColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	c.User = user
	UpsertStatistic(ctx, StatisticTypeComments)
	return c, nil
}

// UpdateComment update the comment by id
func (comment *Comment) Update(ctx context.Context, body string, user *User) error {
	if !comment.isPermit(user) {
		return session.ForbiddenError(ctx)
	}
	body = strings.TrimSpace(body)
	if len(body) < commentBodySizeLimit {
		return session.BadDataError(ctx)
	}
	comment.Body = body
	comment.UpdatedAt = time.Now()
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		cols, posits := durable.PrepareColumnsAndExpressions([]string{"body", "updated_at"}, 1)
		_, err := tx.Exec(ctx, fmt.Sprintf("UPDATE comments SET (%s)=(%s) WHERE comment_id=$1", cols, posits), comment.CommentID, comment.Body, comment.UpdatedAt)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

// ReadComments read all comments, parameters: offset default time.Now()
func ReadComments(ctx context.Context, offset time.Time, topic *Topic, user *User) ([]*Comment, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	query := fmt.Sprintf("SELECT %s FROM comments WHERE updated_at<$1 ORDER BY updated_at DESC LIMIT $2", strings.Join(commentColumns, ","))
	params := []any{offset, LIMIT}
	if topic != nil {
		query = fmt.Sprintf("SELECT %s FROM comments WHERE topic_id=$1 AND created_at<$2 ORDER BY created_at LIMIT $3", strings.Join(commentColumns, ","))
		params = append([]any{topic.TopicID}, params...)
	}
	if user != nil {
		query = fmt.Sprintf("SELECT %s FROM comments WHERE user_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(commentColumns, ","))
		params = append([]any{user.UserID}, params...)
	}

	var comments []*Comment
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, query, params...)
		if err != nil {
			return err
		}
		defer rows.Close()

		var userIds []string
		for rows.Next() {
			comment, err := commentFromRows(rows)
			if err != nil {
				return err
			}
			comment.User = user
			if comment.User == nil {
				userIds = append(userIds, comment.UserID)
			}
			comments = append(comments, comment)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		if len(userIds) > 0 {
			userSet, err := readUserSet(ctx, tx, userIds)
			if err != nil {
				return err
			}
			for i, comment := range comments {
				comments[i].User = userSet[comment.UserID]
			}
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return comments, nil
}

// DeleteComment delete a comment by ID
func (comment *Comment) Delete(ctx context.Context, user *User) error {
	if !comment.isPermit(user) {
		return session.ForbiddenError(ctx)
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		topic, err := findTopic(ctx, tx, comment.TopicID)
		if err != nil {
			return err
		} else if topic == nil {
			return session.BadDataError(ctx)
		}
		count, err := fetchCommentsCount(ctx, tx, comment.TopicID)
		if err != nil {
			return err
		}
		topic.CommentsCount = count - 1
		cols, posits := durable.PrepareColumnsAndExpressions([]string{"comments_count", "updated_at"}, 1)
		query := fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id=$1", cols, posits)
		_, err = tx.Exec(ctx, query, topic.TopicID, topic.CommentsCount, topic.UpdatedAt)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "DELETE FROM comments WHERE comment_id=$1", comment.CommentID)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func ReadComment(ctx context.Context, id string) (*Comment, error) {
	var comment *Comment
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		comment, err = findComment(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return comment, nil
}

func findComment(ctx context.Context, tx pgx.Tx, id string) (*Comment, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM comments WHERE comment_id=$1", strings.Join(commentColumns, ",")), id)
	c, err := commentFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func fetchCommentsCount(ctx context.Context, tx pgx.Tx, topicID string) (int64, error) {
	var count int64
	query := "SELECT count(*) FROM comments"
	params := []any{}
	if uuid.FromStringOrNil(topicID).String() == topicID {
		query = "SELECT count(*) FROM comments WHERE topic_id=$1"
		params = []any{topicID}
	}
	err := tx.QueryRow(ctx, query, params...).Scan(&count)
	return count, err
}

func (comment *Comment) isPermit(user *User) bool {
	if user == nil {
		return false
	}
	if user.isAdmin() {
		return true
	}
	return comment.UserID == user.UserID
}
