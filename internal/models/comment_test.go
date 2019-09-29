package models

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCommentCRUD(t *testing.T) {
	assert := assert.New(t)
	mctx := setupTestContext()
	defer mctx.database.Close()
	defer teardownTestContext(mctx)

	user := createTestUser(mctx, "im.yuqlee@gmail.com", "username", "password")
	assert.NotNil(user)
	category, _ := CreateCategory(mctx, "name", "alias", "Description", 0)
	assert.NotNil(category)
	topic, _ := user.CreateTopic(mctx, "title", "body", category.CategoryID, false)
	assert.NotNil(topic)

	commentCases := []struct {
		topicID string
		body    string
		valid   bool
	}{
		{topic.TopicID, "", false},
		{topic.TopicID, "      ", false},
		{uuid.Must(uuid.NewV4()).String(), "comment body", false},
		{topic.TopicID, "comment body", true},
	}

	for _, tc := range commentCases {
		t.Run(fmt.Sprintf("comment body %s", tc.body), func(t *testing.T) {
			if !tc.valid {
				comment, err := user.CreateComment(mctx, tc.topicID, tc.body)
				assert.NotNil(err)
				assert.Nil(comment)
				return
			}

			comment, err := user.CreateComment(mctx, tc.topicID, tc.body)
			assert.Nil(err)
			assert.NotNil(comment)
			assert.Equal(tc.body, comment.Body)
			new, err := readTestComment(mctx, comment.CommentID)
			assert.Nil(err)
			assert.NotNil(new)
			new, err = readTestComment(mctx, uuid.Must(uuid.NewV4()).String())
			assert.Nil(err)
			assert.Nil(new)
			new, err = user.UpdateComment(mctx, uuid.Must(uuid.NewV4()).String(), "comment body")
			assert.NotNil(err)
			assert.Nil(new)
			new, err = user.UpdateComment(mctx, comment.CommentID, "    ")
			assert.NotNil(err)
			assert.Nil(new)
			new, err = user.UpdateComment(mctx, comment.CommentID, "new comment body")
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal("new comment body", new.Body)
			comments, err := topic.ReadComments(mctx, time.Time{})
			assert.Nil(err)
			assert.Len(comments, 1)
			comments, err = user.ReadComments(mctx, time.Time{})
			assert.Nil(err)
			assert.Len(comments, 1)
			topic, _ = ReadTopic(mctx, topic.TopicID)
			assert.NotNil(topic)
			assert.Equal(int64(1), topic.CommentsCount)
			err = user.DeleteComment(mctx, comment.CommentID)
			assert.Nil(err)
			topic, err = ReadTopic(mctx, topic.TopicID)
			assert.Nil(err)
			assert.NotNil(topic)
			assert.Equal(int64(0), topic.CommentsCount)
			comments, err = topic.ReadComments(mctx, time.Time{})
			assert.Nil(err)
			assert.Len(comments, 0)
			comments, err = user.ReadComments(mctx, time.Time{})
			assert.Nil(err)
			assert.Len(comments, 0)
			new, err = readTestComment(mctx, comment.CommentID)
			assert.Nil(err)
			assert.Nil(new)
		})
	}
}

func readTestComment(mctx *Context, id string) (*Comment, error) {
	ctx := mctx.context
	var comment *Comment
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		comment, err = findComment(ctx, tx, id)
		return err
	})
	return comment, err
}
