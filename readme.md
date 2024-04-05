![alt text](<images/2024-04-04 11_10_48-C__Windows_System32_cmd.exe.png>)

Greetings 
-------------------

Telecord is an advanced crossplatform c2 using discord and telegram api ,it allows multi agent handling with ease

using Telegram and discord APIs can be good for exfiltration and network evasion , this project is built to enhance red teaming operations 

Telecord works by combining the two APIs to get seamless and easy interaction with each agent 

agents support mac,linux and windows

quick overview of how it works :
-------
![alt text](<images/2024-04-04 17_01_58-Telecord - Tableau blanc en ligne.png>)

the agent consist of 2 subagents , the fist is a discord bot and the other is a telegram bot.

since telegram does not allow multipleagents to run at the same time. by default the telegram bot is asleep inside our discord bot until we want to interact with it by sending the !interact command to the discord bot ,once it receives it wakes up the telegram bot ,meaning our session enabling us to execute more commands 

list of commands below:

![alt text](<images/2024-04-04 17_11_51-C__Windows_System32_cmd.exe - go  run cc.go -theme yellow.png>)


Prerequisite
-------
golang
python3
discord acc
telegram acc 

Setup
----------

git clone https://github.com/ELMERIKH/Telecord

go mod tidy

go run Telecord.go 

You will get prompted to enter 2 telegram bot tokens,1 telegram channel id ,1 Discord bot token ,a discord channel id and its webhook id

once done a config.yaml file is created with your settings ,if you want to change something either edit it or delete it and run "go run Telecord.go " another time 

about the discord and telegram setup follow :[Setup Guide](docs/SETUP.md)



⚠️ DISCLAIMER :
----------------------
ME The author takes NO responsibility and/or liability for how you choose to use any of the tools/source code/any files provided. ME The author and anyone affiliated with will not be liable for any losses and/or damages in connection with use of Telecord. By using Telecord or any files included, you understand that you are AGREEING TO USE AT YOUR OWN RISK. Once again Telecord is for EDUCATION and/or RESEARCH purposes ONLY


