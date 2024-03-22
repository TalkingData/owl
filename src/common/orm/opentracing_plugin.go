package orm

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"owl/common/global"
)

const (
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"
)

type OpentracingPlugin struct{}

func (op *OpentracingPlugin) Name() string {
	return "opentracingPlugin"
}

func before(db *gorm.DB) {
	if db.Statement.Context == nil {
		return
	}

	tracer := opentracing.GlobalTracer()

	span, _ := opentracing.StartSpanFromContext(db.Statement.Context, "gorm.DB")

	err := tracer.Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(map[string]string{}))
	if err != nil {
		return
	}

	ext.DBType.Set(span, db.Name())
	db.InstanceSet(global.OpentracingGormKey, span)
	return
}

func after(db *gorm.DB) {
	_span, isExist := db.InstanceGet(global.OpentracingGormKey)
	if !isExist {
		return
	}

	span, ok := _span.(opentracing.Span)
	if !ok {
		return
	}
	defer span.Finish()

	ext.DBStatement.Set(span, db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...))
	span.SetTag("db.rows_affected", db.Statement.RowsAffected)

	if db.Error != nil {
		ext.Error.Set(span, true)
		span.LogFields(tracerLog.Error(db.Error))
	}
}

func (op *OpentracingPlugin) Initialize(db *gorm.DB) (err error) {
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}
