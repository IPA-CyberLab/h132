package lws

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/IPA-CyberLab/h132/lws"
	"github.com/IPA-CyberLab/h132/pb"
)

var createCommand = &cli.Command{
	Name:  "create",
	Usage: "Create a letter writing set",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "name of the letter writing set to be created",
		},
	},
	Action: func(c *cli.Context) error {
		s := zap.S()

		name := c.String("name")
		if name == "" {
			return fmt.Errorf("name is required")
		}

		l, err := lws.ReadLWS()
		if err == nil {
			return fmt.Errorf("Letter writing set (name=%s) already exists!", l.Name)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := lws.CheckWriteAccess(lws.GetLWSDir()); err != nil {
			return err
		}

		l = &pb.LetterWritingSet{
			Name: name,
		}
		if err := lws.WriteLWS(l, os.O_CREATE|os.O_EXCL|os.O_WRONLY); err != nil {
			return err
		}

		s.Infof("Letter writing set (name=%s) successfully created!", name)
		return nil
	},
}
