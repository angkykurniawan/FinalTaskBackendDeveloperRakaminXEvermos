package controllers

import (
	"evermos-mini/usecases"
	"evermos-mini/utils"

	"github.com/gofiber/fiber/v2"
)

type ProvCityController struct {
	provCityUsecase usecases.ProvCityUsecase
}

func NewProvCityController(provCityUsecase usecases.ProvCityUsecase) *ProvCityController {
	return &ProvCityController{provCityUsecase: provCityUsecase}
}

// GET /provcity/listprovincies
func (ctrl *ProvCityController) GetListProvincies(c *fiber.Ctx) error {
	data, err := ctrl.provCityUsecase.GetListProvincies()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data provinsi")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "", data)
}

// GET /provcity/listcities/:prov_id
func (ctrl *ProvCityController) GetListCities(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	data, err := ctrl.provCityUsecase.GetListCities(provID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data kota")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "", data)
}

// GET /provcity/detailprovince/:prov_id
func (ctrl *ProvCityController) GetDetailProvince(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	data, err := ctrl.provCityUsecase.GetDetailProvince(provID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil detail provinsi")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "", data)
}

// GET /provcity/detailcity/:city_id
func (ctrl *ProvCityController) GetDetailCity(c *fiber.Ctx) error {
	cityID := c.Params("city_id")
	data, err := ctrl.provCityUsecase.GetDetailCity(cityID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil detail kota")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "", data)
}
