package app

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/IPA-CyberLab/h132/cmd/h132/common"
	"github.com/IPA-CyberLab/h132/cmd/h132/envelope"
	"github.com/IPA-CyberLab/h132/cmd/h132/keys"
	lwsc "github.com/IPA-CyberLab/h132/cmd/h132/lws"
	"github.com/IPA-CyberLab/h132/lws"
	"github.com/IPA-CyberLab/h132/version"
)

func New() *cli.App {
	app := cli.NewApp()
	app.Name = "h132"
	app.Usage = "File encryption tool"
	app.Authors = []*cli.Author{
		{Name: "yzp0n", Email: "yzp0n@coe.ad.jp"},
	}
	app.Version = fmt.Sprintf("%s.%s", version.Version, version.Commit)
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "Enable verbose output",
		},
	}

	app.Commands = []*cli.Command{
		envelope.Command,
		keys.Command,
		lwsc.Command,
	}
	BeforeImpl := func(c *cli.Context) error {
		var logger *zap.Logger
		if loggeri, ok := app.Metadata["Logger"]; ok {
			logger = loggeri.(*zap.Logger)
		} else {
			cfg := zap.NewProductionConfig()
			if c.Bool("verbose") {
				cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			}
			cfg.DisableCaller = true
			cfg.EncoderConfig.TimeKey = ""
			cfg.Encoding = "console"

			var err error
			logger, err = cfg.Build(
				zap.AddStacktrace(zap.NewAtomicLevelAt(zap.DPanicLevel)))
			if err != nil {
				return err
			}
		}
		zap.ReplaceGlobals(logger)

		if err := os.Setenv("H132_LWS_DIR", lws.GetLWSDir()); err != nil {
			return fmt.Errorf("failed to set H132_LWS_DIR envvar: %w", err)
		}

		return nil
	}
	app.Before = func(c *cli.Context) error {
		if err := BeforeImpl(c); err != nil {
			// Print error message to stderr
			app.Writer = app.ErrWriter

			// Suppress help message on app.Before() failure.
			cli.HelpPrinter = func(_ io.Writer, _ string, _ interface{}) {}
			return err
		}

		return nil
	}
	app.ExitErrHandler = func(c *cli.Context, err error) {
		if errors.Is(err, common.ErrInvalidInput{}) {
			cli.ShowSubcommandHelp(c)
		}
		cli.HandleExitCoder(err)
	}

	app.After = func(c *cli.Context) error {
		zap.L().Sync()
		return nil
	}

	return app
}
