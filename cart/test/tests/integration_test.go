package tests

import (
	"github.com/stretchr/testify/suite"
	testSuite "homework/cart/test/suite"
	"testing"
)

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(testSuite.DeleteCartItemSuite))
	suite.Run(t, new(testSuite.GetCartSuite))
	suite.Run(t, new(testSuite.CartCheckoutSuite))
}
