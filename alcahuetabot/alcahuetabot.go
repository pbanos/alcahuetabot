package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/pbanos/alcahuetabot"
	"github.com/spf13/viper"
)

// type tp struct{}
//
// func (t *tp) PostTweet(status string, v url.Values) (tweet anaconda.Tweet, err error) {
// 	fmt.Printf("Tweeting: %s\n", status)
// 	return anaconda.Tweet{}, nil
// }

type logger struct {
	level int
}

var logLevels = [4]string{"DEBUG", "INFO", "WARN", "ERROR"}

func (l *logger) messageWithLevel(level string, args ...interface{}) []interface{} {
	return append([]interface{}{level}, args...)
}

func (l *logger) Debug(args ...interface{}) {
	if l.level < 1 {
		log.Println(l.messageWithLevel("[DEBUG]", args...))
	}
}
func (l *logger) Info(args ...interface{}) {
	if l.level < 2 {
		log.Println(l.messageWithLevel("[INFO]", args...))
	}
}
func (l *logger) Warn(args ...interface{}) {
	if l.level < 3 {
		log.Println(l.messageWithLevel("[WARN]", args...))
	}
}
func (l *logger) Error(args ...interface{}) {
	if l.level < 4 {
		log.Println(l.messageWithLevel("[ERROR]", args...))
	}
}

func main() {
	readConfig()
	rand.Seed(time.Now().UTC().UnixNano())
	logLevel := -1
	slogLevel := strings.ToUpper(viper.GetString("log_level"))
	for i, level := range logLevels {
		if level == slogLevel {
			logLevel = i
			break
		}
	}
	if logLevel == -1 {
		fmt.Println("Invalid log_level ", viper.GetString("log_level"))
		os.Exit(1)
	}
	l := &logger{logLevel}
	anaconda.SetConsumerKey(viper.GetString("consumer_key"))
	anaconda.SetConsumerSecret(viper.GetString("consumer_secret"))
	api := anaconda.NewTwitterApi(viper.GetString("access_token"), viper.GetString("access_token_secret"))
	minFs := viper.GetInt("minimum_followers")
	l.Info(fmt.Sprintf("Running random matcher with a minimum number of followers of %d", minFs))
	matcher, err := alcahuetabot.NewRandomMatcher(minFs, api)
	if err != nil {
		l.Error(err)
		os.Exit(1)
	}

	matcher = alcahuetabot.DecorateMatcher(matcher, alcahuetabot.MatcherLogging(l))
	// introductions := []string{"Creo que %{calixto} y %{melibea} harían muy buena pareja!! ♡♡♡"}
	// matchmaker, err := alcahuetabot.NewMatchMaker(api, introductions)
	introductions := viper.GetStringSlice("introductions")
	l.Info(fmt.Sprintf("Running matchmaker with %d introductions", len(introductions)))
	matchmaker, err := alcahuetabot.NewMatchMaker(api, introductions)
	if err != nil {
		l.Error(err)
		os.Exit(1)
	}
	matchmaker = alcahuetabot.DecorateMatchMaker(matchmaker, alcahuetabot.MatchMakerLogging(l))
	cronspecs := viper.GetStringSlice("cronspecs") //[]string{"0 * * * * *"}
	l.Info(fmt.Sprintf("Bot will run on the following schedule: %v", cronspecs))
	b, err := alcahuetabot.New(matcher, matchmaker, cronspecs)
	if err != nil {
		l.Error(err)
		os.Exit(1)
	}
	l.Info("Starting...")
	go b.Start()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	l.Warn(fmt.Sprintf("Stopping because of %v...", <-ch))
	b.Stop()
	l.Info("Exiting...")
	os.Exit(0)
}

func readConfig() {
	viper.SetDefault("minimum_followers", 2)
	viper.SetDefault("log_level", "info")
	viper.AutomaticEnv()
	viper.SetConfigName("config")              // name of config file (without extension)
	viper.AddConfigPath("/etc/alcahuetabot/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.alcahuetabot") // call multiple times to add many search paths
	viper.AddConfigPath(".")                   // optionally look for config in the working directory
	err := viper.ReadInConfig()                // Find and read the config file
	if err != nil {                            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
