package core_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/sukhajata/portsapi/service/internal/core"
	"github.com/sukhajata/portsapi/service/mocks"
	pb "github.com/sukhajata/portsapi/service/pkg/proto"
	"testing"
)

func setup(mockCtrl *gomock.Controller) (*core.Service, *mocks.MockSQLEngine){
	mockSQLEngine := mocks.NewMockSQLEngine(mockCtrl)
	
	service := core.NewService(mockSQLEngine)
	return service, mockSQLEngine
}

func TestService_GetPort(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	service, mockDBEngine := setup(mockCtrl)
	
	// arrange
	retVal := `
	{
		"name": "Abu Dhabi",
		"coordinates": [
		  54.37,
		  24.47
		],
		"city": "Abu Dhabi",
		"province": "Abu Z¸aby [Abu Dhabi]",
		"country": "United Arab Emirates",
		"alias": ["Test"],
		"regions": ["Asia", "Middle East"],
		"timezone": "Asia/Dubai",
		"unlocs": [
		  "AEAUH"
		],
		"code": "52001"
	}`
	req := &pb.GetPortRequest {
		Id: "ABCD",
	}
	
	// expect
	mockDBEngine.EXPECT().QueryTextColumn(gomock.Any(), req.Id).Return([]string{retVal}, nil)
	
	// execute
	port, err := service.GetPort(req)
	
	// assert
	require.NoError(t, err)
	require.Equal(t, "52001", port.Code)
}

// TODO - tests for UpsertPorts and GetPorts