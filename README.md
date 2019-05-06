# Ginko
Telegram bot written in Go. Currently WIP. Contributions are welcome.

A modular group management bot, written with the sole purpose of being entirely in Go, for the concurrency and speed that it offers.

You can [find me on telegram](https://t.me/sphericalkat)! I'm usually online, so I can hopefully answer any questions you may have.

## Setting up the bot (Important! Please go through once):

This project is entirely dockerized, so you don't have to go through the hassle of setting up dependencies. All you need is Docker, which you can [install here](https://docs.docker.com/install/).

### Configuration
There are two possible ways of configuring your bot. Through environment variables, or a dotenv file. 

The preferred method is to create a dotenv file named `.env`, as it makes it much easier to see all your configuration settings grouped together. A sample dotenv file called `.sample_env` has been included for convenience.

If you can't, or don't want to use a dotenv file, it is also possible to use environment variables. The following environment variables are supported:

* `BOT_API_KEY` :  Your bot token, as a string
* `BOT_NAME` : The name of your bot, as it appears on telegram
* `OWNER_USERNAME` : Your Telegram username, without the `@`
* `OWNER_ID` : Your Telegram ID
* `SUDO_USERS`: A list of userIDs, separated by spaces, who should have sudo access to the bot

**Note: As of now, all the above fields are required**

## Starting the bot
Once you have docker installed, and created your dotenv file or populated the neccessary environment variables, you can start up your bot. All you need to do is execute two commands.

Run: 
```
docker-compose build
```
to make sure any updates to the bot are reflected in the image.

And then run,
```
docker-compose up
```
to start your bot.

Alternatively, you may use 
```
docker-compose up -d
```
to run the bot in the background.




