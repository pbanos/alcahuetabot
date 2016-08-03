package alcahuetabot

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

type Match interface {
	Calixto() string
	Melibea() string
}

type match struct {
	melibea, calixto string
}

type Matcher interface {
	Match() (match Match, err error)
}

type TwitterFollowerInformer interface {
	GetFollowersIdsAll(v url.Values) (result chan anaconda.FollowersIdsPage)
	GetUsersShowById(id int64, v url.Values) (u anaconda.User, err error)
}

type MatcherDecorator func(Matcher) Matcher
type MatcherFunc func() (Match, error)

func (f MatcherFunc) Match() (Match, error) {
	return f()
}

type randomMatcher struct {
	minLength int
	tfi       TwitterFollowerInformer
}

func (rm *randomMatcher) Match() (m Match, err error) {
	var fIds []int64
	fIdsPages := rm.tfi.GetFollowersIdsAll(nil)
	for fPage := range fIdsPages {
		fIds = append(fIds, fPage.Ids...)
	}
	lenFIds := len(fIds)
	if lenFIds < rm.minLength {
		err = fmt.Errorf("Number of followers (%d) less than minimum for matching (%d)", lenFIds, rm.minLength)
		return
	}
	id := rand.Intn(lenFIds)
	fIds[id], fIds[0] = fIds[0], fIds[id]
	id = rand.Intn(lenFIds-1) + 1
	fIds[id], fIds[1] = fIds[1], fIds[id]
	melibea, err := getScreenNameForID(rm.tfi, fIds[0])
	if err != nil {
		return
	}
	calixto, err := getScreenNameForID(rm.tfi, fIds[1])
	if err != nil {
		return
	}
	m = &match{melibea, calixto}
	return
}

func NewRandomMatcher(minLength int, tfi TwitterFollowerInformer) (Matcher, error) {
	if minLength < 2 {
		return nil, fmt.Errorf("Minimum length for matching must be 2 or more")
	}
	return &randomMatcher{minLength, tfi}, nil
}

func DecorateMatcher(m Matcher, ds ...MatcherDecorator) Matcher {
	decorated := m
	for _, decorate := range ds {
		decorated = decorate(decorated)
	}
	return decorated
}

func (m *match) Melibea() string {
	return m.melibea
}

func (m *match) Calixto() string {
	return m.calixto
}

func getScreenNameForID(tfi TwitterFollowerInformer, id int64) (screenName string, err error) {
	user, err := tfi.GetUsersShowById(id, nil)
	if err != nil {
		return
	}
	screenName = user.ScreenName
	return
}
