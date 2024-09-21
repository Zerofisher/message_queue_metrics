package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"message_queue_metrics/internal/monitor"
)

// MockStorage 现在实现 storage.Storage 接口
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveMetrics(metrics *monitor.Metrics) error {
	args := m.Called(metrics)
	return args.Error(0)
}

func (m *MockStorage) GetMetrics(startTime, endTime time.Time) ([]*monitor.Metrics, error) {
	args := m.Called(startTime, endTime)
	return args.Get(0).([]*monitor.Metrics), args.Error(1)
}

func (m *MockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestGetMetrics(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的存储
	mockStorage := new(MockStorage)

	// 创建测试数据
	testMetrics := []*monitor.Metrics{
		{
			MessageCount: 100,
			PartitionLag: map[int]int64{0: 10, 1: 20},
		},
		{
			MessageCount: 200,
			PartitionLag: map[int]int64{0: 15, 1: 25},
		},
	}

	// 设置模拟存储的预期行为
	mockStorage.On("GetMetrics", mock.Anything, mock.Anything).Return(testMetrics, nil)

	// 创建一个新的 Gin 引擎
	r := gin.New()

	// 创建一个新的 Handler 并注册路由
	handler := NewHandler(mockStorage)
	r.GET("/metrics", handler.GetMetrics)

	// 创建一个测试请求
	req, _ := http.NewRequest("GET", "/metrics?start_time=2023-05-01T00:00:00Z&end_time=2023-05-02T00:00:00Z", nil)

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 服务 HTTP 请求
	r.ServeHTTP(w, req)

	// 检查状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应体
	var response []*monitor.Metrics
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应内容
	assert.Equal(t, testMetrics, response)

	// 验证模拟存储的方法是否被调用
	mockStorage.AssertExpectations(t)
}
