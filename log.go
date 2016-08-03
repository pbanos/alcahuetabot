package alcahuetabot

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}

func MatcherLogging(l Logger) MatcherDecorator {
	return func(matcher Matcher) Matcher {
		return MatcherFunc(func() (m Match, err error) {
			l.Debug("Generating a match")
			m, err = matcher.Match()
			if err != nil {
				l.Error(err)
				return
			}
			l.Info(fmt.Sprintf("Match generated: %v", m))
			return
		})
	}
}

func MatchMakerLogging(l Logger) MatchMakerDecorator {
	return func(matchMaker MatchMaker) MatchMaker {
		return MatchMakerFunc(func(m Match) (t anaconda.Tweet, err error) {
			l.Debug(fmt.Sprintf("Making the match %v", m))
			t, err = matchMaker.MakeMatch(m)
			if err != nil {
				l.Info(err)
				return
			}
			l.Info(fmt.Sprintf("Match %v successfully made", m))
			return
		})
	}
}
