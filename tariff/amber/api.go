package amber

import (
	"time"

	"github.com/evcc-io/evcc/util/request"
)

const ApiURI = "https://api.amber.com.au/v1/"

type Site struct {
	SiteID	string	`json:"id"`
	NMI	string	`json:"nmi"`
	Status	string	`json:"status"`
}

type PriceData struct {
	Type		string		`json:"type"`
	StartTime	time.Time	`json:"startTime"`
	EndTime		time.Time	`json:"endTime"`
	Renewables	float64		`json:"renewables"`
	Price		float64		`json:"perKwh"`
	ChannelType	string		`json:"channelType"`
}

func ApiHeader(bearerToken string) (map[string]string) {
	return map[string]string{
                        "Authorization":        "Bearer " + bearerToken,
                        "Accept":               request.JSONContent,
                }
}

func GetSiteApiUri() string {
	return ApiURI + "sites"
}

func GetCurrentPriceUri(siteID string) string {
	return GetSiteApiUri() + "/" + siteID + "/prices/current?next=1&previous=0&resolution=30"
}
