<img align="left" width="80px" height="80px" src="https://avatars.githubusercontent.com/u/106493015?s=400&u=c6bcdab7055f741a137be13c617e9642f38faef6&v=4">

# Syslifters OffSec Tools
<a href="https://github.com/Syslifters/offsec-tools/">
    <img src="https://img.shields.io/github/stars/Syslifters/offsec-tools?color=yellow&style=flat-square">
</a>
<a href="https://github.com/Syslifters/offsec-tools/releases/latest">
    <img src="https://img.shields.io/github/v/release/Syslifters/offsec-tools?color=green&style=flat-square">
</a>
<a href="https://github.com/Syslifters/offsec-tools/releases/latest">
    <img src="https://img.shields.io/github/release-date/Syslifters/offsec-tools?color=blue&style=flat-square">
</a>
<a href="https://github.com/Syslifters/offsec-tools/releases/latest">
    <img src="https://img.shields.io/github/repo-size/Syslifters/offsec-tools?color=red&style=flat-square">
</a>
<a href="https://www.linkedin.com/company/syslifters/">
    <img src="https://img.shields.io/badge/-Linkedin-blue?style=flat-square&logo=linkedin">
</a>

This repository is intended for pentesters and red teamers using a variety of offensive security tools during their assessments. The repository is a collection of useful tools suitable for assessments in internal environments. We fetch and compile the latest version of each tool on a regular basis and provide it to you as a release.

You don't have to worry about updating and compiling the tools yourself. Just download the latest release and find all the awesome tools you will need in a single archive.

Happy Hacking! :)  
<b>Team Syslifters</b> 🦖  
<a href="https://syslifters.com">https://syslifters.com</a>
<br>
<br>
<h3 align="center">🔎 P.S. looking for a proper pentest reporting tool?</h3>
<h3 align="center">
    <a href="https://cloud.sysreptor.com/demo"><img src="/Reptor_Schild_Demo.svg" width="15%" alt="Demo"></a>
</h3>
<h3 align="center">:rocket: Have a look at <a class="md-button" href="https://docs.sysreptor.com">SysReptor</a></h3>

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
