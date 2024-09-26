package tests

import (
	"github.com/stretchr/testify/suite"
	testSuite "homework/loms/test/suite"
	"testing"
)

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(testSuite.OrderCreateSuite))
	suite.Run(t, new(testSuite.OrderPaySuite))
	suite.Run(t, new(testSuite.OrderCancelSuite))
	suite.Run(t, new(testSuite.StockInfoSuite))
}
