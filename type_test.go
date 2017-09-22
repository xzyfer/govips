package vips_test

import (
	"io/ioutil"
	"testing"

	"github.com/davidbyttow/govips"
	"github.com/stretchr/testify/assert"
)

func TestJpeg(t *testing.T) {
	if testing.Short() {
		return
	}
	buf, _ := ioutil.ReadFile("fixtures/canyon.jpg")
	assert.NotNil(t, buf)

	imageType := vips.DetermineImageType(buf)
	assert.Equal(t, vips.ImageTypeJPEG, imageType)
}
