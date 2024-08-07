package controllers

import (
	"strconv"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/gin-gonic/gin"
)

func getPageParams(c *gin.Context) (dto.PageParams, error) {
	limitQueryParam := c.Query("limit")
	offsetQueryParam := c.Query("offset")

	limit, err := strconv.Atoi(limitQueryParam)
	if limitQueryParam != "" && err != nil {
		return dto.PageParams{}, err
	}

	offset, err := strconv.Atoi(offsetQueryParam)
	if offsetQueryParam != "" && err != nil {
		return dto.PageParams{}, err
	}

	return dto.NewPageParams(offset, limit), nil
}
