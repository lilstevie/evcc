package tariff

import (
	"errors"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/tariff/amber"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/request"
	)

type Amber struct {
	log		*util.Logger
	siteID		string
	apiKey		string
	solarFeed	bool
	data		*util.Monitor[api.Rates]
}

var _ api.Tariff = (*Amber)(nil)

func init() {
	registry.Add("amberelectric", NewAmberFromConfig)
}

func NewAmberFromConfig(other map[string]interface{}) (api.Tariff, error) {
	var cc struct {
		ApiKey		string
		SiteID		string
		NMI		string
		SolarFeed	bool
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	if cc.ApiKey == "" {
		return nil, errors.New("Missing Amber API key")
	}

	if cc.SiteID == "" && cc.NMI == ""{
		return nil, errors.New("Missing Site ID or NMI")
	}

	t := &Amber{
		log:	util.NewLogger("AmberElectric"),
		siteID:	cc.SiteID,
		apiKey:	cc.ApiKey,
		solarFeed: cc.SolarFeed,
		data:	util.NewMonitor[api.Rates](12 * time.Hour),
	}

	t.log.Redact(t.apiKey)

	if t.siteID == "" {
		t.log.Redact(cc.NMI)
		if sid, err := t.getSiteIDFromNMI(cc.NMI); err != nil {
			return nil, err
		} else {
			t.log.Redact(sid)
			t.siteID = sid
		}
	}

	done := make(chan error)
	go t.run(done)
	err := <-done

	return t, err
}

func (t *Amber) getSiteIDFromNMI(nmi string) (string, error) {

	if req, err := request.New(http.MethodGet, amber.GetSiteApiUri(), nil, amber.ApiHeader(t.apiKey)); err == nil {
		client := request.NewHelper(t.log)
		var sites []amber.Site
		if err = client.DoJSON(req, &sites); err != nil {
			return "", err
		}
		for _, site := range sites {
			if site.NMI == nmi && site.Status == "active" {
				return site.SiteID, nil
			}
		}
	} else {
		return "", err
	}

	return "", errors.New("NMI not active on this account")
}

func (t *Amber) run(done chan error) {
	var once sync.Once
	client := request.NewHelper(t.log)
	bo := newBackoff()
	request, _ := request.New(http.MethodGet, amber.GetCurrentPriceUri(t.siteID), nil, amber.ApiHeader(t.apiKey))
	for ; true ; <-time.Tick(5 * time.Minute) {
		var pd []amber.PriceData
		if err := backoff.Retry(func() error {
			return client.DoJSON(request, &pd)
		}, bo); err != nil {
			once.Do(func() { done <- err })
			t.log.ERROR.Println(err)
			continue
		}

		data := make(api.Rates, 0, len(pd))
		for _, r := range pd {
			if r.ChannelType == "general" && !t.solarFeed {
				ar := api.Rate{
					Start:	r.StartTime,
					End:	r.EndTime,
					Price:	r.Price / 1e2,
				}
				data = append(data, ar)
			} else if r.ChannelType == "feedIn" && t.solarFeed {
				ar := api.Rate{
					Start:	r.StartTime,
					End:	r.EndTime,
					Price:	r.Price / 1e2,
				}
				data = append(data, ar)
			}
		}
		data.Sort()
		t.data.Set(data)
		once.Do(func() { close(done)})
	}
}

func (t *Amber) Rates() (api.Rates, error) {
	var res api.Rates
	err := t.data.GetFunc(func (val api.Rates) {
		res = slices.Clone(val)
	})
	return res, err
}

func (t *Amber) Type() api.TariffType {
	return api.TariffTypePriceForecast
}
