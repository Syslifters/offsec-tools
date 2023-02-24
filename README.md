<img align="left" width="80px" height="80px" src="https://user-images.githubusercontent.com/7896159/220596940-eaccedd8-979d-46eb-a014-e4082d7bf725.png">

# OffSec Tools
This repository is intended for pentesters and red teamers using a variety of offensive security tools during their assessments. The repository is a collection of useful tools suitable for assessments in internal environments. We fetch and compile the latest version of each tool on a regular basis and provide it to you as a release.

You don't have to worry about updating and compiling the tools yourself. Just download the latest release and find all the awesome tools you will need in a single archive.

Happy Hacking! :)  
<b>Team Syslifters</b> 🦖  
<a href="https://syslifters.com">https://syslifters.com</a>
<br>
<br>
<br>
<h3 align="center">🔎 PS. looking for a proper pentest reporting tool?</h3>
<h4 align="center">:rocket: Checkout <a class="md-button" href="https://docs.sysreptor.com">SysReptor</a></h4>
<h4 align="center">:bulb: Have a look at our documentation <a class="md-button" href="https://docs.sysreptor.com/">here</a></h4>
<br/>
<br>

## FAQ
### Why do you do that?
Many OffSec tools are shipped with their source code only and therefore need to be compiled manually. This is a very time-consuming task, especially if you want to keep your tools up to date before doing assessments. Better save the time for more important things, right?

### Duhh, but you also put PowerShell scripts in releases. Why?
We don't want to rack our brains every time before an assessment about which tools we need. A release conveniently contains all the tools we need for the assessment as a collection. Noice!!

### Did you backdoor the tools?
No. Cross our heart and hope to die.

### Nah, still don't trust you guys. Can't we build it ourselves?
Oh man, we don't blame you. It's the lot of the security industry. But if you're motivated, you can also create your own build pipeline. We'd be happy to show you how to do that. [Instructions](HOWTO.md)
 and our [gitlab-ci.yml](gitlab-ci.yml) are included in this repository.

### Which version of the tools do you use?
When creating a release, we use the latest version from the official repository of the respective tool. You can check the version in the commit message, which points to the latest commit in the official repository.

### How often do you plan to create a release?
We have fully automated the steps required to create a release using a build pipeline. We therefore plan to create a release once a week.

### I am missing an awesome tool, what should I do?
Just let us know! Open an issue with a link to the repository you want to add. We'll have a look and add it if it's a reasonable fit.
