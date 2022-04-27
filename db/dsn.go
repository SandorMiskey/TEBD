// region: credits
//
// Portions of this code were copied from github.com/go-sql-driver/mysql@1.6.0/dsn.go on April 20, 2022. Original credits:
//
// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2016 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.
//
// endregion: credits
// region: packages

package db

import (
	"bytes"
	"net/url"
	"sort"
	"strings"
)

// endregion: packages
// region: FormatDSN formats the given Config into a DSN string which can be passed to the driver.

func writeDSNParam(buf *bytes.Buffer, hasParam *bool, name, value string) {
	buf.Grow(1 + len(name) + 1 + len(value))
	if !*hasParam {
		*hasParam = true
		buf.WriteByte('?')
	} else {
		buf.WriteByte('&')
	}
	buf.WriteString(name)
	buf.WriteByte('=')
	buf.WriteString(value)
}

func FormatDSN(c *Config) (dsn string) {

	var buf bytes.Buffer

	// region: [username[:password]@]

	if len(c.User) > 0 {
		buf.WriteString(c.User)
		if len(c.Passwd) > 0 {
			buf.WriteByte(':')
			buf.WriteString(c.Passwd)
		}
		buf.WriteByte('@')
	}

	// endregion: [username[:password]@]
	// region: [protocol[(address)]]

	if len(c.Net) > 0 {
		buf.WriteString(c.Net)
		if len(c.Addr) > 0 {
			buf.WriteByte('(')
			buf.WriteString(c.Addr)
			buf.WriteByte(')')
		}
	}

	// endregion: [protocol[(address)]]
	// region: /dbname

	buf.WriteByte('/')
	buf.WriteString(c.DBName)

	// endregion: /dbname
	// region: [?param1=value1&...&paramN=valueN]

	hasParam := false

	if c.Params != nil {
		var ps []string
		for p := range c.Params {
			ps = append(ps, p)
		}
		sort.Strings(ps)
		for _, p := range ps {
			writeDSNParam(&buf, &hasParam, p, url.QueryEscape(c.Params[p]))
		}
	}

	// endregion: [?param1=value1&...&paramN=valueN]

	dsn = buf.String()
	c.DSN = dsn
	return dsn
}

func (c *Config) FormatDSN() {
	FormatDSN(c)
}

// endregion: FormatDSN
// region: ParseDSN parses the DSN string to a *Config

func ParseDSN(dsn string) (c Config, e error) {
	c = Config{}
	e = c.ParseDSN(dsn)
	return
}

func (c *Config) ParseDSN(ds ...string) error {
	if len(ds) == 0 {
		ds = append(ds, c.DSN)
	}
	if len(ds) > 1 {
		return ErrTooManyParameters
	}
	c.DSN = ds[0]

	// region: [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]

	// Find the last '/' (since the password or the net addr might contain a '/')
	foundSlash := false
	for i := len(c.DSN) - 1; i >= 0; i-- {
		if c.DSN[i] == '/' {
			foundSlash = true
			var j, k int

			// region: [user[:password]@][net[(addr)]]

			// left part is empty if i <= 0
			if i > 0 {

				// region: [username[:password]@]

				// Find the last '@' in dsn[:i]
				for j = i; j >= 0; j-- {
					if c.DSN[j] == '@' {
						// Find the first ':' in dsn[:j]
						for k = 0; k < j; k++ {
							if c.DSN[k] == ':' {
								c.Passwd = c.DSN[k+1 : j]
								break
							}
						}
						c.User = c.DSN[:k]
						break
					}
				}

				// endregion: [username[:password]@]
				// region: [protocol[(address)]]

				// Find the first '(' in dsn[j+1:i]
				for k = j + 1; k < i; k++ {
					if c.DSN[k] == '(' {
						// dsn[i-1] must be == ')' if an address is specified
						if c.DSN[i-1] != ')' {
							if strings.ContainsRune(c.DSN[k+1:i], ')') {
								return errInvalidDSNUnescaped
							}
							return errInvalidDSNAddr
						}
						c.Addr = c.DSN[k+1 : i-1]
						break
					}
				}
				c.Net = c.DSN[j+1 : k]

				// endregion: [protocol[(address)]]

			}

			// endregion: [user[:password]@][net[(addr)]]
			// region: dbname[?param1=value1&...&paramN=valueN]

			// Find the first '?' in dsn[i+1:]
			for j = i + 1; j < len(c.DSN); j++ {
				if c.DSN[j] == '?' {
					c.Params = make(map[string]string)
					for _, v := range strings.Split(c.DSN[j+1:], "&") {
						param := strings.SplitN(v, "=", 2)
						if len(param) != 2 {
							continue
						}
						c.Params[param[0]] = param[1]
					}
					break
				}
			}
			c.DBName = c.DSN[i+1 : j]

			// endregion: dbname[?param1=value1&...&paramN=valueN]

			break
		}
	}

	// endregion: [username[:password]@][protocol[(address)]]

	if !foundSlash && len(c.DSN) > 0 {
		return errInvalidDSNNoSlash
	}

	// if err = cfg.normalize(); err != nil {
	// 	return nil, err
	// }

	return nil
}

// endregion: ParseDSN
// region: normalize

/*

func (cfg *Config) normalize() error {
	if cfg.InterpolateParams && unsafeCollations[cfg.Collation] {
		return errInvalidDSNUnsafeCollation
	}

	// Set default network if empty
	if cfg.Net == "" {
		cfg.Net = "tcp"
	}

	// Set default address if empty
	if cfg.Addr == "" {
		switch cfg.Net {
		case "tcp":
			cfg.Addr = "127.0.0.1:3306"
		case "unix":
			cfg.Addr = "/tmp/mysql.sock"
		default:
			return errors.New("default addr for network '" + cfg.Net + "' unknown")
		}
	} else if cfg.Net == "tcp" {
		cfg.Addr = ensureHavePort(cfg.Addr)
	}

	switch cfg.TLSConfig {
	case "false", "":
		// don't set anything
	case "true":
		cfg.tls = &tls.Config{}
	case "skip-verify", "preferred":
		cfg.tls = &tls.Config{InsecureSkipVerify: true}
	default:
		cfg.tls = getTLSConfigClone(cfg.TLSConfig)
		if cfg.tls == nil {
			return errors.New("invalid value / unknown config name: " + cfg.TLSConfig)
		}
	}

	if cfg.tls != nil && cfg.tls.ServerName == "" && !cfg.tls.InsecureSkipVerify {
		host, _, err := net.SplitHostPort(cfg.Addr)
		if err == nil {
			cfg.tls.ServerName = host
		}
	}

	if cfg.ServerPubKey != "" {
		cfg.pubKey = getServerPubKey(cfg.ServerPubKey)
		if cfg.pubKey == nil {
			return errors.New("invalid value / unknown server pub key name: " + cfg.ServerPubKey)
		}
	}

	return nil
}

func ensureHavePort(addr string) string {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		return net.JoinHostPort(addr, "3306")
	}
	return addr
}

*/

// endregion: normalize
