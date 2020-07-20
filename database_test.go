package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
)

var (
	repoTest Repository
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		repoTest, err = NewRepository("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	defer func() {
		repoTest.Close()
	}()

	err = repoTest.Up()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	var meta = ogMetaData{
		site:          "https://www.facebook.com",
		ogDescription: "some description",
		ogHeight:      500,
		ogWidth:       500,
		ogImage:       "https://www.someimage.com",
	}
	err := repoTest.Create(meta)
	loadMeta, err := repoTest.FindBySite(meta.site)
	assert.NoError(t, err)
	assert.Equal(t, meta.site, loadMeta.site)
	assert.Equal(t, meta.ogDescription, loadMeta.ogDescription)
	assert.Equal(t, meta.ogImage, loadMeta.ogImage)
	assert.Equal(t, meta.ogHeight, loadMeta.ogHeight)
	assert.Equal(t, meta.ogWidth, loadMeta.ogWidth)
}

func TestUpdate(t *testing.T) {
	var meta = ogMetaData{
		site:          "https://www.facebook.com",
		ogDescription: "alternative",
		ogHeight:      0,
		ogWidth:       0,
		ogImage:       "https://www.someimage.com/imagery",
	}
	err := repoTest.Update(meta)
	assert.NoError(t, err)
	loadMeta, err := repoTest.FindBySite(meta.site)
	assert.NoError(t, err)
	assert.Equal(t, meta.site, loadMeta.site)
	assert.Equal(t, meta.ogDescription, loadMeta.ogDescription)
	assert.Equal(t, meta.ogImage, loadMeta.ogImage)
	assert.Equal(t, meta.ogHeight, loadMeta.ogHeight)
	assert.Equal(t, meta.ogWidth, loadMeta.ogWidth)
}
