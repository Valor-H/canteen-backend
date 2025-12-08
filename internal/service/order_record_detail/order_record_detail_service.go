package order_record_detail

import (
	"canteen/internal/model"
	"canteen/internal/repository/order_record_detail"
	"errors"
	"time"
)

type OrderRecordDetailService interface {
	GetOrderRecordDetails(startDate, endDate string) ([]model.OrderRecordDetail, error)
}

type orderRecordDetailService struct {
	repo order_record_detail.OrderRecordDetailRepository
}

func NewOrderRecordDetailService(repo order_record_detail.OrderRecordDetailRepository) OrderRecordDetailService {
	return &orderRecordDetailService{repo: repo}
}

// GetOrderRecordDetails 获取指定时间范围内的点餐详情
func (s *orderRecordDetailService) GetOrderRecordDetails(startDate, endDate string) ([]model.OrderRecordDetail, error) {
	// 验证日期格式
	if startDate == "" || endDate == "" {
		return nil, errors.New("开始日期和结束日期不能为空")
	}

	// 验证日期格式并转换
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，请使用YYYY-MM-DD格式")
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误，请使用YYYY-MM-DD格式")
	}

	// 调用仓储层查询数据
	details, err := s.repo.FindByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return details, nil
}