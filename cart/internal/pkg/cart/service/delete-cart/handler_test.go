package delete_cart

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"homework/cart/internal/pkg/cart/service/delete-cart/mock"
	"testing"
)

func TestHandler_DeleteItemsByUser(t *testing.T) {
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)

	addItemHandler := NewHandler(cartRepositoryMock)

	userId := int64(123)

	cartRepositoryMock.DeleteItemsByUserMock.Expect(userId).Return(&userId, nil)

	actualResponse, err := addItemHandler.DeleteItemsByUser(userId)

	expectedResponse := &userId
	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}

func TestHandler_DeleteItemsByUser_WithError(t *testing.T) {
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)

	addItemHandler := NewHandler(cartRepositoryMock)

	userId := int64(0)

	_, err := addItemHandler.DeleteItemsByUser(userId)

	require.Error(t, err)
}
