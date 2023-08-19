package service

import (
	"io"
	"main/models"
	"mime/multipart"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func TestUploadVideo(t *testing.T) {
	defer os.RemoveAll("../service/public/")
	defer os.RemoveAll("./service/public/")
	// Open the test video file
	file, err := os.Open("../public/video/test.mp4")
	assert.NoError(t, err)
	defer file.Close()

	// Create a temporary file to hold the multipart form data
	tempFile, err := os.CreateTemp("", "test.mp4")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write the multipart form data to the temporary file
	writer := multipart.NewWriter(tempFile)
	part, err := writer.CreateFormFile("file", "test.mp4")
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	// Reset the file pointer to the beginning
	_, err = tempFile.Seek(0, io.SeekStart)
	assert.NoError(t, err)

	// Parse the multipart form data
	r := multipart.NewReader(tempFile, writer.Boundary())
	form, err := r.ReadForm(10 << 20)
	assert.NoError(t, err)

	// Get the file header from the form data
	fileHeader := form.File["file"][0]

	// Mock the VideoDao.Add method
	gomonkey.ApplyMethod(reflect.TypeOf(models.VideoDao()), "Add", func(dao *models.VideoDaoStruct, video *models.Video) (*models.Video, error) {
		assert.Equal(t, int64(1), video.AuthorId)
		assert.Equal(t, "test video", video.Title)
		assert.True(t, strings.HasPrefix(video.PlayUrl, "/static/video/"))
		assert.True(t, strings.HasPrefix(video.CoverUrl, "/static/cover/"))
		return video, nil
	})
	filename, err := UploadVideo(1, fileHeader, "test video")

	assert.NoError(t, err)
	assert.NotEmpty(t, filename)
}

func TestUploadVideoWithInvalidFile(t *testing.T) {
	// Create an invalid file with no content
	tempFile, err := os.CreateTemp("", "test.mp4")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// write some content to the file

	_, err = tempFile.WriteString(`--boundary
Content-Disposition: form-data; name="file"; filename="test.mp4"
Content-Type: video/mp4

Some content but not a valid video file
--boundary--`)
	assert.NoError(t, err)

	// Reset the file pointer to the beginning
	_, err = tempFile.Seek(0, io.SeekStart)
	assert.NoError(t, err)

	// Parse the multipart form data
	r := multipart.NewReader(tempFile, "boundary")
	form, err := r.ReadForm(10 << 20)
	assert.NoError(t, err)
	filename, err := UploadVideo(1, form.File["file"][0], "test video")

	assert.Error(t, err)
	assert.Empty(t, filename)
}
