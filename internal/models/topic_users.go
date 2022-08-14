package models

import (
	"context"
	"database/sql"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

//
const (
	TopicUserActionLiked      = "liked"
	TopicUserActionBookmarked = "bookmarked"
)

// TopicUser contains the relationships between topic and user
type TopicUser struct {
	TopicID      string
	UserID       string
	LikedAt      sql.NullTime
	BookmarkedAt sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time

	isNew bool
}

var topicUserColumns = []string{"topic_id", "user_id", "liked_at", "bookmarked_at", "created_at", "updated_at"}

func (tu *TopicUser) values() []interface{} {
	return []interface{}{tu.TopicID, tu.UserID, tu.LikedAt, tu.BookmarkedAt, tu.CreatedAt, tu.UpdatedAt}
}

func topicUserFromRow(row durable.Row) (*TopicUser, error) {
	var tu TopicUser
	err := row.Scan(&tu.TopicID, &tu.UserID, &tu.LikedAt, &tu.BookmarkedAt, &tu.CreatedAt, &tu.UpdatedAt)
	return &tu, err
}

// ActiondBy execute user action, like or bookmark a topic
func (topic *Topic) ActiondBy(ctx context.Context, user *User, action string, state bool) (*Topic, error) {
	if action != TopicUserActionLiked &&
		action != TopicUserActionBookmarked {
		return topic, session.BadDataError(ctx)
	}
	tu, err := readTopicUser(ctx, topic.TopicID, user.UserID)
	if err != nil {
		return topic, session.TransactionError(ctx, err)
	}
	if tu == nil {
		t := time.Now()
		tu = &TopicUser{
			TopicID:   topic.TopicID,
			UserID:    user.UserID,
			CreatedAt: t,
			UpdatedAt: t,

			isNew: true,
		}
	}
	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var lcount, bcount int64
		if action == TopicUserActionLiked {
			tu.LikedAt = sql.NullTime{Time: time.Now(), Valid: state}
			if err := tx.QueryRow(ctx, "SELECT count(*) FROM topic_users WHERE topic_id=$1 AND liked_at IS NOT NULL", topic.TopicID).Scan(&lcount); err != nil {
				return err
			}
			if lcount > 0 {
				topic.LikesCount = lcount - 1
			}
			if state {
				topic.LikesCount = lcount + 1
			}
		}
		if action == TopicUserActionBookmarked {
			tu.BookmarkedAt = sql.NullTime{Time: time.Now(), Valid: state}
			if err := tx.QueryRow(ctx, "SELECT count(*) FROM topic_users WHERE topic_id=$1 AND bookmarked_at IS NOT NULL", topic.TopicID).Scan(&bcount); err != nil {
				return err
			}
			if bcount > 0 {
				topic.BookmarksCount = bcount - 1
			}
			if state {
				topic.BookmarksCount = bcount + 1
			}
		}
		topic.IsLikedBy = tu.LikedAt.Valid
		topic.IsBookmarkedBy = tu.BookmarkedAt.Valid
		if _, err := tx.Exec(ctx, "UPDATE topics SET (likes_count,bookmarks_count)=($1,$2) WHERE topic_id=$3", topic.LikesCount, topic.BookmarksCount, topic.TopicID); err != nil {
			return err
		}
		if tu.isNew {
			rows := [][]interface{}{tu.values()}
			_, err := tx.CopyFrom(ctx, pgx.Identifier{"topic_users"}, topicUserColumns, pgx.CopyFromRows(rows))
			return err
		}
		column := "liked_at"
		if action == TopicUserActionBookmarked {
			column = "bookmarked_at"
		}
		query := fmt.Sprintf("UPDATE topic_users SET %s=$1 WHERE topic_id=$2 AND user_id=$3", column)
		if _, err := session.Database(ctx).Exec(ctx, query, sql.NullTime{Time: time.Now(), Valid: state}, tu.TopicID, tu.UserID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return topic, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func readTopicUser(ctx context.Context, topicID, userID string) (*TopicUser, error) {
	query := fmt.Sprintf("SELECT %s FROM topic_users WHERE topic_id=$1 AND user_id=$2", strings.Join(topicUserColumns, ","))
	row := session.Database(ctx).QueryRow(ctx, query, topicID, userID)
	tu, err := topicUserFromRow(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return tu, err
}

func fillTopicWithAction(ctx context.Context, topic *Topic, user *User) error {
	if user == nil {
		return nil
	}
	tu, err := readTopicUser(ctx, topic.TopicID, user.UserID)
	if err != nil || tu == nil {
		return err
	}
	topic.IsLikedBy, topic.IsBookmarkedBy = tu.LikedAt.Valid, tu.BookmarkedAt.Valid
	return nil
}
