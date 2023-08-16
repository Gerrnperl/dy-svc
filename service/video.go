package service

import (
	"main/config"
	"main/models"
	"main/utils"
	"mime/multipart"
	"strconv"
	"time"
)

type ErrVideoFormat struct {
	format string
}

func (e ErrVideoFormat) Error() string {
	return "invalid video format: " + e.format
}

// UploadVideo
//
// uploads a video file to the server and adds a new video record to the database.
// It takes a user ID, a multipart file header,
// and a title as input, and returns the filename of the uploaded video and an error (if any).
func UploadVideo(userId int64, data *multipart.FileHeader, title string) (filename string, err error) {
	if err = CheckVideo(data); err != nil {
		return "", err
	}
	// Generate a unique filename for the video
	// The filename is the hash of the original filename, the title, the current timestamp and a random salt.
	now := time.Now().UnixMilli()
	filename, _ = utils.HashWithSalt(data.Filename + title + strconv.FormatInt(now, 10))
	ext := utils.GetExt(data.Filename)

	if err = utils.SaveFile(data, "public/video/", filename+"."+ext); err != nil {
		return "", err
	}

	models.VideoDao().Add(&models.Video{
		AuthorId: userId,
		PlayUrl:  "/static/video/" + filename + "." + ext,
		CoverUrl: "",
		Title:    title,
	})

	return filename, nil
}

// CheckVideo checks if the given video file is valid.
//
// TODO: implement this function, we may need to use ffmpeg to check the video format.
func CheckVideo(file *multipart.FileHeader) (err error) {
	return nil
}

// GetPublishList
//
// returns a list of videos published by the given user ID
func GetPublishList(userId int64) (videos []*models.Video, err error) {
	return models.VideoDao().GetByAuthorId(userId)
}

// GetVideosBefore
//
// returns a list of videos created before the given time,
// along with the timestamp of the oldest video and an error (if any).
// The returned videos have their PlayUrl and CoverUrl fields updated with
// the current IP address and port number.
func GetVideosBefore(time int64) (videos []*models.Video, oldest int64, err error) {
	videos, oldest, err = models.VideoDao().GetBefore(time, 30)
	if err != nil {
		return nil, 0, err
	}
	ip, err := utils.GetLocalIP()
	if err != nil {
		return nil, 0, err
	}
	for _, video := range videos {
		video.PlayUrl = "http://" + ip + ":" + config.Port + video.PlayUrl
		video.CoverUrl = "http://" + ip + ":" + config.Port + video.CoverUrl
	}
	return videos, oldest, nil
}
