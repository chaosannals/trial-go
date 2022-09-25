package basis

// 报价单项数据
type QuotationItemData struct {
	FactoryId           int     `json:"factoryId"`
	FactoryName         string  `json:"factoryName"`
	FactoryContact      string  `json:"factoryContact"`
	FactoryTelephone    string  `json:"factoryTelephone"`
	FactoryMobilephone  string  `json:"factoryMobilephone"`
	FactoryEmail        string  `json:"factoryEmail"`
	FactoryQQ           string  `json:"factoryQQ"`
	FactoryFax          string  `json:"factoryFax"`
	FactoryAddress      string  `json:"factoryAddress"`
	ToyId               int     `json:"toyId"`
	ToyNumber           string  `json:"toyNumber"`
	ToyName             string  `json:"toyName"`
	ToyNameEn           string  `json:"toyNameEn"`
	ToyArticleNumber    string  `json:"toyArticleNumber"`
	ToyCategoryId       int     `json:"toyCategoryId"`
	ToyCategory         string  `json:"toyCategory"`
	ToyCategoryEn       string  `json:"toyCategoryEn"`
	ToyUnit             string  `json:"toyUnit"`
	ToyExworksPrice     float64 `json:"toyExworksPrice"`
	ToyQuotedPrice      float64 `json:"toyQuotedPrice"`
	ToyLength           float64 `json:"toyLength"`
	ToyWidth            float64 `json:"toyWidth"`
	ToyHeight           float64 `json:"toyHeight"`
	ToyPackLength       float64 `json:"toyPackLength"`
	ToyPackWidth        float64 `json:"toyPackWidth"`
	ToyPackHeight       float64 `json:"toyPackHeight"`
	ToyPackName         string  `json:"toyPackName"`
	ToyPackNameEn       string  `json:"toyPackNameEn"`
	ToyCartonLength     float64 `json:"toyCartonLength"`
	ToyCartonWidth      float64 `json:"toyCartonWidth"`
	ToyCartonHeight     float64 `json:"toyCartonHeight"`
	ToyCartonCubicMeter float64 `json:"toyCartonCubicMeter"`
	ToyCartonCubicFeet  float64 `json:"toyCartonCubicFeet"`
	ToyCartonQuantity   float64 `json:"toyCartonQuantity"`
	ToyInboxCount       int     `json:"toyInboxCount"`
	ToyGrossWeight      float64 `json:"toyGrossWeight"`
	ToyNetWeight        float64 `json:"toyNetWeight"`
	ToyPhotoUrl         string  `json:"toyPhotoUrl"`
}

// 报价单数据
type QuotationData struct {
	Id              int                  `json:"id"`
	Type            int                  `json:"type"`
	ShowroomName    string               `json:"showroomName"`
	QuoteAt         MyTime               `json:"quoteAt"`
	CreateAt        MyTime               `json:"createAt"`
	Creator         string               `json:"creator"`
	CreatorNickname string               `json:"creatorNickname"`
	Items           []*QuotationItemData `json:"items"`
}
