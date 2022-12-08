package app

import (
	"reflect"
	"time"

	"gormdemo/util"

	_ "github.com/revel/modules"
	"github.com/revel/revel"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	AppVersion string   // 版本
	BuildTime  string   // 构建时间
	Db         *gorm.DB //

	// 新日志实例，没有打印
	Log = revel.RootLog.New("for", "format")
)

func init() {
	Log.SetStackDepth(1)

	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	if r, ok := c.AppController.(util.Injectable); ok {
		t := r.GetType()
		revel.AppLog.Infof("controller type: %s", t.Name())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Type.Kind() == reflect.Pointer {
				revel.AppLog.Infof("%d %s => %s", i, f.Name, "pointer")
			} else {
				revel.AppLog.Infof("%d %s => %s", i, f.Name, f.Type.Name())
			}
		}
		revel.AppLog.Info("====================")
		for i := 0; i < t.NumMethod(); i++ {
			revel.AppLog.Infof("%d %s", i, t.Method(i).Name)
		}

		v := r.GetValue()
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if f.Kind() == reflect.Pointer {
				e := f.Elem()
				if e.IsValid() {
					revel.AppLog.Infof("vp %d => %v", i, e.Type())
				}
			} else {
				revel.AppLog.Infof("v %d => %s", i, f.Type().Name())
			}
		}

		revel.AppLog.Info("====================")
		pv := r.GetPointerValue()
		if inject := pv.MethodByName("Inject"); true {
			arg := make([]reflect.Value, 1)
			arg[0] = reflect.ValueOf(12345)
			inject.Call(arg)
			revel.AppLog.Info("call")
		}
	}

	// revel.AppLog.Info(c.Type.Name())
	// revel.AppLog.Info(c.Type.Namespace)
	// revel.AppLog.Info(c.Type.ShortName())
	// revel.AppLog.Info(c.MethodType.Name)

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func InitDB() {
	var err error
	// driver := revel.Config.StringDefault("db.driver", "mysql")
	cs := revel.Config.StringDefault("db.connect", "root:password@tcp(127.0.0.1:3306)/reveldemo?charset=utf8mb4&parseTime=True&loc=Local")
	if Db, err = gorm.Open(mysql.Open(cs), &gorm.Config{}); err != nil {
		revel.AppLog.Error("Db error", err)
		return
	}

	pool := dbresolver.Register(dbresolver.Config{}).SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(100).
		SetMaxOpenConns(200)
	Db.Use(pool)
	if revel.DevMode == true {

	}
}
