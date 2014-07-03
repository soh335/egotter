package main

import (
	"flag"
)

var configPath = flag.String("config", "", "config")

func main() {
	flag.Parse()

	config, err := NewConfig(*configPath)
	if err != nil {
		Fatal("config error:", err)
	}

	runner := NewRunner()
	workers := make([]Worker, 0)

	notifier := NewImKayacComNotifier(
		config.ImKayacCom.User,
		config.ImKayacCom.Password,
		config.ImKayacCom.Secret,
	)

	matchers, err := NewMatchers(config, runner.StopChan, func(msg string) {
		notifier.Notify(msg)
	})

	if err != nil {
		Fatal("matchers error:", err)
	}

	workers = append(workers, matchers)

	for _, _agent := range config.Agent {
		agent, err := NewAgent(
			config.Twitter.ConsumerKey,
			config.Twitter.ConsumerSecret,
			config.Twitter.Token,
			config.Twitter.TokenSecret,
			runner.StopChan,
			matchers.bytChan,
			_agent.Name,
			_agent.Params,
		)
		if err != nil {
			Fatal("agent error:", err)
		}
		workers = append(workers, agent)
	}

	runner.Run(workers)
}
