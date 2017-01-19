package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	birc "github.com/jpiontek/bitter-irc"
)

type environment struct {
	oauth    string
	username string
	clientid string
	channels []string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	env := getEnvironment()

	if len(env.channels) == 0 {
		fmt.Println("please supply a comma separated list of channels via the -channels flag")
		os.Exit(1)
	}

	errors := make(chan error, len(env.channels))
	// Fan out to connect to channels
	for _, channel := range env.channels {
		go func(chanName string) {
			c := birc.NewTwitchChannel(chanName, env.username, env.oauth, logger, pingHandler)
			err := c.Connect()
			if err != nil {
				errors <- err
				return
			}

			err = c.Authenticate()
			if err != nil {
				errors <- err
				return
			}

			// listen method is blocking, so this will wait until it returns something
			errors <- c.Listen()
		}(channel)
		// Delay to prevent bumping up against rate limiting when connecting to multiple channels
		time.Sleep(350 * time.Millisecond)

	}

	// Fan in errors
	for i := 0; i <= len(env.channels); i++ {
		select {
		case e := <-errors:
			if e == nil {
				break
			}

			switch e.(type) {
			case *birc.ChannelError:
				ce := e.(*birc.ChannelError)
				// If we receive a ChannelError be sure to Disconnect
				ce.Channel.Disconnect()
			}

			fmt.Println(e.Error())
		}
	}
}

func getEnvironment() *environment {
	env := &environment{}

	var c string
	flag.StringVar(&env.username, "username", "", "Your Twitch username")
	flag.StringVar(&env.oauth, "oauth", "", "Your account oauth token for authentication")
	flag.StringVar(&env.clientid, "clientid", "", "Your client ID")
	flag.StringVar(&c, "channels", "", "Comma separated list of channels")
	flag.Parse()

	if c != "" {
		env.channels = strings.Split(c, ",")
	}

	return env
}

// logger is a basic digester that will print out messages to stdout as they come in. It's basically a clone
// of the birc.Logger built-in logging digester but it's here for demonstration purposes.
func logger(m birc.Message, w birc.ChannelWriter) {
	if m.Username != "" && m.Content != "" {
		fmt.Printf("%s [%s] %s\n", m.Time.Format("2006-01-06 15:04:05"), m.Username, m.Content)
	}
}

// pingHandler is a digester that will respond to a !ping command by ponging the user back.
func pingHandler(m birc.Message, w birc.ChannelWriter) {
	if m.Content == "!ping" {
		w.Send(fmt.Sprintf("Pong @%s\n", m.Username))
	}
}
