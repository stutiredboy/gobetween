package server

/**
 * server.go - server creator
 *
 * @author Illarion Kovalchuk
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */

import (
	"errors"

	"github.com/stutiredboy/gobetween/config"
	"github.com/stutiredboy/gobetween/core"
	"github.com/stutiredboy/gobetween/server/tcp"
	"github.com/stutiredboy/gobetween/server/udp"
)

/**
 * Creates new Server based on cfg.Protocol
 */
func New(name string, cfg config.Server) (core.Server, error) {
	switch cfg.Protocol {
	case "tls", "tcp":
		return tcp.New(name, cfg)
	case "udp":
		return udp.New(name, cfg)
	default:
		return nil, errors.New("Can't create server for protocol " + cfg.Protocol)
	}
}
