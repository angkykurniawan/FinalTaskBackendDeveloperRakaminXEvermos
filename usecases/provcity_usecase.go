package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const emsifaBaseURL = "https://www.emsifa.com/api-wilayah-indonesia/api"

type ProvCityUsecase interface {
	GetListProvincies() (interface{}, error)
	GetListCities(provID string) (interface{}, error)
	GetDetailProvince(provID string) (interface{}, error)
	GetDetailCity(cityID string) (interface{}, error)
}

type provCityUsecase struct {
	httpClient *http.Client
}

func NewProvCityUsecase() ProvCityUsecase {
	return &provCityUsecase{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (u *provCityUsecase) fetchJSON(url string) (interface{}, error) {
	resp, err := u.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("layanan wilayah merespon dengan status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (u *provCityUsecase) GetListProvincies() (interface{}, error) {
	return u.fetchJSON(fmt.Sprintf("%s/provinces.json", emsifaBaseURL))
}

func (u *provCityUsecase) GetListCities(provID string) (interface{}, error) {
	return u.fetchJSON(fmt.Sprintf("%s/regencies/%s.json", emsifaBaseURL, provID))
}

func (u *provCityUsecase) GetDetailProvince(provID string) (interface{}, error) {
	return u.fetchJSON(fmt.Sprintf("%s/province/%s.json", emsifaBaseURL, provID))
}

func (u *provCityUsecase) GetDetailCity(cityID string) (interface{}, error) {
	return u.fetchJSON(fmt.Sprintf("%s/regency/%s.json", emsifaBaseURL, cityID))
}
