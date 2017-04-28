package handlers

// Context is the context of the application
import (
	"github.com/clebi/gofin/es"
	finance "github.com/clebi/yfinance"
	elastic "gopkg.in/olivere/elastic.v5"
)

// SchemaDecoder decodes URL query to struct
type SchemaDecoder interface {
	Decode(dst interface{}, src map[string][]string) error
}

//Context contains resources that needs to be access in http handlers
type Context struct {
	es         *elastic.Client
	sh         SchemaDecoder
	historyAPI finance.HistoryAPI
	esStock    es.IEsStock
	esPosition es.IEsPositionStock
}

//NewContext creates a new context for handlers
func NewContext(
	es *elastic.Client,
	sh SchemaDecoder,
	historyAPI finance.HistoryAPI,
	esStock es.IEsStock,
	esPosition es.IEsPositionStock) *Context {
	return &Context{
		es:         es,
		sh:         sh,
		historyAPI: historyAPI,
		esStock:    esStock,
		esPosition: esPosition,
	}
}
