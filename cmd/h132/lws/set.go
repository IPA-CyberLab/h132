package lws

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/IPA-CyberLab/h132/cmd/h132/common"
	"github.com/IPA-CyberLab/h132/lws"
)

var setCommand = &cli.Command{
	Name:      "set",
	Usage:     "Configure a property value of the letter writing set",
	ArgsUsage: "PROPERTY_NAME [VALUE]",
	Action: func(c *cli.Context) error {
		s := zap.S()

		propertyName := c.Args().First()
		if propertyName == "" {
			return common.ErrInvalidInput{Msg: "property name is required"}
		}
		value := strings.Join(c.Args().Slice()[1:], " ")

		s.Debugf("Setting property %s=%q", propertyName, value)

		l, err := lws.ReadLWS()
		if err != nil {
			return err
		}

		if err := lws.SetProperty(l, propertyName, value); err != nil {
			return err
		}
		if err := lws.UpdateLWS(l, 0); err != nil {
			return fmt.Errorf("failed to update lws: %w", err)
		}

		s.Infof("Updated LWS property %s=%q", propertyName, value)
		return nil
	},
}
