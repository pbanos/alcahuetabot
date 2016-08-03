package alcahuetabot

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

type MatchMaker interface {
	MakeMatch(m Match) (tweet anaconda.Tweet, err error)
}

type TweetPoster interface {
	PostTweet(status string, v url.Values) (tweet anaconda.Tweet, err error)
}

type MatchMakerDecorator func(MatchMaker) MatchMaker
type MatchMakerFunc func(Match) (anaconda.Tweet, error)

func (f MatchMakerFunc) MakeMatch(m Match) (anaconda.Tweet, error) {
	return f(m)
}

type matchMaker struct {
	tweetPoster   TweetPoster
	introductions []string
}

func NewMatchMaker(tweetPoster TweetPoster, introductions []string) (mm MatchMaker, err error) {
	if len(introductions) < 1 {
		return nil, fmt.Errorf("Cannot build match maker without at least one introduction template")
	}
	return &matchMaker{tweetPoster, introductions}, nil
}

func DecorateMatchMaker(m MatchMaker, ds ...MatchMakerDecorator) MatchMaker {
	decorated := m
	for _, decorate := range ds {
		decorated = decorate(decorated)
	}
	return decorated
}

func (mm *matchMaker) MakeMatch(m Match) (anaconda.Tweet, error) {
	introduction := mm.introductions[rand.Intn(len(mm.introductions))]
	introduction = strings.Replace(introduction, "%{melibea}", fmt.Sprintf("@%s", m.Melibea()), 1)
	introduction = strings.Replace(introduction, "%{calixto}", fmt.Sprintf("@%s", m.Calixto()), 1)
	return mm.tweetPoster.PostTweet(introduction, nil)
}
