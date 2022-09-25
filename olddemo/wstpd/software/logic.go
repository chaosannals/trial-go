package software

import "trial/wstpd/basis"

type SoftwareLogic struct {
}

// 导入报价单
func (i *SoftwareLogic) ImportQuotation(request *basis.SoftwareRequest) *basis.SoftwareResponse {
	// TODO
	return &basis.SoftwareResponse{
		Id:   request.Id,
		Code: 0,
		Data: "导入报价单成功",
	}
}
