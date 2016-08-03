package alcahuetabot

import "github.com/robfig/cron"

type Bot interface {
	Start()
	Stop()
	Run()
}

type bot struct {
	*cron.Cron
	matcher    Matcher
	matchmaker MatchMaker
}

func New(m Matcher, mm MatchMaker, cronspecs []string) (b Bot, err error) {
	c := cron.New()
	b = &bot{c, m, mm}
	for _, cronspec := range cronspecs {
		err = c.AddJob(cronspec, b)
		if err != nil {
			b = nil
		}
	}
	return
}

func (b *bot) Run() {
	match, err := b.matcher.Match()
	if err != nil {
		return
	}
	b.matchmaker.MakeMatch(match)
}
