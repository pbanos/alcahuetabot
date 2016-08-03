# alcahuetabot

Matchmaker twitter bot written in Go. It periodically chooses two of its followers and introduces one another posting a tweet and mentioning both of them on it.

## Config format

Uses [viper](https://github.com/spf13/viper) underneath. Needs following values

* consumer_key - The consumer key for the Twitter app the bot will be using
* consumer_secret - The consumer secret for the Twitter app the bot will be using
* access_token - The access token for the Twitter user the bot will be managing
* access_token_secret - The access token secret for the Twitter user the bot will be managing
* minimum_followers - The minimum number of followers required to make a match, must be greater or equal to 2. Set to 2 by default, you should probably set this to a higher value to avoid frequent repetition of the same match.
* log_level - The level of logging. It defaults to debug and must be one of debug, info, warn or error.
* cronspecs - A collection of cron-style interval strings that determine when the bot must awake and make matches. The syntax for this specs can be found [here](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format).

  For example, "0 * * * * *" configures the bot to try and make a match once a minute.

* introductions - A collection of strings that serve as templates for the tweets that will introduce the matches. The templates should have the strings '%{calixto}' and '%{melibea}' as placeholders for the twitter users to introduce.

  For example, given the introduction template 'I think %{calixto} should ask %{melibea} out' and the twitter users user1 and user2, the bot would tweet:
  ```I think @user1 should ask @user2 out```.

It looks for a config file in the following locations:
* the /etc/alcahuetabot directory
* an .alcahuetabot directory in the user's home
* the working directory
