package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Carlosyoot/go-estudos/config"
	"github.com/Carlosyoot/go-estudos/database"
	"github.com/gin-gonic/gin"
)

const chunkSize = 200

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

func GetHandlerStreamChunk(ctx *gin.Context) {
	if config.DB == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "DB não inicializado"})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 55*time.Second)
	defer cancel()

	rows, cols, err := database.QuerySimplesV2(reqCtx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  "falha ao executar consulta",
			"detail": err.Error(),
		})
		return
	}
	defer rows.Close()

	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Status(http.StatusOK)

	bw := bufio.NewWriterSize(ctx.Writer, 64*1024)
	defer bw.Flush()
	flusher, _ := ctx.Writer.(http.Flusher)

	_, _ = bw.Write([]byte(`{"data":[`))
	firstGlobal := true

	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	chunk := make([]map[string]any, 0, chunkSize)

	flushChunk := func() error {
		if len(chunk) == 0 {
			return nil
		}
		for i := 0; i < len(chunk); i++ {
			if !firstGlobal {
				if _, err := bw.Write([]byte{','}); err != nil {
					return err
				}
			} else {
				firstGlobal = false
			}
			b, err := json.Marshal(chunk[i])
			if err != nil {
				return err
			}
			if _, err = bw.Write(b); err != nil {
				return err
			}
		}
		chunk = chunk[:0]
		if err := bw.Flush(); err != nil {
			return err
		}
		if flusher != nil {
			flusher.Flush()
		}
		return nil
	}

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			_, _ = bw.Write([]byte(`],"error":"scan error"}`))
			bw.Flush()
			if flusher != nil {
				flusher.Flush()
			}
			return
		}

		obj := make(map[string]any, len(cols))
		for i, c := range cols {
			key := strings.ToLower(c)
			switch v := vals[i].(type) {
			case []byte:
				obj[key] = string(v)
			default:
				obj[key] = v
			}
		}

		chunk = append(chunk, obj)
		if len(chunk) >= chunkSize {
			if err := flushChunk(); err != nil {
				_, _ = bw.Write([]byte(`],"error":"flush error"}`))
				bw.Flush()
				if flusher != nil {
					flusher.Flush()
				}
				return
			}
		}
	}

	if err := rows.Err(); err != nil {
		_, _ = bw.Write([]byte(`],"error":"rows error"}`))
		bw.Flush()
		if flusher != nil {
			flusher.Flush()
		}
		return
	}

	if err := flushChunk(); err != nil {
		_, _ = bw.Write([]byte(`],"error":"flush error"}`))
		bw.Flush()
		if flusher != nil {
			flusher.Flush()
		}
		return
	}

	_, _ = bw.Write([]byte(`]}`))
	bw.Flush()
	if flusher != nil {
		flusher.Flush()
	}
}
