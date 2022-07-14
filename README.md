# ec2-recovery-lambda
Golang lambda for recovering EC2 instances when `Status check` are failed.
If AWS EC2 instance fails validation and the check status is in the error state, the lambda restarts it by turning it off and on. If necessary, you can enable notifications in Slack. 
# Build 
```
go build ec2-recovery-lambda.go
zip ec2-recovery-lambda.zip ec2-recovery-lambda
```
# Enable Slack Notification
### Create Slack Bot
The first thing we need to do is to create the Slack application. Visit the [slack website](https://api.slack.com/apps?new_app=1) to create the application. Select the `From scratch` option. 
You will be presented with the option to add a Name to the application and the Workspace to allow the application to be used. You should be able to see all workspaces that you are connected to. Select the appropriate workspace.
Select the Bot option.
After clicking Bots you will be redirected to a Help information page, select the option to add scopes. The first thing we need to add to the application is the actual permissions to perform anything.
After pressing `Review Scopes to Add`, scroll down to Bot scopes and start adding the 4 scopes:
```
channels:history
chat:write
chat:write.customize
incoming-webhook
```
After adding the scopes we are ready to install the application. Once you click Allow you will see long strings, one OAuth token, and one Webhook URL. Remember the location of these, or save them on another safe storage. Then we need to invite the Application into a channel that we want him to be available in.
Go there and start typing a command message which is done by starting the message with `/`. We can invite the bot by typing `/invite @NameOfYourbot`.
### Connecting Lambda to Slack
In the lambda console, navigate to the function of your interest, click on the `Configuration` tab and click on `Environment variables` from the left menu.
Click on `Edit`, then click on `Add environment variable`. Enter this variables:
```
SLACK_ENABLE=true
SLACK_AUTH_TOKEN=your_auth_token_here
SLACK_CHANNEL_ID=ec2-recovery-lambda
```
Then click on `Save`.

### Deploy
You can use `terraform` for deploying this `Lambda` function. All configuration files you can find in `terraform-deploy` folder.