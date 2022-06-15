package log

import (
	"log"
	"os"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/pio"
	"github.com/peashoot/sunlight/config"
)

func Init(app *iris.Application) {
	log.SetOutput(config.AppRunLogFileWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	app.Logger().Level = golog.DebugLevel
	debugSetting()
	infoSetting()
	warnSetting()
	errorSetting()
	fatalSetting()
	app.Logger().TimeFormat = "2006-01-02 15:04:05.999"
	app.Logger().SetOutput(os.Stdout).SetOutput(config.AppRunLogFileWriter)
}

func debugSetting() {
	level := golog.Levels[golog.DebugLevel]
	level.Name = "debug"         // default
	level.Title = "[DBUG]"       // default
	level.ColorCode = pio.Yellow // default
}

func infoSetting() {
	level := golog.Levels[golog.InfoLevel]
	level.Name = "info"        // default
	level.Title = "[INFO]"     // default
	level.ColorCode = pio.Blue // default
}

func warnSetting() {
	level := golog.Levels[golog.WarnLevel]
	level.Name = "warn"           // default
	level.Title = "[WARN]"        // default
	level.ColorCode = pio.Magenta // default
}

func errorSetting() {
	level := golog.Levels[golog.ErrorLevel]
	level.Name = "error"      // default
	level.Title = "[ERRO]"    // default
	level.ColorCode = pio.Red // default
}

func fatalSetting() {
	level := golog.Levels[golog.FatalLevel]
	level.Name = "fatal"      // default
	level.Title = "[FTAL]"    // default
	level.ColorCode = pio.Red // default
}
