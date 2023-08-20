package service

import (
	"encoding/json"
	"errors"
	"main/config"
	"main/models"
	"main/utils"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type ErrVideoFormat struct {
	format string
}

func (e ErrVideoFormat) Error() string {
	return "invalid video format: " + e.format
}

type VideoInfo struct {
	Id            int64       `json:"id,omitempty" gorm:"primarykey"`
	Author        UserProfile `json:"author,omitempty"`
	PlayUrl       string      `json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"`
	FavoriteCount int64       `json:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count,omitempty"`
	Title         string      `json:"title,omitempty"`
	CreatedAt     time.Time   `json:"created_at,omitempty"`
}

// UploadVideo 上传视频
//
// uploads a video file to the server and adds a new video record to the database.
// It will check the video format and extract the cover image from the video file.
// It takes a user ID, a multipart file header,
// and a title as input, and returns the filename of the uploaded video and an error (if any).
func UploadVideo(userId int64, data *multipart.FileHeader, title string) (filename string, err error) {
	// Generate a unique filename for the video
	// The filename is the hash of the original filename, the title, the current timestamp and a random salt.
	now := time.Now().UnixMilli()
	filename, _ = utils.HashWithSalt(data.Filename + title + strconv.FormatInt(now, 10))
	ext := utils.GetExt(data.Filename)

	if err = utils.SaveFile(data, "public/video/", filename+"."+ext); err != nil {
		return "", err
	}

	checkCh := make(chan bool)
	extCh := make(chan string)
	errCh := make(chan error, 2)

	go func() {
		err := CheckVideo(filename + "." + ext)
		if err != nil {
			errCh <- err
		} else {
			checkCh <- true
		}
	}()

	go func() {
		coverFilename, err := extractCover(filename + "." + ext)
		if err != nil {
			errCh <- err
		} else {
			extCh <- coverFilename
		}
	}()

	var coverFilename string

	select {
	case err := <-errCh:
		utils.RemoveFile("public/video/" + filename + "." + ext)
		return "", err
	case <-checkCh:
		coverFilename = <-extCh
	case coverFilename = <-extCh:
		<-checkCh
	}

	models.VideoDao().Add(&models.Video{
		AuthorId: userId,
		PlayUrl:  "/static/video/" + filename + "." + ext,
		CoverUrl: "/static/cover/" + coverFilename,
		Title:    title,
	})

	return filename, nil
}

// extractCover 从视频文件中提取封面
//
// extracts the first frame of a video file and saves it as a JPEG image file.
// The function takes the filename of the video file (with extension) as input and returns
// the filename of the generated cover image file (with extension) on success. If an error
// occurs during the extraction process, the function returns an error.
func extractCover(filename string) (cover string, err error) {
	// check if "public/cover" exists, if not, create it
	os.MkdirAll("public/cover", os.ModePerm)
	src := "public/video/" + filename
	now := time.Now().UnixMilli()
	targetFilename, _ := utils.HashWithSalt(filename + strconv.FormatInt(now, 10))
	target := "public/cover/" + targetFilename + ".jpg"
	if err = ffmpeg.Input(src).Output(target, ffmpeg.KwArgs{"ss": "00:00:00.000", "vframes": 1}).Run(); err != nil {
		return "", err
	}
	return targetFilename + ".jpg", nil
}

type probeInfo struct {
	Format struct {
		FormatName string `json:"format_name"`
	} `json:"format"`
	Streams []struct {
		CodecName string `json:"codec_name"`
	} `json:"streams"`
}

// CheckVideo 检查视频文件
//
// checks if the given video file is valid.
// It uses ffmpeg to probe the video file.
func CheckVideo(filename string) (err error) {
	src := "public/video/" + filename
	infoJson, err := ffmpeg.Probe(src)
	if err != nil {
		return errors.New("invalid video file")
	}
	// infoJson to struct, use json.Unmarshal
	var info probeInfo

	err = json.Unmarshal([]byte(infoJson), &info)

	if err != nil {
		return errors.New("invalid video file")
	}
	if len(info.Streams) == 0 {
		return ErrVideoFormat{info.Format.FormatName}
	}
	return nil
}

// GetPublishList 获取视频列表
//
// returns a list of videos published by the given user ID
func GetPublishList(userId int64) (videos []*models.Video, err error) {
	videos, err = models.VideoDao().GetByAuthorId(userId)
	if err != nil {
		return nil, err
	}
	if err = AdjustVideosUrl(videos); err != nil {
		return nil, err
	}
	return videos, nil
}

// GetVideosBefore 获取视频列表
//
// returns a list of videos created before the given time,
// along with the timestamp of the oldest video and an error (if any).
// The returned videos have their PlayUrl and CoverUrl fields updated with
// the current IP address and port number.
func GetVideosBefore(time int64, requestId int64) (videos []*VideoInfo, oldest int64, err error) {
	rawVideos, oldest, err := models.VideoDao().GetBefore(time, 30)
	if err != nil {
		return nil, 0, err
	}
	if err = AdjustVideosUrl(rawVideos); err != nil {
		return nil, 0, err
	}
	for _, rawVideo := range rawVideos {
		userProfile, err := GetUserProfile(rawVideo.AuthorId, requestId)
		if err != nil {
			return nil, 0, err
		}
		videos = append(videos, &VideoInfo{
			Id:        rawVideo.Id,
			Author:    *userProfile,
			PlayUrl:   rawVideo.PlayUrl,
			CoverUrl:  rawVideo.CoverUrl,
			Title:     rawVideo.Title,
			CreatedAt: rawVideo.CreatedAt,
		})
	}
	return videos, oldest, nil
}

// AdjustVideosUrl 调整视频相关URL
//
// takes a slice of *models.Video and modifies the PlayUrl and CoverUrl fields of
// each video to include the local IP address and the configured port number.
// this function is designed to fix the problem that the demo app does not support relative URLs.
func AdjustVideosUrl(videos []*models.Video) (err error) {
	ip, err := utils.GetLocalIP()
	if err != nil {
		return err
	}
	for _, video := range videos {
		video.PlayUrl = "http://" + ip + ":" + config.Port + video.PlayUrl
		video.CoverUrl = "http://" + ip + ":" + config.Port + video.CoverUrl
	}
	return nil
}
