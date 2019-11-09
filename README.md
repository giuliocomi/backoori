# backoori
Tool aided persistence via Windows URI schemes abuse
<br/>
<a href="https://raw.githubusercontent.com/empijei/wapty/master/LICENSE" rel="nofollow"><img src="https://camo.githubusercontent.com/dcb3a3de32cb31ae6a7edf80d88747f989878809/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f6c6963656e73652d47504c76332d626c75652e737667" alt="License" data-canonical-src="https://img.shields.io/badge/license-GPLv3-blue.svg" style="max-width:100%;"></a>
<img alt="Twitter Follow" src="https://img.shields.io/twitter/follow/giulio_comi?label=Follow&style=social">

## Why
Backoori ("Backdoor the URIs") is a Proof of Concept tool aimed to automate the fileless URI persistence technique in Windows 10 targets.

## Abstract of the Research behind the tool
The widespread adoption of custom URI protocols to launch specific Windows Universal App can be diverted to a nefarious purpose. The URI schemes in Windows 10 can be abused in such a way to maintain persistence via the 'Living off the Land' approach. Backdooring a compromised Windows account in userland context is a matter of seconds. The operation is concealed to the unaware victim thanks to the URI intents being transparently proxyed to the legitimate default application.
The subtle fileless payloads can be triggered in many contexts, from the Narrator available in the Windows logon screen (an undocumented Accessibility Feature abuse technique that set off this whole research) to the classical web attack surface.

All this research started with a novel Accessibility Feature Abuse I discuss here: 

https://www.secjuice.com/abusing-windows-10-for-fileless-persistence/

The tool will be demo at BlackHat Europe Arsenal 2019:

https://www.blackhat.com/eu-19/arsenal/schedule/#backoori-tool-aided-persistence-via-windows-uri-schemes-abuse-18131

### Demo videos
1) URI persistence technique: User triggered persistence scenario
https://www.youtube.com/watch?v=oR9cSs6Sw4g

2) URI persistence technique: Hijacking multiple Universal Apps
https://www.youtube.com/watch?v=KLtDuhccfec&t=49s

3) URI persistence technique: Going beyond User triggered persistence via web 'attack surface'
https://www.youtube.com/watch?v=W6FqUx8vi5c

### Features
1) Implements the Windows 10 URI persistence technique
2) Standalone
3) 0 dependencies

### Installation
```
go get github.com/giuliocomi/backoori
go run main.go
```
#### Cross-Compile
* Windows x64: $ env GOOS=windows GOARCH=amd64 go build -o backoori main.go

* Linux x64  : $ env GOOS=linux GOARCH=amd64  go build -o backoori main.go

* MacOs x64  : $ env GOOS=darwin GOARCH=amd64  go build -o backoori main.go

[Cross-Compile Instructions](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04#step-4-%E2%80%94-building-executables-for-different-architectures)

### Usage

```
Backoori0.92: tool aided persistence via Windows URI schemes abuse
| Generate a ready-to-launch Powershell agent that will backdoor specific Universal URI Apps with fileless payloads of your choice.
|  -help
|        Display help details
|  -online string
|       Provide 'true' if wants agent to fetch the payloads via the webserver, 'false' otherwise to store the payloads directly in the agent PS file (default "false")
|  -payloads string
|       Provide the JSON filename containing the payloads to use in the backdoored gadgets (default "./resources/payloads_sample.json")
|  -protocols string
|       Provide the JSON filename containing the URI protocols to backdoor on the target system (default "./resources/uri_protocols_sample.json")
|  -psversion int
|       Provide the Powershell version that the agent will use for the payloads (recommended v2) (default 3)
```

### Example Output
* *Golang cli*

![](https://imgur.com/zAI1Rdf.png)

* *Powershell agent output*
![](https://imgur.com/jYWo83T.png)

### Possible next steps
* [ ] Adding logging and symmetric payload encryption to the web server that deploys the gadgets
* [ ] Support gadget interactions
* [ ] Convert this standalone project to a Metasploit module

## Issues
Spot a bug? Please create an issue here on GitHub (https://github.com/giuliocomi/backoori/issues)

## License
This project is licensed under the  GNU general public license Version 3.
