package routes

import (
	"finance-control-api/database"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
)

type ReportPeriodParams struct {
	Year  *int `query:"year" validate:"omitempty,min=1"`
	Month *int `query:"month" validate:"omitempty,min=1,max=12"`
}

type ReportLineParams struct {
	Year  *int    `query:"year" validate:"omitempty,min=1"`
	Month *int    `query:"month" validate:"omitempty,min=1,max=12"`
	Type  *string `query:"type" validate:"omitempty,oneof=income expense"`
}

type ReportRadarParams struct {
	Year     *int    `query:"year" validate:"omitempty,min=1"`
	Month    *int    `query:"month" validate:"omitempty,min=1,max=12"`
	Type     *string `query:"type" validate:"omitempty,oneof=income expense"`
	AreaIDs  string  `query:"area_ids"`
}

// GetKPIs godoc
//
//	@Summary		Relatório - KPIs
//	@Description	Retorna entradas, saídas e saldo do mês.
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Param			year	query		int	false	"Ano (padrão mês atual)"
//	@Param			month	query		int	false	"Mês (1-12, padrão mês atual)"
//	@Success		200		{object}	models.ReportKPIs
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/reports/kpis [get]
func GetKPIs(reportGroup *echo.Group) {
	reportGroup.GET("/kpis", func(c *echo.Context) error {
		userID, err := currentUserID(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": err.Error()})
		}

		var params ReportPeriodParams
		if !bindAndValidateQuery(c, &params) {
			return nil
		}

		if (params.Year == nil) != (params.Month == nil) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": []string{"Os campos year e month devem ser enviados juntos ou omitidos juntos"}})
		}

		y, m := 0, 0
		if params.Year != nil {
			y = *params.Year
		}
		if params.Month != nil {
			m = *params.Month
		}

		kpis, err := database.GetKPIs(userID, y, m)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, kpis)
	})
}

// GetLineChart godoc
//
//	@Summary		Relatório - Gráfico de linha
//	@Description	Retorna um ponto por dia do mês com valor e mediana dos últimos 6 meses.
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Param			year	query		int		false	"Ano (padrão mês atual)"
//	@Param			month	query		int		false	"Mês (1-12, padrão mês atual)"
//	@Param			type	query		string	false	"Tipo (income ou expense, padrão expense)"	Enums(income, expense)
//	@Success		200		{object}	models.LineChart
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/reports/line-chart [get]
func GetLineChart(reportGroup *echo.Group) {
	reportGroup.GET("/line-chart", func(c *echo.Context) error {
		userID, err := currentUserID(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": err.Error()})
		}

		var params ReportLineParams
		if !bindAndValidateQuery(c, &params) {
			return nil
		}

		if (params.Year == nil) != (params.Month == nil) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": []string{"Os campos year e month devem ser enviados juntos ou omitidos juntos"}})
		}

		y, m := 0, 0
		if params.Year != nil {
			y = *params.Year
		}
		if params.Month != nil {
			m = *params.Month
		}

		txType := "expense"
		if params.Type != nil {
			txType = *params.Type
		}

		chart, err := database.GetLineChart(userID, y, m, txType)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, chart)
	})
}

// GetRadarChart godoc
//
//	@Summary		Relatório - Gráfico de radar
//	@Description	Retorna valor e mediana por área ativa.
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Param			year		query		int		false	"Ano (padrão mês atual)"
//	@Param			month		query		int		false	"Mês (1-12, padrão mês atual)"
//	@Param			type		query		string	false	"Tipo (income ou expense, padrão expense)"	Enums(income, expense)
//	@Param			area_ids	query		string	false	"Lista de IDs de áreas separadas por vírgula"
//	@Success		200			{object}	models.RadarChart
//	@Failure		400			{object}	map[string]interface{}
//	@Failure		401			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]interface{}
//	@Router			/reports/radar-chart [get]
func GetRadarChart(reportGroup *echo.Group) {
	reportGroup.GET("/radar-chart", func(c *echo.Context) error {
		userID, err := currentUserID(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": err.Error()})
		}

		var params ReportRadarParams
		if !bindAndValidateQuery(c, &params) {
			return nil
		}

		if (params.Year == nil) != (params.Month == nil) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": []string{"Os campos year e month devem ser enviados juntos ou omitidos juntos"}})
		}

		y, m := 0, 0
		if params.Year != nil {
			y = *params.Year
		}
		if params.Month != nil {
			m = *params.Month
		}

		txType := "expense"
		if params.Type != nil {
			txType = *params.Type
		}

		areaIDs := make([]int, 0)
		if params.AreaIDs != "" {
			parts := strings.Split(params.AreaIDs, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				id, err := strconv.Atoi(p)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": []string{"O parâmetro area_ids deve conter apenas números separados por vírgula"}})
				}
				if id <= 0 {
					return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": []string{"O parâmetro area_ids deve conter apenas números positivos"}})
				}
				areaIDs = append(areaIDs, id)
			}
		}

		chart, err := database.GetRadarChart(userID, y, m, txType, areaIDs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, chart)
	})
}

// GetBudgetReport godoc
//
//	@Summary		Relatório - Orçamentos
//	@Description	Retorna progresso dos orçamentos do mês para áreas ativas.
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Param			year	query		int	false	"Ano (padrão mês atual)"
//	@Param			month	query		int	false	"Mês (1-12, padrão mês atual)"
//	@Success		200		{object}	models.BudgetReport
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/reports/budgets [get]
func GetBudgetReport(reportGroup *echo.Group) {
	reportGroup.GET("/budgets", func(c *echo.Context) error {
		userID, err := currentUserID(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": err.Error()})
		}

		var params ReportPeriodParams
		if !bindAndValidateQuery(c, &params) {
			return nil
		}

		if (params.Year == nil) != (params.Month == nil) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": []string{"Os campos year e month devem ser enviados juntos ou omitidos juntos"}})
		}

		y, m := 0, 0
		if params.Year != nil {
			y = *params.Year
		}
		if params.Month != nil {
			m = *params.Month
		}

		report, err := database.GetBudgetReport(userID, y, m)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, report)
	})
}
