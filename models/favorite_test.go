package models

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestFavoriteDao_Action_Add(t *testing.T) {
	favorite := &Favorite{
		UserId:  1,
		VideoId: 2,
	}

	// Expect the query to check if the favorite exists to return 0
	mock.ExpectQuery("SELECT count(*) FROM `favorite` WHERE (user_id = ? AND video_id = ?) AND `favorite`.`deleted_at` IS NULL").
		WithArgs(favorite.UserId, favorite.VideoId).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))

	// Expect the query to insert the favorite
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `favorite` (`user_id`,`video_id`,`created_at`,`deleted_at`) VALUES (?,?,?,?)").
		WithArgs(favorite.UserId, favorite.VideoId, sqlmock.AnyArg(), nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Expect the query to update the FavoriteCount of the video
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `video` SET `favorite_count`=favorite_count + ?,`updated_at`=? WHERE id = ? AND `video`.`deleted_at` IS NULL").
		WithArgs(1, sqlmock.AnyArg(), favorite.VideoId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Expect the query to update the TotalFavorited of the author
	mock.ExpectQuery("SELECT * FROM `video` WHERE id = ? AND `video`.`deleted_at` IS NULL ORDER BY `video`.`id` LIMIT 1").
		WithArgs(favorite.VideoId).
		WillReturnRows(sqlmock.NewRows([]string{"author_id"}).AddRow(3))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user` SET `total_favorited`=total_favorited + ?,`updated_at`=? WHERE id = ? AND `user`.`deleted_at` IS NULL").
		WithArgs(1, sqlmock.AnyArg(), 3).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Expect the query to update the FavoriteCount of the user
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user` SET `favorite_count`=favorite_count + ?,`updated_at`=? WHERE id = ? AND `user`.`deleted_at` IS NULL").
		WithArgs(1, sqlmock.AnyArg(), favorite.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := FavoriteDao().Action(favorite, true)
	require.NoError(t, err)
}

func TestFavoriteDao_Action_Remove(t *testing.T) {
	favorite := &Favorite{
		UserId:  1,
		VideoId: 2,
	}

	// Expect the query to delete the favorite
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `favorite` WHERE user_id = ? AND video_id = ?").
		WithArgs(favorite.UserId, favorite.VideoId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Expect the query to update the FavoriteCount of the video
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `video` SET `favorite_count`=favorite_count + ?,`updated_at`=? WHERE id = ? AND `video`.`deleted_at` IS NULL").
		WithArgs(-1, sqlmock.AnyArg(), favorite.VideoId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Expect the query to update the TotalFavorited of the author
	mock.ExpectQuery("SELECT * FROM `video` WHERE id = ? AND `video`.`deleted_at` IS NULL ORDER BY `video`.`id` LIMIT 1").
		WithArgs(favorite.VideoId).
		WillReturnRows(sqlmock.NewRows([]string{"author_id"}).AddRow(3))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user` SET `total_favorited`=total_favorited + ?,`updated_at`=? WHERE id = ? AND `user`.`deleted_at` IS NULL").
		WithArgs(-1, sqlmock.AnyArg(), 3).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Expect the query to update the FavoriteCount of the user
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user` SET `favorite_count`=favorite_count + ?,`updated_at`=? WHERE id = ? AND `user`.`deleted_at` IS NULL").
		WithArgs(-1, sqlmock.AnyArg(), favorite.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := FavoriteDao().Action(favorite, false)
	require.NoError(t, err)
}
