# Small File Exchanger

## Small fast app in Go just to send files to other local devices with zero-setup. 

Created in mind to use it for 5 minutes and forgot forever (unlike FTP/SMB protocols which u have to setup and authorize). 


![Client-side image](https://github.com/KrisWilson/sfe/blob/master/client.png)

Created in mind of simplicity and zero-setup file share/retrieve server/client CLI app


Server can:
> Communicate in HTTP (Post/Get/Put/Errors), Authorize (login/pass/token), Manage SQLite DB for authorizing, Send dir tree in JSON, Send files/dirs to client, Receive files/dirs from clients

Do some configs like:
* Server Port
* Server DB name
* Server Sharing folder

Client can:
> Communicate in HTTP, Do multiple downloads (even thousands of them at one time), Authorize and explore server, Download files and directories, Upload files and also directories, All from CLI

Also do some configs
* Connecting IP:Port
* Login/Password
* Downloading folder


Tested platforms: i386/amd64 for Linux and Windows

Get started:
```
git clone https://github.com/KrisWilson/sfe
go get sfe/cmd
go get sfe/client
go get sfe/listener
go run main.go
```

![Client-side Help CLI](https://github.com/KrisWilson/sfe/blob/master/help.png)







TODO:
* Redis for authorization and the most downloaded data from share folder
* Load-Balancer and quene manager, maximimum download/upload config, number of maximum connections
* Add option for HTTP to send/download files in parts for big files and low RAM machines


