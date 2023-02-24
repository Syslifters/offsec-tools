# How-to Setup OffSec Tools Build Pipeline 
This how-to provides instructions for setting up the OffSec Tools Build Pipeline on your GitLab instance to automatically compile the tools of your choice for you.

## Prerequisites
* GitLab instance
* Windows host
* Remote Repository hosted on GitHub or a GitLab instance

## 1. Create a new project
Create a new project on your GitLab instance and add the [.gitlab-ci.yml](gitlab-ci.yml) to the root of your project. The only purpose of this project is to build the tools defined in the CI script and push the repositories and compiled tools to a dedicated remote repository.

## 2. Create a new personal access token
Create a personal access token for the remote repository where the repositories and compiled tools will be pushed to. See official [documentation](https://docs.github.com/en/enterprise-server@3.4/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) on how to do that for GitHub.

Scopes needed for personal access token:
- read_repository
- write_repository

## 2. Define necessary CI/CD variables
Define the following CI/CD variables in your GitLab project. They will be needed in order to successfully run the CI script.

* **CI_PROJECT_PATH**: Project path of remote repository to push the tools into (e.g. `/Syslifters/offsec-tools.git`)
* **CI_SERVER_HOST**: Hostname of the GitLab/GitHub instance (e.g. `github.com`)
* **GIT_TOKEN**: Personal access token previously created.
* **GIT_USER_EMAIL**: Username associated with the personal access token.
* **GIT_USERNAME**: Username of the user used to commit the changes.

## 3. Prepare GitLab Runner
You will need a dedicated Windows system that will compile the tools and function as GitLab Runner. 

### 3.1 Install and Register GitLab Runner
Install the GitLab Runner agent on a Windows system using the PowerShell scripts below or follow the official documentation for step-by-step instructions: [https://docs.gitlab.com/runner/install/windows.html](https://docs.gitlab.com/runner/install/windows.html)

**1. Download Runner Binary**
```
# Run PowerShell: https://docs.microsoft.com/en-us/powershell/scripting/windows-powershell/starting-windows-powershell?view=powershell-7#with-administrative-privileges-run-as-administrator
# Create a folder somewhere on your system, for example: C:\GitLab-Runner
New-Item -Path 'C:\GitLab-Runner' -ItemType Directory

# Change to the folder
cd 'C:\GitLab-Runner'

# Download binary
Invoke-WebRequest -Uri "https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-windows-amd64.exe" -OutFile "gitlab-runner.exe"
```

**2. Register Runner**
```
./gitlab-runner.exe register --url <GITLAB_INSTANCE_URL> --registration-token <REGISTRATION_TOKEN>
```

Hint: you will find the registration token in the _CI/CD Settings_ in the _Specific Runner_ section of your builder project. 

**3. Install and Start Runner-Service**
```
.\gitlab-runner.exe install
.\gitlab-runner.exe start
```


### 3.2 Edit config.toml
In the GitLab Runner installation directory, define the following settings in the _config.toml_ configuration file as follows

* **executor**: shell
* **shell**: powershell

### 3.2 Install required software
- Visual Studio 2019 or up
- Visual Studio 2013 (required by some tools)
- nuget ([https://www.nuget.org/downloads](https://www.nuget.org/downloads))

### 3.3 Set/Modify environment variables
Set/modify the following environment variables. You may need to restart the Gitlab Runner service for the changes to take effect.

* **PATH**
  1. Add installation path of _MSBuild.exe_ of Visual Studio 2019 or up. E.g. `C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\MSBuild\Current\Bin`
  2. Add installation path of _nuget.exe_

* **MSBUILD_VS13**\
Create a new environment variable _MSBUILD_VS13_ and set its value to the installation path of MSBuild.exe of Visual Studio 2013. E.g. `C:\Program Files (x86)\MSBuild\12.0\Bin`

## 4. Optional: Schedule Pipeline
Create a schedule in GitLab to repeatedly trigger the pipeline. This is very practical if you always want to have the latest version of the tools.

## FINISH 🏁
Your build pipeline should now be ready. You can manually trigger the build pipeline in the GitLab CI/CD view. This will cause the Windows GitLab Runner to retrieve the repositories in the CI script, compile them and commit everything to the remote repository you specified. Each time you trigger the pipeline, it also creates a release with all the compiled tools.

Happy Compiling and Hacking!  
Team Syslifters 🦖
