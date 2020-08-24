# Let Me Google That For You telegram bot with AWS Lambda's

The idea behind is very simple, instead of answering a question, you get to be a smartass and answer it through Google!

The bot is available 24/7 [@lmgt4u_bot](https://t.me/lmgt4ybot) (unless I decided to turn this off.)

Up and running on AWS Cloud and deployed with serverless framework.

## How-to use
Step 0: Setup env variables
Step 1: Deploy bot on AWS
```bash
$ make deploy
```
Step 2: Set webhook URL
```
curl --request POST --url https://api.telegram.org/bot<token>/setWebhook --header 'content-type: application/json' --data '{"url": "<URL>"}'
```


## Useful links
* https://www.serverless.com/
* https://aws.amazon.com/lambda/
* https://blog.golang.org/using-go-modules