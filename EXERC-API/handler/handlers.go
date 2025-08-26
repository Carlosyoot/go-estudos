package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Carlosyoot/go-estudos/config"
	"github.com/Carlosyoot/go-estudos/database"
	"github.com/gin-gonic/gin"
)

func GetHandler(ctx *gin.Context) {

	if config.DB == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "DB não inicializado"})
		return
	}

	data, err := database.QuerySimples()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  "falha ao buscar dados",
			"detail": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func GetHandlerStream(ctx *gin.Context) {
	if config.DB == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "DB não inicializado"})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 55*time.Second)
	defer cancel()

	rows, cols, err := database.QuerySimplesV2(reqCtx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao executar consulta", "detail": err.Error()})
		return
	}
	defer rows.Close()

	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Status(http.StatusOK)

	w := ctx.Writer
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	_, _ = w.Write([]byte(`{"data":[`))

	first := true
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			_, _ = w.Write([]byte(`],"error":"scan error"}`))
			return
		}

		obj := make(map[string]any, len(cols))
		for i, c := range cols {
			key := strings.ToLower(c)

			if b, ok := vals[i].([]byte); ok {
				obj[key] = string(b)
			} else {
				obj[key] = vals[i]
			}
		}

		if !first {
			_, _ = w.Write([]byte(","))
		}
		first = false

		if err := enc.Encode(&obj); err != nil {
			_, _ = w.Write([]byte(`],"error":"encode error"}`))
			return
		}
	}

	if err := rows.Err(); err != nil {
		_, _ = w.Write([]byte(`],"error":"rows error"}`))
		return
	}

	_, _ = w.Write([]byte(`]}`))

}
