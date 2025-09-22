// internal/bg/bg.go
package bg

import (
	"context"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/dyaksa/archer"
	"github.com/spf13/viper"
)

type Jobs struct{ Client *archer.Client }

func archerOptionsFromDSN(dsn string) (*archer.Options, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	var user, pass string
	if u.User != nil {
		user = u.User.Username()
		pass, _ = u.User.Password()
	}
	host := u.Host
	if !strings.Contains(host, ":") {
		host = net.JoinHostPort(host, "5432")
	}

	return &archer.Options{
		Addr:     host,
		User:     user,
		Password: pass,
		DBName:   strings.TrimPrefix(u.Path, "/"),
		SSL:      u.Query().Get("sslmode"), // forward sslmode
	}, nil
}

func NewJobs() (*Jobs, error) {
	opts, err := archerOptionsFromDSN(viper.GetString("database.dsn"))
	if err != nil {
		return nil, err
	}

	instances := viper.GetInt("archer.instances")
	if instances <= 0 {
		instances = 1
	}
	timeoutSec := viper.GetInt("archer.timeoutSec")
	if timeoutSec <= 0 {
		timeoutSec = 60
	}

	// LOG what we’re connecting to (sanitized) so you can confirm DB/host
	log.Printf("[archer] addr=%s db=%s user=%s ssl=%s", opts.Addr, opts.DBName, opts.User, opts.SSL)

	c := archer.NewClient(
		opts,
		archer.WithSetTableName("jobs"), // <- ensure correct table
		archer.WithSleepInterval(1*time.Second), // fast poll while debugging
		archer.WithErrHandler(func(err error) { // bubble up worker SQL errors
			log.Printf("[archer] ERROR: %v", err)
		}),
	)

	c.Register(
		"bootstrap_bastion",
		BastionBootstrap,
		archer.WithInstances(instances),
		archer.WithTimeout(time.Duration(timeoutSec)*time.Second),
	)

	return &Jobs{Client: c}, nil
}

func (j *Jobs) Start() error { return j.Client.Start() }
func (j *Jobs) Stop()        { j.Client.Stop() }

func (j *Jobs) Enqueue(ctx context.Context, id, queue string, args any, opts ...archer.FnOptions) (any, error) {
	return j.Client.Schedule(ctx, id, queue, args, opts...)
}
