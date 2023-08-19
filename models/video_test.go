package models

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestVideoDao_Add(t *testing.T) {

	video := &Video{
		AuthorId: 1,
		PlayUrl:  "https://example.com/video.mp4",
		Title:    "Test Video",
	}

	// Expect the query to insert the video
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user` SET `work_count`=work_count + ?,`updated_at`=? WHERE id = ? AND `user`.`deleted_at` IS NULL").
		WithArgs(1, sqlmock.AnyArg(), video.AuthorId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `video` (`created_at`,`updated_at`,`deleted_at`,`author_id`,`play_url`,`cover_url`,`favorite_count`,`comment_count`,`title`) VALUES (?,?,?,?,?,?,?,?,?)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, video.AuthorId, video.PlayUrl, video.CoverUrl, 0, 0, video.Title).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := VideoDao().Add(video)

	require.NoError(t, err)
	assert.Equal(t, video, result)
}

func TestVideoDao_Add_InvalidAuthor(t *testing.T) {

	video := &Video{
		AuthorId: 1,
		PlayUrl:  "https://example.com/video.mp4",
		Title:    "Test Video",
	}

	// Expect the query to check if the author exists
	mock.ExpectQuery("SELECT count(*) FROM `user` WHERE id = ? AND `user`.`deleted_at` IS NULL").
		WithArgs(video.AuthorId).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))

	result, err := VideoDao().Add(video)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestVideoDao_Add_DatabaseError(t *testing.T) {

	video := &Video{
		AuthorId: 1,
		PlayUrl:  "https://example.com/video.mp4",
		Title:    "Test Video",
	}

	// Expect the query to check if the author exists
	mock.ExpectQuery("SELECT count(*) FROM `user` WHERE id = ? AND `user`.`deleted_at` IS NULL").
		WithArgs(video.AuthorId).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))

	// Expect the query to insert the video to fail
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `video`").
		WithArgs(video.AuthorId, video.PlayUrl, video.CoverUrl, video.FavoriteCount, video.CommentCount, video.Title).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	result, err := VideoDao().Add(video)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestVideoDao_GetBefore(t *testing.T) {
	video1 := &Video{
		Id:    2,
		Title: "Video 2",
		Model: gorm.Model{
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
	}
	video2 := &Video{
		Id:    3,
		Title: "Video 3",
		Model: gorm.Model{
			CreatedAt: time.Now().Add(-3 * time.Hour),
		},
	}

	// Expect the query to retrieve videos before the given timestamp
	mock.ExpectQuery("SELECT * FROM `video` WHERE created_at < ? AND `video`.`deleted_at` IS NULL ORDER BY created_at desc LIMIT 2").
		WithArgs(time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "created_at"}).
			AddRow(video1.Id, video1.Title, video1.CreatedAt).
			AddRow(video2.Id, video2.Title, video2.CreatedAt))
	result, oldest, err := VideoDao().GetBefore(time.Now().Add(-2*time.Hour).Unix(), 2)

	require.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, video1.Id, result[0].Id)
	assert.Equal(t, video2.Id, result[1].Id)
	assert.Equal(t, video2.CreatedAt.Unix(), oldest)
}
