General Use Case
-----------------

Themes
-------

you can change Telecord theme color:
![alt text](<../images/2024-04-05 04_51_55-C__Windows_System32_cmd.exe.png>)

!kitty
---------- 

spawn acat

help
-------------
list commands:

![alt text](<../images/2024-04-05 04_58_22-telecord.go - telegramc2 - Visual Studio Code.png>)

!generate
----------

mac - windows -linux

example:

![alt text](<../images/2024-04-05 05_04_41-USECASE.md - Telecord - Visual Studio Code.png>)

upon agent executing you will receive a notification:

![alt text](<../images/2024-04-05 05_07_52-C__Windows_System32_cmd.exe - go  run Telecord.go -theme teal.png>)

!list
----------

listing agents:

![alt text](<../images/2024-04-05 05_12_02-C__Windows_System32_cmd.exe - go  run Telecord.go -theme teal.png>)

!interact
-------------

interact with agents by id,this will get the telegram sub agent up and runing :

![alt text](<../images/2024-04-05 05_12_30-C__Windows_System32_cmd.exe - go  run Telecord.go -theme teal.png>)

/shell
--------------

execute any shell command (take in mind the platform you are on):
![alt text](<../images/2024-04-05 05_14_01-C__Windows_System32_cmd.exe - go  run Telecord.go -theme teal.png>)

/cd
-----------------

change working directory

examples: /cd .. , /cd Desktop ...

/stop
-------------

make telegram sub agent go to sleep:

![alt text](<../images/2024-04-05 05_14_29-Kali-Linux-2021.3-vmware-amd64 - VMware Workstation 17 Player (Non-commercial us.png>)

!kill
---------------

kill agent (stop process):

![alt text](<../images/2024-04-05 05_15_28-Kali-Linux-2021.3-vmware-amd64 - VMware Workstation 17 Player (Non-commercial us.png>)


/screenshot | /screenshare
---------------

will open an html file with your default browser, by default any screenshot taken willbe stored in ./screenshots folder as screenshot.jpg ,so any /screenshot command will overwrite the previous one

![alt text](<../images/2024-04-04 15_44_27-.png>)

/downlaod | /upload
-----------------
download and upload files viathe telgram api 

examples:

/download kk.txt = download kk.txt to the server

/upload kk.txt =upload kk.txt from local directory to working directory of the agent 

exit
-----------

exit telecord
