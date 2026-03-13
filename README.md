Small File Exchanger

Small app in Go just to send files to other local devices. 

Created in mind to use it for 5 minutes and forgot forever (unlike FTP/SMB protocols which u have to setup and authorize). 


Created in mind of simplicity and zero-setup file share/retrieve server/client CLI app


Server can:
> Communicate in HTTP (Post/Get/Errors), Authorize (login/pass/token), Manage SQLite DB for authorizing, Send dir tree, Send files to client
Config
* Server Port
* Server DB name
* Server sharing folder

Client can:
> Communicate in HTTP, authorize and explore server, download files
Config
* Connecting IP:Port
* Login/Password
* Downloading folder


Get started:

> git clone https://github.com/KrisWilson/sfe
>
> go run main.go
>
> 
