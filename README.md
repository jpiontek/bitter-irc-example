# Bitter IRC Example Application

This is an example application that utilizes the [bitter-irc](https://github.com/jpiontek/bitter-irc) library to build a small "bot".

## Running

To run execute the following command in the repository root after cloning to your local machine.

```
go run main.go -oauth=token -clientid=id -username=user -channels=example,channels
```

As you can see you need to pass in a few flags:

* oauth
* clientid
* usrename
* channels

These are pretty self-explanatory. "oauth" is your oauth token. "channels" is a comma separated list of values to join multiple channels.

## Explanation

The application will connect to multiple Twitch IRC channels and start pumping the messages to stdout via rhe logger digester. 
It also has a small reactive digester, named pingHandler, that will respond to a user that executes the !ping command in a channel.

When running you'll see log messages in your terminal that look something like this:

> 2017-01-10 12:32:10 [user] Kappa 

And an example of the pingHandler:

> 2017-01-10 12:32:10 [user] !ping  
> 2017-01-10 12:32:10 [bot] pong @user

## Notes

All digesters **must** be threadsafe. They will be called from multiple go routines as new messages come in. Keep this in mind 
if you create a digester that interacts with any type of unsafe struct. 

