package dbrepo

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/util"
	"github.com/stretchr/testify/require"
)

func createBook(t *testing.T) models.Book {
	book := models.Book{
		Title:   util.RandomString(8),
		Price:   uint(util.RandomInt(2, 10)),
		Stock:   uint(util.RandomInt(5, 10)),
		Picture: fmt.Sprintf("%s%s", "http://", util.RandomString(8)),
	}

	book2 := book

	err := testRepo.CreateBook(&book2)
	require.NoError(t, err)

	require.Equal(t, book.Title, book2.Title)
	require.Equal(t, book.Price, book2.Price)
	require.Equal(t, book.Stock, book2.Stock)
	require.Equal(t, book.Picture, book2.Picture)

	return book2
}

func TestGetBook(t *testing.T) {
	book := createBook(t)
	book2, err := testRepo.GetBook(book.ID)
	require.NoError(t, err)

	require.Equal(t, book, book2)
}

func TestCreateBook(t *testing.T) {
	createBook(t)
}

// NOTE: 测试下单图书，并发问题
func TestOrderBook(t *testing.T) {
	var wg sync.WaitGroup

	book := createBook(t)

	fmt.Println("init >>", book.Stock)

	var listUsers []models.User
	for i := 0; i < 25; i++ {
		user := createUser(t)
		listUsers = append(listUsers, user)
	}

	wg.Add(len(listUsers))
	for _, user := range listUsers {
		// 模拟其他业务在查询 图书
		go func() {
			_, err := testRepo.GetBook(book.ID)
			if err != nil {
				log.Printf("GetBook err: %v\n", err)
			} else {
				log.Println("其他业务查询成功")
			}
		}()

		go func(user models.User) {
			defer wg.Done()
			err := testRepo.TxUserBuyBook(user.ID, book.ID)
			if err != nil {
				log.Printf("TxUserBuyBook err: %v - %t\n", err, errors.Is(err, ErrUnderstock))
			}
		}(user)
	}

	wg.Wait()
	book2, err := testRepo.GetBook(book.ID)
	require.NoError(t, err)
	require.NotEqual(t, book.Stock, book2.Stock)
	require.Equal(t, book2.Stock, uint(0))

	fmt.Println("last >>", book2.Stock)
}
