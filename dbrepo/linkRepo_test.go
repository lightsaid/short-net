package dbrepo

import (
	"testing"
	"time"

	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/util"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createLink(t *testing.T, userIds ...uint) models.Link {
	longUrl := "http://" + util.RandomString(10)
	hash := util.RandomString(6)
	expired_at := time.Now().Add(time.Minute)

	var user models.User
	if len(userIds) > 0 {
		user.ID = userIds[0]
	} else {
		user = createUser(t)
	}

	link := models.Link{
		UserID:    user.ID,
		LongURL:   longUrl,
		ShortHash: hash,
		ExpiredAt: expired_at,
	}

	err := testRepo.CreateLink(&link)
	require.NoError(t, err)
	require.True(t, link.ID > 0)
	require.Equal(t, longUrl, link.LongURL)
	require.Equal(t, hash, link.ShortHash)
	require.Equal(t, uint(0), link.Click)
	require.WithinDuration(t, time.Now(), link.CreatedAt, time.Second)
	require.WithinDuration(t, time.Now(), link.UpdatedAt, time.Second)
	require.WithinDuration(t, expired_at, link.ExpiredAt, time.Second)

	return link
}

func TestCreateLink(t *testing.T) {
	_ = createLink(t)
}

func TestGetLinkByID(t *testing.T) {
	link := createLink(t)

	qLink, err := testRepo.GetLinkByID(link.ID)
	require.NoError(t, err)

	require.True(t, qLink.ID > 0)
	require.Equal(t, link.LongURL, qLink.LongURL)
	require.Equal(t, link.ShortHash, qLink.ShortHash)
	require.Equal(t, link.Click, qLink.Click)
	require.WithinDuration(t, link.CreatedAt, qLink.CreatedAt, time.Second)
	require.WithinDuration(t, link.UpdatedAt, qLink.UpdatedAt, time.Second)
	require.WithinDuration(t, link.ExpiredAt, qLink.ExpiredAt, time.Second)
}

func TestGetLinkByHash(t *testing.T) {
	link := createLink(t)

	qLink, err := testRepo.GetLinkByHash(link.ShortHash)
	require.NoError(t, err)

	require.True(t, qLink.ID > 0)
	require.Equal(t, link.LongURL, qLink.LongURL)
	require.Equal(t, link.ShortHash, qLink.ShortHash)
	require.Equal(t, link.Click, qLink.Click)
	require.WithinDuration(t, link.CreatedAt, qLink.CreatedAt, time.Second)
	require.WithinDuration(t, link.UpdatedAt, qLink.UpdatedAt, time.Second)
	require.WithinDuration(t, link.ExpiredAt, qLink.ExpiredAt, time.Second)
}

func TestUpdateLinkByID(t *testing.T) {
	link := createLink(t)
	long_url := util.RandomString(10)

	err := testRepo.UpdateLinkByID(link.ID, models.Link{LongURL: long_url})
	require.NoError(t, err)

	qLink, err := testRepo.GetLinkByID(link.ID)
	require.NoError(t, err)
	require.True(t, qLink.ID > 0)

	require.Equal(t, long_url, qLink.LongURL)
	require.Equal(t, link.ShortHash, qLink.ShortHash)
	require.Equal(t, link.Click, qLink.Click)
	require.WithinDuration(t, link.CreatedAt, qLink.CreatedAt, time.Second)
	require.WithinDuration(t, link.UpdatedAt, qLink.UpdatedAt, time.Second)
	require.WithinDuration(t, link.ExpiredAt, qLink.ExpiredAt, time.Second)
}

func TestDeleteLinkByID(t *testing.T) {
	link := createLink(t)

	err := testRepo.DeleteLinkByID(link.ID)
	require.NoError(t, err)

	_, err = testRepo.GetLinkByID(link.ID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestListLinksByUserID(t *testing.T) {
	user := createUser(t)
	count := 5

	for i := 0; i < count; i++ {
		longUrl := "http://" + util.RandomString(10)
		hash := util.RandomString(6)
		expired_at := time.Now().Add(time.Minute)

		link := models.Link{
			UserID:    user.ID,
			LongURL:   longUrl,
			ShortHash: hash,
			ExpiredAt: expired_at,
		}
		err := testRepo.CreateLink(&link)
		require.NoError(t, err)
	}

	lists, err := testRepo.ListLinksByUserID(user.ID, Filters{Page: 0, Size: count})
	require.NoError(t, err)
	require.True(t, len(lists) == count)

	for _, l := range lists {
		require.NotEmpty(t, l)
		require.True(t, len(l.LongURL) > 0)
		require.True(t, len(l.ShortHash) > 0)
	}
}
