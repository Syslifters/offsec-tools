function Invoke-DNSUpdate
{
    <#
    .SYNOPSIS
    This function performs secure and nonsecure DNS dynamic updates against an AD domain controller. Authentication
    for secure updates is performed through Kerberos GSS-TSIG. 

    Author: Kevin Robertson (@kevin_robertson)
    License: BSD 3-Clause  
    
    .DESCRIPTION
    This function can be used to add/delete dynamic DNS records through secure or nonsecure dynamic updates against an
    AD domain controller. A, AAAA, CNAME, MX, PTR, SRV, and TXT records are currently supported. Invoke-DNSUpdate is modeled
    after BIND`s nsupdate tool when using the '-g' or 'gsstsig' options for secure updates or no authentication for
    nonsecure updates. 

    By default, Active Directory-integrated zones have secure dynamic updates enabled with authenticated users having
    'Create all child objects' permission. Records that do not exist in an AD zone can be added/deleted with a standard
    user account. Existing records created by default or created by other users impose limitations. For example, creating
    records that apply to the root of the zone or creating additional SRV records for kerberos/ldap will likely be blocked
    due to existing records. Note however that older existing dynamic records can sometimes be hijacked. Subdomain folders
    can also be created.

    With secure dynamic updates, this function supports only GSS-TSIG through Kerberos AES256-CTS-HMAC-SHA1-96 using
    two separate methods. By default, the function will have Windows perform all Kerberos steps up until the AP-REQ
    is sent to DNS on the DC. This method will work with either the current session context or with specified credentials.
    The second method performs Kerberos authentication using just PowerShell code over a TCPClient connection. This method
    will accept a password or AES256 hash and will not place any tickets in the client side cache.

    In the event that a zone is configured for nonsecure dynamic updates, you should have full control over the zone.

    Note that wpad and isatap are on a block list by default starting with Server 2008. Although the records can be added
    with both secure and nonsecure dynamic updates, AD DNS will not answer requests for wpad and isatap if they are listed
    on the block list. 

    .PARAMETER Credential
    PSCredential object that will be used to to perform dynamic update.

    .PARAMETER DomainController
    Domain controller to target in FQDN format.

    .PARAMETER Realm
    Kerberos realm.

    .PARAMETER Username
    Username of user with DNS secure dynamic update access. If using a machine account, the trailing '$' must be
    included if one is shown in the SAMAccountName attribute. Note, this is an alternative to using Credential and
    is mainly included as part of pass the hash functionality.

    .PARAMETER Password
    Password of user with DNS secure dynamic update access. The password must be in the form of a secure string.
    Note, this is an alternative to using Credential and is mainly included as part of pass the hash
    functionality.

    .PARAMETER Hash
    AES256 password hash for user with DNS secure dynamic update access. Note that this will use Kerberos
    authentication built on top of TCPClient. Note, this is an alternative to using Credential and is mainly
    included as part of pass the hash functionality.

    .PARAMETER Security
    Default = Secure: (Auto/Nonsecure/Secure) Dynamic update security type. Auto will attempt to use nonsecure. If
    nonsecure fails, secure will be used. This is the standard dynamic update behavior. Secure is the default
    because it generates less traffic with default setups. 

    .PARAMETER DNSName
    DNS record name.

    .PARAMETER DNSData
    DNS records data. For most record types this will be the destination hostname or IP address. For TXT records
    this can be used for data. If deleting a record, do not set this parameter.

    .PARAMETER DNSPort
    DNS SRV record port.

    .PARAMETER DNSPreference
    DNS MX record preference.

    .PARAMETER DNSPriority
    DNS SRV record priority.

    .PARAMETER DNSTTL
    Default = 600: DNS record TTL.

    .PARAMETER DNSType
    Default = A: DNS record type. This function supports A, AAAA, CNAME, MX, PTR, SRV, and TXT.

    .PARAMETER DNSWeight
    DNS SRV record weight.

    .PARAMETER DNSZone
    DNS zone.

    .PARAMETER RecordCheck
    Check for an existing matching record before attempting to add or delete.

    .PARAMETER TCPClientAuth
    Switch to force usage of the TCPClient based Kerberos authentication. Note, usernames are case sensitive with
    this switch.

    .EXAMPLE
    Invoke-DNSUpdate -DNSName www -DNSData 192.168.100.125
    Add A Record

    .EXAMPLE
    Invoke-DNSUpdate -DNSType AAAA -DNSName www.test.local -DNSData 2001:0db8:85a3:0000:0000:8a2e:0370:7334
    Add AAAA Record

    .EXAMPLE
    Invoke-DNSUpdate -DNSType CNAME -DNSName www.test.local -DNSData system.test.local
    Add CNAME Record

    .EXAMPLE
    Invoke-DNSUpdate -DNSType MX -DNSName test.local -DNSData 192.168.100.125 -DNSPreference 10
    Add MX Record

    .EXAMPLE
    Invoke-DNSUpdate -DNSType PTR -DNSName 125.100.168.192.in-addr.arpa -DNSData www.test.local -DNSZone 100.168.192.in-addr.arpa
    Add PTR Record - there is a good chance this will be denied if there is an existing record for an IP

    .EXAMPLE
    Invoke-DNSUpdate -DNSType SRV -DNSName _autodiscover._tcp.lab.local -DNSData system.test.local -DNSPriority 100 -DNSWeight 80 -DNSPort 443
    Add SRV Record

    .EXAMPLE
    Invoke-DNSUpdate -DNSType TXT -DNSName host.test.local -DNSData "some text"
    Add TXT Record

    .EXAMPLE
    Invoke-DNSUpdate -DNSType TXT -DNSName host.test.local
    Delete TXT record - all deletes follow the same format, just specify DNSType and DNSName

    .EXAMPLE
    Invoke-DNSUpdate -DNSType A -DNSName www.test.local -Username testuser
    Add A record using another account

    .EXAMPLE
    Invoke-DNSUpdate -DNSType A -DNSName www.test.local -Username testuser -Hash 0C27E0A5B0D69640B40DDED4A28EB3BB0D157659EBED2816A41A8228E98D111B
    Add A record using another account and an AES256 hash

    .LINK
    https://github.com/Kevin-Robertson/Powermad
    #>

    [CmdletBinding()]
    param
    (
        [parameter(Mandatory=$false)][String]$DomainController,
        [parameter(Mandatory=$false)][String]$Realm,
        [parameter(Mandatory=$false)][String]$Username,
        [parameter(Mandatory=$false)][System.Security.SecureString]$Password,
        [parameter(Mandatory=$false)][ValidateScript({$_.Length -eq 64})][String]$Hash,
        [parameter(Mandatory=$false)][String]$Zone,
        [parameter(Mandatory=$false)][Int]$DNSTTL = 600,
        [parameter(Mandatory=$false)][Int]$DNSPreference,
        [parameter(Mandatory=$false)][Int]$DNSPriority,
        [parameter(Mandatory=$false)][Int]$DNSWeight,
        [parameter(Mandatory=$false)][Int]$DNSPort,
        [parameter(Mandatory=$false)][ValidateSet("Auto","Nonsecure","Secure")][String]$Security = "Secure",
        [parameter(Mandatory=$false)][ValidateSet("A","AAAA","CNAME","MX","PTR","SRV","TXT")][String]$DNSType = "A",
        [parameter(Mandatory=$true)][String]$DNSName,
        [parameter(Mandatory=$false)][ValidateScript({$_.Length -le 255})][String]$DNSData,
        [parameter(Mandatory=$false)][Switch]$RecordCheck,
        [parameter(Mandatory=$false)][Switch]$TCPClientAuth,
        [parameter(Mandatory=$false)][System.Management.Automation.PSCredential]$Credential,
        [parameter(ValueFromRemainingArguments=$true)]$invalid_parameter
    )

    if($invalid_parameter)
    {
        Write-Output "[-] $($invalid_parameter) is not a valid parameter"
        throw
    }

    if($TCPClientAuth -and (!$Credential -and !$Username))
    {
        Write-Output "[-] TCPClientAuth requires a username or PSCredential"
        throw
    }

    switch ($DNSType)
    {

        'MX'
        {

            if(!$DNSPreference)
            {
                Write-Output "[-] MX records require a DNSPreference"
                throw
            }

        }

        'PTR'
        {

            if(!$Zone)
            {
                Write-Output "[-] PTR records require a DNSZone"
                throw
            }

        }

        'SRV'
        {

            if(!$DNSPriority -and !$DNSWeight -and !$DNSPort -and $DNSData)
            {
                Write-Output "[-] DNSType SRV requires DNSPriority, DNSWeight, and DNSPort"
                throw
            }
    
            if($DNSName -notlike '*._tcp.*' -and $DNSName -notlike '*._udp.*')
            {
                Write-Output "[-] DNSName doesn't contain a protocol"
                throw
            }

        }
        
    }

    if($Security -ne 'Nonsecure' -and $Username -and !$Hash)
    {
        $password = Read-Host -Prompt "Enter password" -AsSecureString  
    }

    if(!$DistinguishedName -and (!$DomainController -or !$Domain -or !$Realm -or !$Zone))
    {

        try
        {
            $current_domain = [System.DirectoryServices.ActiveDirectory.Domain]::GetCurrentDomain()
        }
        catch
        {
            Write-Output "[-] $($_.Exception.Message)"
            throw
        }

    }

    if(!$DomainController)
    {
        $DomainController = $current_domain.PdcRoleOwner.Name
        Write-Verbose "[+] Domain Controller = $DomainController"
    }

    if(!$Domain)
    {
        $Domain = $current_domain.Name
        Write-Verbose "[+] Domain = $Domain"
    }

    if(!$Realm)
    {
        $Realm = $current_domain.Name
        Write-Verbose "[+] Kerberos Realm = $Realm"
    }

    if(!$Zone)
    {
        $Zone = $current_domain.Name
        Write-Verbose "[+] DNS Zone = $Zone"
    }

    $Zone = $Zone.ToLower()

    if($TCPClientAuth -or $Hash -or ($TCPClientAuth -and $Credential))
    {
        $kerberos_tcpclient = $true
        $Realm = $Realm.ToUpper()

        if($Credential)
        {
            $Username = $Credential.Username
            $Password = $Credential.Password
        }

        if($username -like "*\*")
        {
            $username = $username.SubString(($username.IndexOf("\") + 1),($username.Length - ($username.IndexOf("\") + 1)))
        }

        if($username -like "*@*")
        {
            $username = $username.SubString(0,($username.IndexOf("@")))
        }

        if($Username.EndsWith("$"))
        {
            $salt = $Realm + "host" + $Username.SubString(0,$Username.Length - 1) + "." + $Realm.ToLower()        
        }
        else
        {
            $salt = $Realm + $Username    
        }

        Write-Verbose "[+] Salt $salt"
    }
    
    function ConvertFrom-PacketOrderedDictionary($OrderedDictionary)
    {

        ForEach($field in $OrderedDictionary.Values)
        {
            $byte_array += $field
        }

        return [Byte[]]$byte_array
    }

    function Get-KerberosAES256UsageKey
    {
        param([String]$key_type,[Int]$usage_number,[Byte[]]$base_key)

        $padding = 0x00 * 16

        if($key_type -eq 'checksum')
        {
            switch($usage_number) 
            {
                25 {[Byte[]]$usage_constant = 0x5d,0xfb,0x7d,0xbf,0x53,0x68,0xce,0x69,0x98,0x4b,0xa5,0xd2,0xe6,0x43,0x34,0xba + $padding}
            }
        }
        elseif($key_type -eq 'encrypt')
        {

            switch($usage_number) 
            {
                1 {[Byte[]]$usage_constant = 0xae,0x2c,0x16,0x0b,0x04,0xad,0x50,0x06,0xab,0x55,0xaa,0xd5,0x6a,0x80,0x35,0x5a + $padding}
                3 {[Byte[]]$usage_constant = 0xbe,0x34,0x9a,0x4d,0x24,0xbe,0x50,0x0e,0xaf,0x57,0xab,0xd5,0xea,0x80,0x75,0x7a + $padding}
                4 {[Byte[]]$usage_constant = 0xc5,0xb7,0xdc,0x6e,0x34,0xc7,0x51,0x12,0xb1,0x58,0xac,0x56,0x2a,0x80,0x95,0x8a + $padding}
                7 {[Byte[]]$usage_constant = 0xde,0x44,0xa2,0xd1,0x64,0xe0,0x51,0x1e,0xb7,0x5b,0xad,0xd6,0xea,0x80,0xf5,0xba + $padding}
                11 {[Byte[]]$usage_constant = 0xfe,0x54,0xaa,0x55,0xa5,0x02,0x52,0x2f,0xbf,0x5f,0xaf,0xd7,0xea,0x81,0x75,0xfa + $padding}
                12 {[Byte[]]$usage_constant = 0x05,0xd7,0xec,0x76,0xb5,0x0b,0x53,0x33,0xc1,0x60,0xb0,0x58,0x2a,0x81,0x96,0x0b + $padding}
            }
                
        }
        elseif($key_type -eq 'integrity') 
        {
            
            switch($usage_number) 
            {
                1 {[Byte[]]$usage_constant = 0x5b,0x58,0x2c,0x16,0x0a,0x5a,0xa8,0x05,0x56,0xab,0x55,0xaa,0xd5,0x40,0x2a,0xb5 + $padding}
                4 {[Byte[]]$usage_constant = 0x72,0xe3,0xf2,0x79,0x3a,0x74,0xa9,0x11,0x5c,0xae,0x57,0x2b,0x95,0x40,0x8a,0xe5 + $padding}
                7 {[Byte[]]$usage_constant = 0x8b,0x70,0xb8,0xdc,0x6a,0x8d,0xa9,0x1d,0x62,0xb1,0x58,0xac,0x55,0x40,0xeb,0x15 + $padding}
                11 {[Byte[]]$usage_constant = 0xab,0x80,0xc0,0x60,0xaa,0xaf,0xaa,0x2e,0x6a,0xb5,0x5a,0xad,0x55,0x41,0x6b,0x55 + $padding}
            }

        }

        $AES = New-Object "System.Security.Cryptography.AesManaged"
        $AES.Mode = [System.Security.Cryptography.CipherMode]::CBC
        $AES.Padding = [System.Security.Cryptography.PaddingMode]::Zeros
        $AES.IV = 0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00
        $AES.KeySize = 256
        $AES.Key = $base_key
        $AES_encryptor = $AES.CreateEncryptor()
        $usage_key = $AES_encryptor.TransformFinalBlock($usage_constant,0,$usage_constant.Length)

        return $usage_key
    }

    # TCPClient Kerberos start - this section can be removed if not using a hash or -TCPClientAuth
    function Get-KerberosAES256BaseKey
    {
        param([String]$salt,[System.Security.SecureString]$password)

        $password_BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($password)
        $password_cleartext = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($password_BSTR)
        [Byte[]]$salt = [System.Text.Encoding]::UTF8.GetBytes($salt)
        [Byte[]]$password_cleartext = [System.Text.Encoding]::UTF8.GetBytes($password_cleartext)
        $constant = 0x6B,0x65,0x72,0x62,0x65,0x72,0x6F,0x73,0x7B,0x9B,0x5B,0x2B,0x93,0x13,0x2B,0x93,0x5C,0x9B,0xDC,0xDA,0xD9,0x5C,0x98,0x99,0xC4,0xCA,0xE4,0xDE,0xE6,0xD6,0xCA,0xE4
        $PBKDF2 = New-Object Security.Cryptography.Rfc2898DeriveBytes($password_cleartext,$salt,4096)
        Remove-Variable password_cleartext
        $PBKDF2_key = $PBKDF2.GetBytes(32)
        $AES = New-Object "System.Security.Cryptography.AesManaged"
        $AES.Mode = [System.Security.Cryptography.CipherMode]::CBC
        $AES.Padding = [System.Security.Cryptography.PaddingMode]::None
        $AES.IV = 0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00
        $AES.KeySize = 256
        $AES.Key = $PBKDF2_key
        $AES_encryptor = $AES.CreateEncryptor()
        $base_key_part_1 = $AES_encryptor.TransformFinalBlock($constant,0,$constant.Length)
        $base_key_part_2 = $AES_encryptor.TransformFinalBlock($base_key_part_1,0,$base_key_part_1.Length)
        $base_key = $base_key_part_1[0..15] + $base_key_part_2[0..15]

        return $base_key
    }

    function New-PacketKerberosASREQ()
    {
        param([Byte[]]$Username,[Byte[]]$Realm,[Byte[]]$NameString,[Byte[]]$Nonce,[Byte[]]$PAC,[Byte[]]$PACSignature)

        $timestamp = Get-Date
        $till = $timestamp.AddYears(20)
        $timestamp = ("{0:u}" -f $timestamp) -replace "-","" -replace " ","" -replace ":",""
        $till = ("{0:u}" -f $till) -replace "-","" -replace " ","" -replace ":",""
        [Byte[]]$timestamp = [System.Text.Encoding]::UTF8.GetBytes($timestamp)
        [Byte[]]$till = [System.Text.Encoding]::UTF8.GetBytes($till)

        if($PAC)
        {
            $pac_extra_length = 78
        }

        [Byte[]]$namestring1_length = Get-ASN1LengthArray $NameString.Count
        [Byte[]]$namestring_length = Get-ASN1LengthArray ($NameString.Count + $namestring1_length.Count + 6)
        [Byte[]]$namestring_length2 = Get-ASN1LengthArray ($NameString.Count + $namestring1_length.Count + $namestring_length.Count + 7)
        [Byte[]]$sname_length = Get-ASN1LengthArray ($NameString.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + 13)
        [Byte[]]$sname_length2 = Get-ASN1LengthArray ($NameString.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + 14)
        [Byte[]]$realm_length = Get-ASN1LengthArray $Realm.Count
        [Byte[]]$realm_length2 = Get-ASN1LengthArray ($Realm.Count + $realm_length.Count + 1)
        [Byte[]]$cname_length = Get-ASN1LengthArray $Username.Count
        [Byte[]]$cname_length2 = Get-ASN1LengthArray ($Username.Count + $cname_length.Count + 1)
        [Byte[]]$cname_length3 = Get-ASN1LengthArray ($Username.Count + $cname_length.Count + $cname_length2.Count + 2)
        [Byte[]]$cname_length4 = Get-ASN1LengthArray ($Username.Count + $cname_length.Count + $cname_length2.Count + $cname_length3.Count + 8)
        [Byte[]]$cname_length5 = Get-ASN1LengthArray ($Username.Count + $cname_length.Count + $cname_length2.Count + $cname_length3.Count + $cname_length4.Count + 9)
        $grouped_length = $address_length.Count + $address_length2.Count + $address_length3.Count + $address_length4.Count + $address_length5.Count + $NameString.Count +
            $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + $sname_length2.Count + $Realm.Count + $realm_length.Count +
            $realm_length2.Count + $Username.Count + $cname_length.Count + $cname_length2.Count + $cname_length3.Count + $cname_length4.Count + $cname_length5.Count
        [Byte[]]$reqbody_length = Get-ASN1LengthArrayLong ($grouped_length + 86)
        [Byte[]]$reqbody_length2 = Get-ASN1LengthArrayLong ($grouped_length + $reqbody_length.Count + 87)
        [Byte[]]$message_length = Get-ASN1LengthArrayLong ($grouped_length + $reqbody_length.Count + $reqbody_length2.Count + $pac_extra_length + 114)
        [Byte[]]$message_length2 = Get-ASN1LengthArrayLong ($grouped_length + $reqbody_length.Count + $reqbody_length2.Count + $message_length.Count + $pac_extra_length + 115)
        [Byte[]]$asreq_length = [System.BitConverter]::GetBytes($grouped_length + $reqbody_length.Count + $reqbody_length2.Count + $message_length.Count + $message_length2.Count +
            $pac_extra_length + 116)[3..0]

        $KerberosASREQ = New-Object System.Collections.Specialized.OrderedDictionary
        $KerberosASREQ.Add("Length",$asreq_length)
        $KerberosASREQ.Add("Message_Encoding",[Byte[]](0x6a) + $message_length2 + [Byte[]](0x30) + $message_length)
        $KerberosASREQ.Add("Message_PVNO_Encoding",[Byte[]](0xa1,0x03,0x02,0x01))
        $KerberosASREQ.Add("Message_PVNO",[Byte[]](0x05))
        $KerberosASREQ.Add("Message_MSGType_Encoding",[Byte[]](0xa2,0x03,0x02,0x01))
        $KerberosASREQ.Add("Message_MSGType",[Byte[]](0x0a))

        if($PAC)
        {
            $KerberosASREQ.Add("Message_PAData_Encoding",[Byte[]](0xa3,0x5c,0x30,0x5a,0x30,0x4c,0xa1,0x03,0x02,0x01,0x02))
            $KerberosASREQ.Add("Message_PAData0_Type_Encoding",[Byte[]](0xa2,0x45,0x04,0x43,0x30,0x41,0xa0,0x03,0x02,0x01))
            $KerberosASREQ.Add("Message_PAData0_Type",[Byte[]](0x12))
            $KerberosASREQ.Add("Message_PAData0_Value_Encoding",[Byte[]](0xa2,0x3a,0x04,0x38))
            $KerberosASREQ.Add("Message_PAData0_Value",$PAC)
            $KerberosASREQ.Add("Message_PAData0_Signature",$PACSignature)
            $KerberosASREQ.Add("Message_PAData1_Type_Encoding",[Byte[]](0x30,0x0a,0xa1,0x04,0x02,0x02))
        }
        else
        {
            $KerberosASREQ.Add("Message_PAData_Encoding",[Byte[]](0xa3,0x0e,0x30,0x0c,0x30,0x0a))
            $KerberosASREQ.Add("Message_PAData1_Type_Encoding",[Byte[]](0xa1,0x04,0x02,0x02))
        }

        $KerberosASREQ.Add("Message_PAData1_Type",[Byte[]](0x00,0x95))
        $KerberosASREQ.Add("Message_PAData1_Value_Encoding",[Byte[]](0xa2,0x02,0x04))
        $KerberosASREQ.Add("Message_PAData1_Value",[Byte[]](0x00))
        $KerberosASREQ.Add("Message_REQBody_Encoding",[Byte[]](0xa4) + $reqbody_length2 + [Byte[]](0x30) + $reqbody_length)
        $KerberosASREQ.Add("Message_REQBody_KDCOptions_Encoding",[Byte[]](0xa0,0x07,0x03,0x05))
        $KerberosASREQ.Add("Message_REQBody_KDCOptions_Padding",[Byte[]](0x00))
        $KerberosASREQ.Add("Message_REQBody_KDCOptions",[Byte[]](0x50,0x00,0x00,0x00))
        $KerberosASREQ.Add("Message_REQBody_CName_Encoding",[Byte[]](0xa1) + $cname_length5 + [Byte[]](0x30) + $cname_length4)
        $KerberosASREQ.Add("Message_REQBody_CName_NameType_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosASREQ.Add("Message_REQBody_CName_NameType",[Byte[]](0x01))
        $KerberosASREQ.Add("Message_REQBody_CName_NameString_Encoding",[Byte[]](0xa1) + $cname_length3 + [Byte[]](0x30) + $cname_length2 + [Byte[]](0x1b) + $cname_length)
        $KerberosASREQ.Add("Message_REQBody_CName_NameString",$Username)
        $KerberosASREQ.Add("Message_REQBody_Realm_Encoding",[Byte[]](0xa2) + $realm_length2 + [Byte[]](0x1b) + $realm_length)
        $KerberosASREQ.Add("Message_REQBody_Realm",$Realm)
        $KerberosASREQ.Add("Message_REQBody_SName_Encoding",[Byte[]](0xa3) + $sname_length2 + [Byte[]](0x30) + $sname_length)
        $KerberosASREQ.Add("Message_REQBody_SName_NameType_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosASREQ.Add("Message_REQBody_SName_NameType",[Byte[]](0x01))
        $KerberosASREQ.Add("Message_REQBody_SName_NameString_Encoding",[Byte[]](0xa1) + $namestring_length2 + [Byte[]](0x30) + $namestring_length)
        $KerberosASREQ.Add("Message_REQBody_SName_NameString0_Encoding",[Byte[]](0x1b,0x03))
        $KerberosASREQ.Add("Message_REQBody_SName_NameString0",[Byte[]](0x44,0x4e,0x53))
        $KerberosASREQ.Add("Message_REQBody_SName_NameString1_Encoding",[Byte[]](0x1b) + $namestring1_length) #50
        $KerberosASREQ.Add("Message_REQBody_SName_NameString1",$NameString)
        $KerberosASREQ.Add("Message_REQBody_Till_Encoding",[Byte[]](0xa5,0x11,0x18,0x0f))
        $KerberosASREQ.Add("Message_REQBody_Till",$till)
        $KerberosASREQ.Add("Message_REQBody_Nonce_Encoding",[Byte[]](0xa7,0x06,0x02,0x04))
        $KerberosASREQ.Add("Message_REQBody_Nonce",$Nonce)
        $KerberosASREQ.Add("Message_REQBody_EType_Encoding",[Byte[]](0xa8,0x15,0x30,0x13))
        $KerberosASREQ.Add("Message_REQBody_EType",[Byte[]](0x02,0x01,0x12,0x02,0x01,0x11,0x02,0x01,0x17,0x02,0x01,0x18,0x02,0x02,0xff,0x79,0x02,0x01,0x03))

        return $KerberosASREQ
    }

    function New-PacketKerberosAPREQ()
    {
        param([Byte[]]$Realm,[Byte[]]$SPN,[Byte[]]$KVNO,[Byte[]]$Ticket,[Byte[]]$Authenticator,[Byte[]]$AuthenticatorSignature)

        $Authenticator += $AuthenticatorSignature
        $parameter_length = $Realm.Count + $SPN.Count + $ticket.Count + $Authenticator.Count
        [Byte[]]$authenticator_length = Get-ASN1LengthArrayLong $Authenticator.Count
        [Byte[]]$authenticator_length2 = Get-ASN1LengthArrayLong ($Authenticator.Count + $authenticator_length.Count + 1)
        [Byte[]]$Authenticator_length3 = Get-ASN1LengthArrayLong ($authenticator.Count + $authenticator_length.Count + $authenticator_length2.Count + 7)
        [Byte[]]$authenticator_length4 = Get-ASN1LengthArrayLong ($Authenticator.Count + $authenticator_length.Count + $authenticator_length2.Count + $authenticator_length3.Count + 8)
        [Byte[]]$ticket_length = Get-ASN1LengthArrayLong $ticket.Count
        [Byte[]]$ticket_length2 = Get-ASN1LengthArrayLong ($ticket.Count + $ticket_length.Count + 1)
        [Byte[]]$ticket_length3 = Get-ASN1LengthArrayLong ($ticket.Count + $ticket_length.Count + $ticket_length2.Count + 12)
        [Byte[]]$ticket_length4 = Get-ASN1LengthArrayLong ($ticket.Count + $ticket_length.Count + $ticket_length2.Count  + $ticket_length3.Count + 13)
        [Byte[]]$namestring1_length = Get-ASN1LengthArray $SPN.Count
        [Byte[]]$namestring_length = Get-ASN1LengthArray ($SPN.Count + $namestring_length.Count + 4)
        [Byte[]]$namestring_length2 = Get-ASN1LengthArray ($SPN.Count + $namestring1_length.Count + $namestring_length.Count + 5)
        [Byte[]]$sname_length = Get-ASN1LengthArray ($SPN.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + 4)
        [Byte[]]$sname_length2 = Get-ASN1LengthArray ($SPN.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + 5)
        [Byte[]]$sname_length3 = Get-ASN1LengthArray ($SPN.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + $sname_length2.Count + 11)
        [Byte[]]$sname_length4 = Get-ASN1LengthArray ($SPN.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + $sname_length2.Count +
            $sname_length3.Count + 12)
        [Byte[]]$realm_length = Get-ASN1LengthArray $Realm.Count
        [Byte[]]$realm_length2 = Get-ASN1LengthArray ($Realm.Count + $realm_length.Count + 1)
        [Byte[]]$ticket_length5 = Get-ASN1LengthArrayLong ($ticket.Count + $ticket_length.Count + $ticket_length2.Count + $ticket_length3.Count + $ticket_length4.Count +
            $SPN.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + $sname_length2.Count +
            $sname_length3.Count + $sname_length4.Count + $Realm.Count + $realm_length.Count + $realm_length2.Count + 34)
        [Byte[]]$ticket_length6 = Get-ASN1LengthArrayLong ($ticket.Count + $ticket_length.Count + $ticket_length2.Count + $ticket_length3.Count + $ticket_length4.Count +
            $SPN.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + $sname_length2.Count +
            $sname_length3.Count + $sname_length4.Count + $Realm.Count + $realm_length.Count + $realm_length2.Count + $ticket_length5.Count + 35)
        [Byte[]]$ticket_length7 = Get-ASN1LengthArrayLong ($ticket.Count + $ticket_length.Count + $ticket_length2.Count + $ticket_length3.Count + $ticket_length4.Count +
            $SPN.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + $sname_length2.Count +
            $sname_length3.Count + $sname_length4.Count + $Realm.Count + $realm_length.Count + $realm_length2.Count + $ticket_length5.Count + $ticket_length6.Count + 36)
        [Byte[]]$apreq_length = Get-ASN1LengthArrayLong ($parameter_length + $ticket_length.Count + $ticket_length2.Count + $ticket_length3.Count +
            $ticket_length4.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + $sname_length2.Count +
            $sname_length3.Count + $sname_length4.Count + $realm_length.Count + $realm_length2.Count + $ticket_length5.Count + $ticket_length6.Count + $ticket_length7.Count + 73)
        [Byte[]]$apreq_length2 = Get-ASN1LengthArrayLong ($parameter_length + $ticket_length.Count + $ticket_length2.Count + $ticket_length3.Count +
            $ticket_length4.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + $sname_length2.Count +
            $sname_length3.Count + $sname_length4.Count + $realm_length.Count + $realm_length2.Count + $ticket_length5.Count + $ticket_length6.Count + $ticket_length7.Count +
            $apreq_length.Count + 74)
        [Byte[]]$length = Get-ASN1LengthArrayLong ($parameter_length + $ticket_length.Count + $ticket_length2.Count + $ticket_length3.Count +
            $ticket_length4.Count + $namestring1_length.Count + $namestring_length.Count + $namestring_length2.Count + $sname_length.Count + $sname_length2.Count +
            $sname_length3.Count + $sname_length4.Count + $realm_length.Count + $realm_length2.Count + $ticket_length5.Count + $ticket_length6.Count + $ticket_length7.Count +
            $apreq_length.Count + $apreq_length2.Count + 88)
        
        $KerberosAPREQ = New-Object System.Collections.Specialized.OrderedDictionary
        $KerberosAPREQ.Add("Length",([Byte[]](0x60) + $length))
        $KerberosAPREQ.Add("MechToken_ThisMech",[Byte[]](0x06,0x09,0x2a,0x86,0x48,0x86,0xf7,0x12,0x01,0x02,0x02))
        $KerberosAPREQ.Add("MechToken_TokenID",[Byte[]](0x01,0x00))
        $KerberosAPREQ.Add("APReq_Encoding",[Byte[]](0x6e) + $apreq_length2 + [Byte[]](0x30) + $apreq_length)
        $KerberosAPREQ.Add("PVNO_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosAPREQ.Add("PVNO",[Byte[]]0x05)
        $KerberosAPREQ.Add("MSGType_Encoding",[Byte[]](0xa1,0x03,0x02,0x01))
        $KerberosAPREQ.Add("MSGType",[Byte[]](0x0e))
        $KerberosAPREQ.Add("Padding_Encoding",[Byte[]](0xa2,0x07,0x03,0x05))
        $KerberosAPREQ.Add("Padding",[Byte[]](0x00))
        $KerberosAPREQ.Add("APOptions",[Byte[]](0x20,0x00,0x00,0x00))
        $KerberosAPREQ.Add("Ticket_Encoding",[Byte[]](0xa3) + $ticket_length7 + [Byte[]](0x61) + $ticket_length6 + [Byte[]](0x30) + $ticket_length5)
        $KerberosAPREQ.Add("Ticket_TKTVNO_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosAPREQ.Add("Ticket_TKTVNO",[Byte[]](0x05))
        $KerberosAPREQ.Add("Ticket_Realm_Encoding",[Byte[]](0xa1) + $realm_length2 + [Byte[]](0x1b) + $realm_length)
        $KerberosAPREQ.Add("Ticket_Realm",$Realm)
        $KerberosAPREQ.Add("Ticket_SName_Encoding",[Byte[]](0xa2) + $sname_length4 + [Byte[]](0x30) + $sname_length3)
        $KerberosAPREQ.Add("Ticket_SName_NameType_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosAPREQ.Add("Ticket_SName_NameType",[Byte[]](0x01))
        $KerberosAPREQ.Add("Ticket_SName_NameString_Encoding",[Byte[]](0xa1) + $sname_length2 + [Byte[]](0x30) + $sname_length)
        $KerberosAPREQ.Add("Ticket_SName_NameString0_Encoding",[Byte[]](0x1b,0x03))
        $KerberosAPREQ.Add("Ticket_SName_NameString0",[Byte[]](0x44,0x4e,0x53))
        $KerberosAPREQ.Add("Ticket_SName_NameString1_Encoding",[Byte[]](0x1b) + $namestring1_length)
        $KerberosAPREQ.Add("Ticket_SName_NameString1",$SPN)
        $KerberosAPREQ.Add("Ticket_EncPart_Encoding",[Byte[]](0xa3) + $ticket_length4 + [Byte[]](0x30) + $ticket_length3)
        $KerberosAPREQ.Add("Ticket_EncPart_EType_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosAPREQ.Add("Ticket_EncPart_EType",[Byte[]](0x12))
        $KerberosAPREQ.Add("Ticket_EncPart_KVNO_Encoding",[Byte[]](0xa1,0x03,0x02,0x01))
        $KerberosAPREQ.Add("Ticket_EncPart_KVNO",$KVNO)
        $KerberosAPREQ.Add("Ticket_EncPart_Cipher_Encoding",[Byte[]](0xa2) + $ticket_length2 + [Byte[]](0x04) + $ticket_length)
        $KerberosAPREQ.Add("Ticket_EncPart_Cipher",$ticket)
        $KerberosAPREQ.Add("Authenticator_Encoding",[Byte[]](0xa4) + $authenticator_length4 + [Byte[]](0x30) + $authenticator_length3)
        $KerberosAPREQ.Add("Authenticator_EType_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosAPREQ.Add("Authenticator_EType",[Byte[]](0x12))
        $KerberosAPREQ.Add("Authenticator_Cipher_Encoding",[Byte[]](0xa2) + $authenticator_length2 + [Byte[]](0x04) + $authenticator_length)
        $KerberosAPREQ.Add("Authenticator_Cipher",$Authenticator)

        return $KerberosAPREQ
    }

    function Unprotect-KerberosASREP
    {
        param([Byte[]]$ke_key,[Byte[]]$encrypted_data)

        $final_block_length = [Math]::Truncate($encrypted_data.Count % 16)
        [Byte[]]$final_block = $encrypted_data[($encrypted_data.Count - $final_block_length)..$encrypted_data.Count]
        [Byte[]]$penultimate_block = $encrypted_data[($encrypted_data.Count - $final_block_length - 16)..($encrypted_data.Count - $final_block_length - 1)]
        $AES = New-Object "System.Security.Cryptography.AesManaged"
        $AES.Mode = [System.Security.Cryptography.CipherMode]::CBC
        $AES.Padding = [System.Security.Cryptography.PaddingMode]::Zeros
        $AES.IV = 0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00
        $AES.KeySize = 256
        $AES.Key = $ke_key
        $AES_decryptor = $AES.CreateDecryptor()
        $penultimate_block_cleartext = $AES_decryptor.TransformFinalBlock($penultimate_block,0,$penultimate_block.Length)
        [Byte[]]$final_block_padding = $penultimate_block_cleartext[$final_block_length..$penultimate_block_cleartext.Count]
        $final_block += $final_block_padding
        [Byte[]]$cts_encrypted_data = $encrypted_data[0..($encrypted_data.Count - $final_block_length - 17)] + $final_block + $penultimate_block
        [Byte[]]$cleartext = $AES_decryptor.TransformFinalBlock($cts_encrypted_data,0,$cts_encrypted_data.Length)

        return $cleartext
    }

    function New-KerberosPACTimestamp
    {
        param([Byte[]]$ke_key)

        [Byte[]]$timestamp = Get-KerberosTimestampUTC
        [String]$confounder = [String](1..16 | ForEach-Object {"{0:X2}" -f (Get-Random -Minimum 1 -Maximum 255)})
        [Byte[]]$confounder = $confounder.Split(" ") | ForEach-Object{[Char][System.Convert]::ToInt16($_,16)}

        [Byte[]]$PAC_timestamp = $confounder +
                                    0x30,0x1a,0xa0,0x11,0x18,0x0f +
                                    $timestamp + 
                                    0xa1,0x05,0x02,0x03,0x01,0x70,0x16
        
        return $PAC_timestamp
    }

    function New-KerberosAuthenticator
    {
        param([Byte[]]$Realm,[Byte[]]$Username,[Byte[]]$SubKey,[Byte[]]$SequenceNumber)

        $parameter_length = $Realm.Count + $Username.Count + $SubKey.Count
        [Byte[]]$subkey_length = Get-ASN1LengthArray $SubKey.Count
        [Byte[]]$subkey_length2 = Get-ASN1LengthArray ($SubKey.Count + $subkey_length.Count + 1)
        [Byte[]]$subkey_length3 = Get-ASN1LengthArray ($SubKey.Count + $subkey_length.Count + $subkey_length2.Count + 7)
        [Byte[]]$subkey_length4 = Get-ASN1LengthArray ($SubKey.Count + $subkey_length.Count + $subkey_length2.Count + $subkey_length3.Count + 8)
        [Byte[]]$cname_length = Get-ASN1LengthArray $Username.Count
        [Byte[]]$cname_length2 = Get-ASN1LengthArray ($Username.Count + $cname_length.Count + 1)
        [Byte[]]$cname_length3 = Get-ASN1LengthArray ($Username.Count + $cname_length.Count + $cname_length2.Count + 2)
        [Byte[]]$cname_length4 = Get-ASN1LengthArray ($Username.Count + $cname_length.Count + $cname_length2.Count + $cname_length3.Count + 8)
        [Byte[]]$cname_length5 = Get-ASN1LengthArray ($Username.Count + $cname_length.Count + $cname_length2.Count + $cname_length3.Count + $cname_length4.Count + 9)
        [Byte[]]$crealm_length = Get-ASN1LengthArray $Realm.Count
        [Byte[]]$crealm_length2 = Get-ASN1LengthArray ($Realm.Count + $crealm_length.Count + 1)
        [Byte[]]$authenticator_length = Get-ASN1LengthArrayLong ($parameter_length + 99 + $crealm_length.Count + $crealm_length2.Count +
            $cname_length.Count + $cname_length2.Count + $cname_length3.Count + $cname_length4.Count + $cname_length5.Count + $subkey_length.Count + 
            $subkey_length2.Count + $subkey_length3.Count + $subkey_length4.Count)
        [Byte[]]$authenticator_length2 = Get-ASN1LengthArrayLong ($parameter_length + 100 + $crealm_length.Count + $crealm_length2.Count +
            $cname_length.Count + $cname_length2.Count + $cname_length3.Count + $cname_length4.Count + $cname_length5.Count + $subkey_length.Count + 
            $subkey_length2.Count + $subkey_length3.Count + $subkey_length4.Count + $authenticator_length.Count)

        $KerberosAuthenticator = New-Object System.Collections.Specialized.OrderedDictionary
        $KerberosAuthenticator.Add("Encoding",[Byte[]](0x62) + $authenticator_length2 + [Byte[]](0x30) + $authenticator_length)
        $KerberosAuthenticator.Add("AuthenticatorVNO_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosAuthenticator.Add("AuthenticatorVNO",[Byte[]](0x05))
        $KerberosAuthenticator.Add("CRealm_Encoding",[Byte[]](0xa1) + $crealm_length2 + [Byte[]](0x1b) + $crealm_length)
        $KerberosAuthenticator.Add("CRealm",$Realm)
        $KerberosAuthenticator.Add("CName_Encoding",[Byte[]](0xa2) + $cname_length5 + [Byte[]](0x30) + $cname_length4)
        $KerberosAuthenticator.Add("CName_NameType_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosAuthenticator.Add("CName_NameType",[Byte[]](0x01))
        $KerberosAuthenticator.Add("CName_CNameString_Encoding",[Byte[]](0xa1) + $cname_length3 + [Byte[]](0x30) +
            $cname_length2 + [Byte[]](0x1b) + $cname_length)
        $KerberosAuthenticator.Add("CName_CNameString",$Username)
        $KerberosAuthenticator.Add("CKSum_Encoding",[Byte[]](0xa3,0x25,0x30,0x23,0xa0,0x05,0x02,0x03))
        $KerberosAuthenticator.Add("CKSum_CKSumType",[Byte[]](0x00,0x80,0x03))
        $KerberosAuthenticator.Add("CKSum_Length_Encoding",[Byte[]](0xa1,0x1a,0x04,0x18))
        $KerberosAuthenticator.Add("CKSum_Length",[Byte[]](0x10,0x00,0x00,0x00))
        $KerberosAuthenticator.Add("CKSum_Bnd",[Byte[]](0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00))
        $KerberosAuthenticator.Add("CKSum_Flags",[Byte[]](0x36,0x01,0x00,0x00))
        $KerberosAuthenticator.Add("CKSum_CUSec_Encoding",[Byte[]](0xa4,0x05,0x02,0x03))
        $KerberosAuthenticator.Add("CKSum_CUSec",(Get-KerberosMicrosecond))
        $KerberosAuthenticator.Add("CKSum_CTime_Encoding",[Byte[]](0xa5,0x11,0x18,0x0f))
        $KerberosAuthenticator.Add("CKSum_CTime",(Get-KerberosTimestampUTC))
        $KerberosAuthenticator.Add("CKSum_Subkey_Encoding",[Byte[]](0xa6) + $subkey_length4 + [Byte[]](0x30) + $subkey_length3)
        $KerberosAuthenticator.Add("CKSum_Subkey_KeyType_Encoding",[Byte[]](0xa0,0x03,0x02,0x01))
        $KerberosAuthenticator.Add("CKSum_Subkey_KeyType",[Byte[]](0x12))
        $KerberosAuthenticator.Add("CKSum_Subkey_KeyValue_Encoding",[Byte[]](0xa1) + $subkey_length2 + [Byte[]](0x04) + $subkey_length)
        $KerberosAuthenticator.Add("CKSum_Subkey_KeyValue",$SubKey)
        $KerberosAuthenticator.Add("CKSum_SEQNumber_Encoding",[Byte[]](0xa7,0x06,0x02,0x04))
        $KerberosAuthenticator.Add("CKSum_SEQNumber",$SequenceNumber)

        return $KerberosAuthenticator
    }

    function Get-KerberosTimestampUTC
    {
        [DateTime]$timestamp = (Get-Date).ToUniversalTime()
        [String]$timestamp = ("{0:u}" -f $timestamp) -replace "-","" -replace " ","" -replace ":",""
        [Byte[]]$timestamp = [System.Text.Encoding]::UTF8.GetBytes($timestamp)

        return $timestamp
    }

    function Get-KerberosMicrosecond
    {
        [Int]$microseconds = Get-Date -Format ffffff
        [Byte[]]$microseconds = [System.Bitconverter]::GetBytes($microseconds)[0..2]

        return $microseconds
    }

    function Protect-KerberosAES256CTS
    {
        param([Byte[]]$ke_key,[Byte[]]$data)
        
        $AES = New-Object "System.Security.Cryptography.AesManaged"
        $AES.Mode = [System.Security.Cryptography.CipherMode]::CBC
        $AES.Padding = [System.Security.Cryptography.PaddingMode]::Zeros
        $IV = 0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00 
        $AES.IV = $IV
        $AES.KeySize = 256
        $AES.Key = $ke_key
        $AES_encryptor = $AES.CreateEncryptor()
        $data_encrypted = $AES_encryptor.TransformFinalBlock($data,0,$data.Length)
        $block_count = [Math]::Ceiling($data_encrypted.Count / 16)
        
        if($block_count -gt 2)
        {
            $data_encrypted = $data_encrypted[0..($data_encrypted.Count - 33)] + $data_encrypted[($data_encrypted.Count - 16)..$data_encrypted.Count] +
                $data_encrypted[($data_encrypted.Count - 32)..($data_encrypted.Count - 17)]
        }
        elseif($blocks -eq 2)
        {
            $data_encrypted = $data_encrypted[16..31] + $data_encrypted[0..15]
        }
        
        $final_block_length = [Math]::Truncate($data.Count % 16)
        
        if($final_block_length -ne 0)
        {
            $remove_count = 16 - $final_block_length
            $data_encrypted = $data_encrypted[0..($data_encrypted.Count - $remove_count - 1)]
        }
        
        return $data_encrypted
    }
    # TCPClient Kerberos end
    
    function Get-KerberosHMACSHA1
    {
        param([Byte[]]$key,[Byte[]]$data)

        $HMAC_SHA1 = New-Object System.Security.Cryptography.HMACSHA1
        $HMAC_SHA1.key = $key
        $hash = $HMAC_SHA1.ComputeHash($data)
        $hash = $hash[0..11]

        return $hash
    }
    
    function Get-ASN1LengthArray
    {
        param([Int]$Length)

        [Byte[]]$asn1 = [System.BitConverter]::GetBytes($Length)

        if($asn1[1] -eq 0)
        {
            $asn1 = $asn1[0]
        }
        else 
        {
            $asn1 = $asn1[1,0]
        }

        return $asn1
    }
    
    function Get-ASN1LengthArrayLong
    {
        param([Int]$Length)

        [Byte[]]$asn1 = [System.BitConverter]::GetBytes($Length)

        if($asn1[1] -eq 0)
        {
            $asn1 = $asn1[0]
            $asn1 = [Byte[]]0x81 + $asn1
        }
        else 
        {
            $asn1 = $asn1[1,0]
            $asn1 = [Byte[]]0x82 + $asn1
        }

        return $asn1
    }
    
    function New-RandomByteArray
    {
        param([Int]$Length,[Int]$Minimum=1,[Int]$Maximum=255)

        [String]$random = [String](1..$Length | ForEach-Object {"{0:X2}" -f (Get-Random -Minimum $Minimum -Maximum $Maximum)})
        [Byte[]]$random = $random.Split(" ") | ForEach-Object{[Char][System.Convert]::ToInt16($_,16)}

        return $random
    }
    
    function New-DNSNameArray
    {
        param([String]$Name)

        $character_array = $Name.ToCharArray()
        [Array]$index_array = 0..($character_array.Count - 1) | Where-Object {$character_array[$_] -eq '.'}

        if($index_array.Count -gt 0)
        {

            $name_start = 0

            ForEach ($index in $index_array)
            {
                $name_end = $index - $name_start
                [Byte[]]$name_array += $name_end
                [Byte[]]$name_array += [System.Text.Encoding]::UTF8.GetBytes($Name.Substring($name_start,$name_end))
                $name_start = $index + 1
            }

            [Byte[]]$name_array += ($Name.Length - $name_start)
            [Byte[]]$name_array += [System.Text.Encoding]::UTF8.GetBytes($Name.Substring($name_start))
        }
        else
        {
            [Byte[]]$name_array = $Name.Length
            [Byte[]]$name_array += [System.Text.Encoding]::UTF8.GetBytes($Name.Substring($name_start))
        }

        return $name_array
    }

    function New-PacketDNSQuery
    {
        param([String]$Name,[String]$Type)

        switch ($Type)
        {

            'A'
            {[Byte[]]$type = 0x00,0x01}

            'AAAA'
            {[Byte[]]$type = 0x00,0x1c}

            'CNAME'
            {[Byte[]]$type = 0x00,0x05}
            
            'MX'
            {[Byte[]]$type = 0x00,0x0f}

            'PTR'
            {[Byte[]]$type = 0x00,0x0c}

            'SRV'
            {[Byte[]]$type = 0x00,0x21}

            'TXT'
            {[Byte[]]$type = 0x00,0x10}

        }

        [Byte[]]$name = (New-DNSNameArray $Name) + 0x00
        [Byte[]]$length = [System.BitConverter]::GetBytes($Name.Count + 16)[1,0]
        [Byte[]]$transaction_ID = New-RandomByteArray 2
        $DNSQuery = New-Object System.Collections.Specialized.OrderedDictionary
        $DNSQuery.Add("Length",$length)
        $DNSQuery.Add("TransactionID",$transaction_ID)
        $DNSQuery.Add("Flags",[Byte[]](0x01,0x00))
        $DNSQuery.Add("Questions",[Byte[]](0x00,0x01))
        $DNSQuery.Add("AnswerRRs",[Byte[]](0x00,0x00))
        $DNSQuery.Add("AuthorityRRs",[Byte[]](0x00,0x00))
        $DNSQuery.Add("AdditionalRRs",[Byte[]](0x00,0x00))
        $DNSQuery.Add("Queries_Name",$name)
        $DNSQuery.Add("Queries_Type",$type)
        $DNSQuery.Add("Queries_Class",[Byte[]](0x00,0x01))

        return $DNSQuery
    }
    
    function New-PacketDNSQueryTKEY
    {
        param([Byte[]]$Name,[byte[]]$Type,[Byte[]]$APReq)
    
        [Byte[]]$transaction_ID = New-RandomByteArray 2

        if($APReq)
        {
            $mechtoken_length = Get-ASN1LengthArrayLong ($APReq.Count)
            $mechtoken_length2 = Get-ASN1LengthArrayLong ($APReq.Count + $mechtoken_length.Count + 1)
            $innercontexttoken_length = Get-ASN1LengthArrayLong ($APReq.Count + $mechtoken_length.Count + $mechtoken_length2.Count + 17) # 31
            $innercontexttoken_length2 = Get-ASN1LengthArrayLong ($APReq.Count + $mechtoken_length.Count + $mechtoken_length2.Count +
                $innercontexttoken_length.Count + 18)
            $spnego_length = Get-ASN1LengthArrayLong ($APReq.Count + $mechtoken_length.Count + $mechtoken_length2.Count +
                $innercontexttoken_length.Count + $innercontexttoken_length2.Count + 27)
            $grouped_length = $APReq.Count + $mechtoken_length.Count + $mechtoken_length2.Count + $innercontexttoken_length.Count +
                $innercontexttoken_length2.Count + $spnego_length.Count + 25
            $key_size = [System.BitConverter]::GetBytes($grouped_length + 3)[1,0]
            $RD_length = [System.BitConverter]::GetBytes($grouped_length + $key_size.Count + 27)[1,0]
            $inception = [int64](([datetime]::UtcNow)-(Get-Date "1/1/1970")).TotalSeconds
            $inception = [System.BitConverter]::GetBytes($inception)
            $inception = $inception[3..0]
        }

        if($APReq)
        {
            [Byte[]]$length = [System.BitConverter]::GetBytes($grouped_length + $Name.Count + 57)[1,0]
        }
        else
        {
            [Byte[]]$length = [System.BitConverter]::GetBytes($Name.Count + 16)[1,0]
        }

        $DNSQueryTKEY = New-Object System.Collections.Specialized.OrderedDictionary
        $DNSQueryTKEY.Add("Length",$length)
        $DNSQueryTKEY.Add("TransactionID",$transaction_ID)
        $DNSQueryTKEY.Add("Flags",[Byte[]](0x00,0x00))
        $DNSQueryTKEY.Add("Questions",[Byte[]](0x00,0x01))
        $DNSQueryTKEY.Add("AnswerRRs",[Byte[]](0x00,0x00))
        $DNSQueryTKEY.Add("AuthorityRRs",[Byte[]](0x00,0x00))

        if($apreq)
        {
            $DNSQueryTKEY.Add("AdditionalRRs",[Byte[]](0x00,0x01))
        }
        else
        {
            $DNSQueryTKEY.Add("AdditionalRRs",[Byte[]](0x00,0x00))
        }

        $DNSQueryTKEY.Add("Queries_Name",$Name)
        $DNSQueryTKEY.Add("Queries_Type",$Type)
        $DNSQueryTKEY.Add("Queries_Class",[Byte[]](0x00,0xff))

        if($apreq)
        {
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_Name",[Byte[]](0xc0,0x0c))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_Type",[Byte[]](0x00,0xf9))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_Class",[Byte[]](0x00,0xff))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_TTL",[Byte[]](0x00,0x00,0x00,0x00))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RDLength",$RD_length)
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_Algorithm",[Byte[]](0x08,0x67,0x73,0x73,0x2d,0x74,0x73,0x69,0x67,0x00))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_Inception",$inception)
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_Expiration",$inception)
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_Mode",[Byte[]](0x00,0x03))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_Error",[Byte[]](0x00,0x00))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_KeySize",$key_size)
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_SPNego_Encoding",[Byte[]](0x60) + $spnego_length)
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_SPNego_ThisMech",[Byte[]](0x06,0x06,0x2b,0x06,0x01,0x05,0x05,0x02))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_SPNego_InnerContextToken_Encoding",[Byte[]](0xa0) + $innercontexttoken_length2 + [Byte[]](0x30) +
                $innercontexttoken_length)
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_SPNego_InnerContextToken_MechTypes_Encoding",[Byte[]](0xa0,0x0d,0x30,0x0b))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_SPNego_InnerContextToken_MechType0",[Byte[]](0x06,0x09,0x2a,0x86,0x48,0x86,0xf7,0x12,0x01,0x02,0x02))
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_SPNego_InnerContextToken_MechToken_Encoding",[Byte[]](0xa2) + $mechtoken_length2 + [Byte[]](0x04) +
                $mechtoken_length)
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_SPNego_InnerContextToken_MechToken_Token",$APReq)
            $DNSQueryTKEY.Add("Queries_AdditionalRecords_RData_OtherSize",[Byte[]](0x00,0x00))
        }

        return $DNSQueryTKEY
    }
    
    function New-PacketDNSUpdate
    {
        param([Byte[]]$TransactionID,[String]$Zone,[String]$Name,[String]$Type,[Int]$TTL,[Int]$Preference,[Int]$Priority,[Int]$Weight,[Int]$Port,[String]$Data,[Byte[]]$TimeSigned,[Byte[]]$TKeyname,[Byte[]]$MAC)

        if($Data)
        {
            $add = $true
            [Byte[]]$class = 0x00,0x01
        }
        else
        {
            [Byte[]]$class = 0x00,0xff
            $TTL = 0
        }

        switch ($Type) 
        {

            'A'
            {
                [Byte[]]$type = 0x00,0x01
                
                if($Data -and [Bool]($Data -as [System.Net.IPAddress]))
                {
                    [Byte[]]$data = ([System.Net.IPAddress][String]([System.Net.IPAddress]$Data)).GetAddressBytes()
                }
                elseif($Data)
                {
                    [Byte[]]$data = [System.Text.Encoding]::UTF8.GetBytes($Data)
                }

            }

            'AAAA'
            {
                [Byte[]]$type = 0x00,0x1c
                
                if($Data -and [Bool]($Data -as [System.Net.IPAddress]))
                {
                    [Byte[]]$data = ([System.Net.IPAddress][String]([System.Net.IPAddress]$Data)).GetAddressBytes()
                }
                elseif($Data)
                {
                    [Byte[]]$data = [System.Text.Encoding]::UTF8.GetBytes($Data)
                }

            }

            'CNAME'
            {
                [Byte[]]$type = 0x00,0x05

                if($Data -and [Bool]($Data -as [System.Net.IPAddress]))
                {
                    [Byte[]]$data = (New-DNSNameArray $Data) + 0x00
                }
                elseif($Data)
                {
                    [Byte[]]$data = (New-DNSNameArray ($Data -replace ('.' + $Zone),'')) + 0xc0,0x0c
                }

            }

            'MX' 
            {
                $MX = $true
                [Byte[]]$type = 0x00,0x0f

                if($Data)
                {
                    $extra_length = 2
                    [Byte[]]$preference = [System.Bitconverter]::GetBytes($Preference)[1,0]
                }

                if($Data -and [Bool]($Data -as [System.Net.IPAddress]))
                {
                    [Byte[]]$data = (New-DNSNameArray $Data) + 0x00
                }
                elseif($Data)
                {
                    [Byte[]]$data = (New-DNSNameArray ($Data -replace ('.' + $Zone),'')) + 0xc0,0x0c
                }

            }

            'PTR'
            {
                [Byte[]]$type = 0x00,0x0c

                if($Data)
                {
                    [Byte[]]$data = (New-DNSNameArray $Data) + 0x00
                }

            }

            'SRV'
            {
                $SRV = $true
                [Byte[]]$type = 0x00,0x21
                
                if($Data)
                {
                    [Byte[]]$priority = [System.Bitconverter]::GetBytes($Priority)[1,0]
                    [Byte[]]$weight = [System.Bitconverter]::GetBytes($Weight)[1,0]
                    [Byte[]]$port = [System.Bitconverter]::GetBytes($Port)[1,0]
                    $extra_length = 6
                    [Byte[]]$data = (New-DNSNameArray $Data) + 0x00
                }

            }

            'TXT'
            {
                $TXT = $true
                [Byte[]]$type = 0x00,0x10
                [Byte[]]$TXT_length = [System.BitConverter]::GetBytes($Data.Length)[0]

                if($Data)
                {
                    $extra_length = 1
                    [Byte[]]$data = [System.Text.Encoding]::UTF8.GetBytes($Data)
                }

            }

        }

        if($Name -eq $Zone)
        {
            [Byte[]]$name = 0xc0,0x0c
        }
        else
        {
            [Byte[]]$name = (New-DNSNameArray ($Name -replace ('.' + $Zone),'')) + 0xc0,0x0c
        }
        
        [Byte[]]$Zone = (New-DNSNameArray $Zone) + 0x00  
        [Byte[]]$TTL = [System.Bitconverter]::GetBytes($TTL)[3..0]
        [Byte[]]$data_length = [System.BitConverter]::GetBytes($data.Length + $extra_length)[1,0]

        if($MAC)
        {
            [Byte[]]$length = [System.BitConverter]::GetBytes($Zone.Count + $name.Count + $data.Length + $TKeyname.Count + $MAC.Count + 62 + $extra_length)[1,0]
        }
        elseif(!$TKeyname)
        {
            [Byte[]]$length = [System.BitConverter]::GetBytes($Zone.Count + $name.Count + $data.Length + 26 + $extra_length)[1,0]
        }

        $DNSUpdate = New-Object System.Collections.Specialized.OrderedDictionary

        if(!$TKeyname -or $MAC)
        {
            $DNSUpdate.Add("Length",$length)
        }

        $DNSUpdate.Add("TransactionID",$TransactionID)
        $DNSUpdate.Add("Flags",[Byte[]](0x28,0x00))
        $DNSUpdate.Add("Zones",[Byte[]](0x00,0x01))
        $DNSUpdate.Add("Prerequisites",[Byte[]](0x00,0x00))
        $DNSUpdate.Add("Updates",[Byte[]](0x00,0x01))

        if($MAC)
        {
            $DNSUpdate.Add("AdditionalRRs",[Byte[]](0x00,0x01))
        }
        else
        {
            $DNSUpdate.Add("AdditiionalRRs",[Byte[]](0x00,0x00))
        }

        $DNSUpdate.Add("Zone_Name",$Zone)
        $DNSUpdate.Add("Zone_Type",[Byte[]](0x00,0x06))
        $DNSUpdate.Add("Zone_Class",[Byte[]](0x00,0x01))
        $DNSUpdate.Add("Updates_Name",$name)
        $DNSUpdate.Add("Updates_Type",$type)
        $DNSUpdate.Add("Updates_Class",$class)
        $DNSUpdate.Add("Updates_TTL",$TTL)
        $DNSUpdate.Add("Updates_DataLength",$data_length)

        if($MX)
        {
            $DNSUpdate.Add("Updates_TXTLength",$preference)
        }

        if($TXT -and $add)
        {
            $DNSUpdate.Add("Updates_TXTLength",$TXT_length)
        }

        if($SRV -and $add)
        {
            $DNSUpdate.Add("Updates_Priority",$priority)
            $DNSUpdate.Add("Updates_Weight",$weight)
            $DNSUpdate.Add("Updates_Port",$port)
        }

        if($add)
        {
            $DNSUpdate.Add("Updates_Address",$data)
        }

        if($TKeyname)
        {
            $DNSUpdate.Add("AdditionalRecords_Name",$TKeyname)

            if($MAC)
            {
                $DNSUpdate.Add("AdditionalRecords_Type",[Byte[]](0x00,0xfa))
            }

            $DNSUpdate.Add("AdditionalRecords_Class",[Byte[]](0x00,0xff))
            $DNSUpdate.Add("AdditionalRecords_TTL",[Byte[]](0x00,0x00,0x00,0x00))

            if($MAC)
            {
                $DNSUpdate.Add("AdditionalRecords_DataLength",[Byte[]](0x00,0x36))
            }

            $DNSUpdate.Add("AdditionalRecords_AlgorithmName",[Byte[]](0x08,0x67,0x73,0x73,0x2d,0x74,0x73,0x69,0x67,0x00))
            $DNSUpdate.Add("AdditionalRecords_TimeSigned",$TimeSigned)
            $DNSUpdate.Add("AdditionalRecords_Fudge",[Byte[]](0x01,0x2c))

            if($MAC)
            {
                $DNSUpdate.Add("AdditionalRecords_MACSize",[Byte[]](0x00,0x1c))
                $DNSUpdate.Add("AdditionalRecords_MAC",$MAC)
                $DNSUpdate.Add("AdditionalRecords_OriginalID",$TransactionID)
            }

            $DNSUpdate.Add("AdditionalRecords_Error",[Byte[]](0x00,0x00))
            $DNSUpdate.Add("AdditionalRecords_OtherLength",[Byte[]](0x00,0x00))
        }

        return $DNSUpdate
    }
    
    function New-PacketDNSUpdateMAC
    {
        param([Byte[]]$Flags,[Byte[]]$SequenceNumber,[Byte[]]$Checksum)

        $DNSUpdateMAC = New-Object System.Collections.Specialized.OrderedDictionary
        $DNSUpdateMAC.Add("DNSUpdateMAC_TokenID",[Byte[]](0x04,0x04))
        $DNSUpdateMAC.Add("DNSUpdateMAC_Flags",$Flags)
        $DNSUpdateMAC.Add("DNSUpdateMAC_Filler",[Byte[]](0xff,0xff,0xff,0xff,0xff))
        $DNSUpdateMAC.Add("DNSUpdateMAC_SequenceNumber",[Byte[]](0x00,0x00,0x00,0x00) + $SequenceNumber)

        if($Checksum)
        {
            $DNSUpdateMAC.Add("DNSUpdateMAC_Checksum",$Checksum)
        }

        return $DNSUpdateMAC
    }

    function Get-DNSUpdateResponseStatus
    {
        param([Byte[]]$DNSClientReceive)

        $DNS_response_flags = [System.BitConverter]::ToString($DNSClientReceive[4..5])
        $DNS_response_flags = $DNS_response_flags -replace "-",""

        switch ($DNS_response_flags)
        {
            'A800' {$DNS_update_response_status = "[+] DNS update successful"}
            'A801' {$DNS_update_response_status = ("[-] format error 0x" + $DNS_response_flags)}
            'A802' {$DNS_update_response_status = ("[-] failed to complete 0x" + $DNS_response_flags)}
            'A804' {$DNS_update_response_status = ("[-] not implemented 0x" + $DNS_response_flags)}
            'A805' {$DNS_update_response_status = ("[-] update refused 0x" + $DNS_response_flags)}
            Default {$DNS_update_response_status = ("[-] DNS update was not successful 0x" + $DNS_response_flags)}
        }

        return $DNS_update_response_status
    }

    if($RecordCheck)
    {

        if($DNSType -ne 'MX' -and $DNSName -notlike '*.*')
        {
            $query_name = $DNSName + "." + $Zone
        }
        else
        {
            $query_name = $DNSName    
        }

        $DNS_client = New-Object System.Net.Sockets.TCPClient
        $DNS_client.Client.ReceiveTimeout = 3000

        try
        {
            $DNS_client.Connect($DomainController,"53")
            $DNS_client_stream = $DNS_client.GetStream()
            $DNS_client_receive = New-Object System.Byte[] 2048
            $packet_DNSQuery = New-PacketDNSQuery $query_name $DNSType
            [Byte[]]$DNS_client_send = ConvertFrom-PacketOrderedDictionary $packet_DNSQuery
            $DNS_client_stream.Write($DNS_client_send,0,$DNS_client_send.Length) > $null
            $DNS_client_stream.Flush()   
            $DNS_client_stream.Read($DNS_client_receive,0,$DNS_client_receive.Length) > $null
            $DNS_client.Close()
            $DNS_client_stream.Close()

            if($DNS_client_receive[9] -ne 0)
            {
                $DNS_record_exists = $true
                Write-Output "[-] $DNSName of record type $DNSType already exists"
            }

        }
        catch
        {
            Write-Output "[-] $DomainController did not respond on TCP port 53"
        }

    }

    if(!$RecordCheck -or ($RecordCheck -and !$DNS_record_exists))
    {
        $DNS_client = New-Object System.Net.Sockets.TCPClient
        $DNS_client.Client.ReceiveTimeout = 3000

        if($Security -ne 'Secure')
        {

            try
            {
                $DNS_client.Connect($DomainController,"53")
            }
            catch
            {
                Write-Output "$DomainController did not respond on TCP port 53"
            }

            if($DNS_client.Connected)
            {
                $DNS_client_stream = $DNS_client.GetStream()
                $DNS_client_receive = New-Object System.Byte[] 2048
                [Byte[]]$transaction_id = New-RandomByteArray 2
                $packet_DNSUpdate = New-PacketDNSUpdate $transaction_ID $Zone $DNSName $DNSType $DNSTTL $DNSPreference $DNSPriority $DNSWeight $DNSPort $DNSData
                [Byte[]]$DNSUpdate = ConvertFrom-PacketOrderedDictionary $packet_DNSUpdate
                $DNS_client_send = $DNSUpdate
                $DNS_client_stream.Write($DNS_client_send,0,$DNS_client_send.Length) > $null
                $DNS_client_stream.Flush()   
                $DNS_client_stream.Read($DNS_client_receive,0,$DNS_client_receive.Length) > $null
                $DNS_update_response_status = Get-DNSUpdateResponseStatus $DNS_client_receive
                Write-Output $DNS_update_response_status
                $DNS_client.Close()
                $DNS_client_stream.Close()
            }

        }

        if($Security -eq 'Secure' -or ($Security -eq 'Auto' -and $DNS_update_response_status -like '*0xA805'))
        {
            $tkey = "6" + ((0..9) | Get-Random -Count 2) + "-ms-7.1-" + ((0..9) | Get-Random -Count 4) + "." + ((0..9) | Get-Random -Count 8) +
                "-" + ((0..9) | Get-Random -Count 4) + "-11e7-" + ((0..9) | Get-Random -Count 4) + "-000c296694e0"
            $tkey = $tkey -replace " ",""
            Write-Verbose "[+] TKEY name $tkey"
            [Byte[]]$tkey_name = [System.Text.Encoding]::UTF8.GetBytes($tkey)
            $tkey_name = [Byte[]]0x08 + $tkey_name + 0x00
            $tkey_name[9] = 0x06
            $tkey_name[16] = 0x24

            if($kerberos_tcpclient)
            {
                $kerberos_client = New-Object System.Net.Sockets.TCPClient
                $kerberos_client.Client.ReceiveTimeout = 3000
                $domain_controller = [System.Text.Encoding]::UTF8.GetBytes($DomainController)
                $kerberos_username = [System.Text.Encoding]::UTF8.GetBytes($Username)
                $kerberos_realm = [System.Text.Encoding]::UTF8.GetBytes($Realm)

                try
                {
                    $kerberos_client.Connect($DomainController,"88")
                }
                catch
                {
                    Write-Output "$DomainController did not respond on TCP port 88"
                }

            }

            if(!$kerberos_tcpclient -or $kerberos_client.Connected)
            {

                if($kerberos_tcpclient)
                {

                    if($Hash)
                    {
                        $base_key = (&{for ($i = 0;$i -lt $hash.Length;$i += 2){$hash.SubString($i,2)}}) -join "-"
                        $base_key = $base_key.Split("-") | ForEach-Object{[Char][System.Convert]::ToInt16($_,16)}
                    }
                    else
                    {
                        $base_key = Get-KerberosAES256BaseKey $salt $password
                    }

                    $ke_key = Get-KerberosAES256UsageKey encrypt 1 $base_key
                    $ki_key = Get-KerberosAES256UsageKey integrity 1 $base_key
                    $nonce = New-RandomByteArray 4        
                    $kerberos_client_stream = $kerberos_client.GetStream()
                    $kerberos_client_receive = New-Object System.Byte[] 2048
                    $packet_AS_REQ = New-PacketKerberosASREQ $kerberos_username $kerberos_realm $domain_controller $nonce
                    $AS_REQ = ConvertFrom-PacketOrderedDictionary $packet_AS_REQ
                    $kerberos_client_send = $AS_REQ
                    $kerberos_client_stream.Write($kerberos_client_send,0,$kerberos_client_send.Length) > $null
                    $kerberos_client_stream.Flush()
                    $kerberos_client_stream.Read($kerberos_client_receive,0,$kerberos_client_receive.Length) > $null
                    [Byte[]]$PAC_Timestamp = New-KerberosPACTimestamp $ke_key
                    [Byte[]]$PAC_ENC_Timestamp = Protect-KerberosAES256CTS $ke_key $PAC_Timestamp
                    [Byte[]]$PAC_Timestamp_Signature = Get-KerberosHMACSHA1 $ki_key $PAC_Timestamp
                    $packet_AS_REQ = New-PacketKerberosASREQ $kerberos_username $kerberos_realm $domain_controller $nonce $PAC_ENC_Timestamp $PAC_Timestamp_Signature
                    $AS_REQ = ConvertFrom-PacketOrderedDictionary $packet_AS_REQ
                    $kerberos_client_send = $AS_REQ
                    $kerberos_client_stream.Write($kerberos_client_send,0,$kerberos_client_send.Length) > $null
                    $kerberos_client_stream.Flush()   
                    $kerberos_client_stream.Read($kerberos_client_receive,0,$kerberos_client_receive.Length) > $null
                    $asrep_payload = [System.BitConverter]::ToString($kerberos_client_receive)
                    $asrep_payload = $asrep_payload -replace "-",""
                    $kerberos_client.Close()
                    $kerberos_client_stream.Close()
                }
                else
                {
                    
                    try
                    {

                        $Null = [System.Reflection.Assembly]::LoadWithPartialName("System.IdentityModel")

                        if($username -or $Credential)
                        {

                            if(!$Credential)
                            {
                                $Credential = New-Object System.Management.Automation.PSCredential ($username,$Password)
                            }

                            $network_creds = $Credential.GetNetworkCredential()
                            $network_creds.Domain = $domain
                            $token = New-Object  System.IdentityModel.Selectors.KerberosSecurityTokenProvider ("DNS/$DomainController",[System.Security.Principal.TokenImpersonationLevel]::Impersonation,$network_creds)
                            $ticket = $token.GetToken([System.TimeSpan]::FromMinutes(1))
                        }
                        else
                        {
                            $ticket = New-Object  System.IdentityModel.Tokens.KerberosRequestorSecurityToken ("DNS/$DomainController")
                        }

                        $asrep_key = $ticket.SecurityKey.GetSymmetricKey()
                        $kerberos_client_receive = $Ticket.GetRequest()
                        $asrep_payload = [System.BitConverter]::ToString($kerberos_client_receive)
                        $asrep_payload = $asrep_payload -replace "-",""
                    }
                    catch
                    {
                        $auth_success = $false
                    }

                }

                if($asrep_key -or ($asrep_payload.Length -gt 0 -and $asrep_payload -like '*A003020105A10302010B*'))
                {
                    Write-Verbose "[+] Kerberos preauthentication successful"
                    $auth_success = $true  
                }
                elseif($asrep_payload.Length -gt 0 -and $asrep_payload -like '*A003020105A10302011E*')
                {
                    Write-Output ("[-] Kerberos preauthentication error 0x" + $asrep_payload.Substring(96,2))
                    $auth_success = $false
                }
                else
                {
                    Write-Output "[-] Kerberos authentication failure"
                    $auth_success = $false
                }

                if($auth_success)
                {
                    $ticket_index = $asrep_payload.IndexOf("A003020112A1030201")
                    $ticket_kvno = $kerberos_client_receive[($ticket_index / 2 + 9)]
                    
                    if($asrep_payload.Substring($ticket_index + 22,2) -eq '82')
                    {
                        $ticket_length = ([System.BitConverter]::ToUInt16($kerberos_client_receive[($ticket_index / 2 + 13)..($ticket_index / 2 + 12)],0)) - 4
                    }
                    else
                    {
                        $ticket_length = $kerberos_client_receive[($ticket_index / 2 + 12)] - 3
                    }

                    $ticket = $Kerberos_client_receive[($ticket_index / 2 + 18)..($ticket_index/2 + 17 + $ticket_length)]

                    if($kerberos_tcpclient)
                    {
                        $cipher_index = $asrep_payload.Substring($ticket_index + 1).IndexOf("A003020112A1030201") + $ticket_index + 1

                        if($asrep_payload.Substring($cipher_index + 22,2) -eq '82')
                        {
                            $cipher_length = ([System.BitConverter]::ToUInt16($kerberos_client_receive[($cipher_index / 2 + 13)..($cipher_index / 2 + 12)],0)) - 4
                        }
                        else
                        {
                            $cipher_length = $kerberos_client_receive[($cipher_length / 2 + 12)] - 3
                        }

                        $cipher = $kerberos_client_receive[($cipher_index / 2 + 18)..($cipher_index / 2 + 17 + $cipher_length)]
                        $ke_key = Get-KerberosAES256UsageKey encrypt 3 $base_key
                        $asrep_cleartext = Unprotect-KerberosASREP $ke_key $cipher[0..($cipher.Count - 13)]
                        $kerberos_session_key = $asrep_cleartext[37..68]
                        $ke_key = Get-KerberosAES256UsageKey encrypt 11 $kerberos_session_key
                        $ki_key = Get-KerberosAES256UsageKey integrity 11 $kerberos_session_key
                        [Byte[]]$subkey = New-RandomByteArray 32
                        [Byte[]]$sequence_number = New-RandomByteArray 4
                        $packet_authenticator = New-KerberosAuthenticator $kerberos_realm $kerberos_username $subkey $sequence_number
                        [Byte[]]$authenticator = ConvertFrom-PacketOrderedDictionary $packet_authenticator
                        $authenticator = (New-RandomByteArray 16) + $authenticator
                        $authenticator_encrypted = Protect-KerberosAES256CTS $ke_key $authenticator
                        $authenticator_signature = Get-KerberosHMACSHA1 $ki_key $authenticator
                        $packet_apreq = New-PacketKerberosAPREQ $kerberos_realm $domain_controller $ticket_kvno $ticket $authenticator_encrypted $authenticator_signature
                        [Byte[]]$apreq = ConvertFrom-PacketOrderedDictionary $packet_apreq
                        [Byte[]]$mac_flags = 0x04
                    }
                    else
                    {
                        [Byte[]]$apreq = $kerberos_client_receive
                        [Byte[]]$mac_flags = 0x00
                    }
                        
                    $packet_DNSQuery = New-PacketDNSQueryTKEY $tkey_name 0x00,0xf9 $apreq
                    $DNSQueryTKEY = ConvertFrom-PacketOrderedDictionary $packet_DNSQuery
                    $DNS_client = New-Object System.Net.Sockets.TCPClient
                    $DNS_client.Client.ReceiveTimeout = 3000

                    try
                    {
                        $DNS_client.Connect($DomainController,"53")
                    }
                    catch
                    {
                        Write-Output "$DomainController did not respond on TCP port 53"
                    }

                    if($DNS_client.Connected)
                    {
                        $DNS_client_stream = $DNS_client.GetStream()
                        $DNS_client_receive = New-Object System.Byte[] 2048
                        $DNS_client_send = $DNSQueryTKEY
                        $DNS_client_stream.Write($DNS_client_send,0,$DNS_client_send.Length) > $null
                        $DNS_client_stream.Flush()   
                        $DNS_client_stream.Read($DNS_client_receive,0,$DNS_client_receive.Length) > $null
                        $tkey_payload = [System.BitConverter]::ToString($DNS_client_receive)
                        $tkey_payload = $tkey_payload -replace "-",""
                        
                        if($tkey_payload.Substring(8,4) -eq '8000')
                        {
                            Write-Verbose "[+] Kerberos TKEY query successful"
                            $TKEY_success = $true         
                        }
                        else
                        {
                            Write-Output ("[-] Kerberos TKEY query error 0x" + $tkey_payload.Substring(8,4))
                            $TKEY_success = $false
                        }

                        if($TKEY_success)
                        {

                            if($kerberos_tcpclient)
                            {
                                $cipher_index = $tkey_payload.IndexOf("A003020112A2")
                                $cipher_length = $DNS_client_receive[($cipher_index / 2 + 8)]
                                $cipher = $DNS_client_receive[($cipher_index / 2 + 9)..($cipher_index / 2 + 8 + $cipher_length)]
                                $ke_key = Get-KerberosAES256UsageKey encrypt 12 $kerberos_session_key
                                $tkey_cleartext = Unprotect-KerberosASREP $ke_key $cipher[0..($cipher.Count - 13)]
                                $acceptor_subkey = $tkey_cleartext[59..90]
                            }
                            else
                            {
                                $sequence_index = $tkey_payload.IndexOf("FFFFFFFFFF00000000")
                                $sequence_number = $DNS_client_receive[($sequence_index / 2 + 9)..($sequence_index / 2 + 12)]
                                $acceptor_subkey = $asrep_key
                            }

                            $kc_key = Get-KerberosAES256UsageKey checksum 25 $acceptor_subkey
                            $time_signed = [Int](([DateTime]::UtcNow)-(Get-Date "1/1/1970")).TotalSeconds
                            $time_signed = [System.BitConverter]::GetBytes($time_signed)
                            $time_signed = 0x00,0x00 + $time_signed[3..0]
                            [Byte[]]$transaction_id = New-RandomByteArray 2
                            $packet_DNSUpdate = New-PacketDNSUpdate $transaction_ID $Zone $DNSName $DNSType $DNSTTL $DNSPreference $DNSPriority $DNSWeight $DNSPort $DNSData $time_signed $tkey_name
                            [Byte[]]$DNSUpdateTSIG = ConvertFrom-PacketOrderedDictionary $packet_DNSUpdate
                            $packet_DNSUpdateMAC = New-PacketDNSUpdateMAC $mac_flags $sequence_number
                            [Byte[]]$DNSUpdateMAC = ConvertFrom-PacketOrderedDictionary $packet_DNSUpdateMAC
                            $DNSUpdateTSIG += $DNSUpdateMAC
                            $checksum = Get-KerberosHMACSHA1 $kc_key $DNSUpdateTSIG
                            $packet_DNSUpdateMAC = New-PacketDNSUpdateMAC $mac_flags $sequence_number $checksum
                            [Byte[]]$DNSUpdateMAC = ConvertFrom-PacketOrderedDictionary $packet_DNSUpdateMAC
                            $packet_DNSUpdate = New-PacketDNSUpdate $transaction_ID $Zone $DNSName $DNSType $DNSTTL $DNSPreference $DNSPriority $DNSWeight $DNSPort $DNSData $time_signed $tkey_name $DNSUpdateMAC
                            [Byte[]]$DNSUpdateTSIG = ConvertFrom-PacketOrderedDictionary $packet_DNSUpdate
                            $DNS_client_send = $DNSUpdateTSIG
                            $DNS_client_stream.Write($DNS_client_send,0,$DNS_client_send.Length) > $null
                            $DNS_client_stream.Flush()   
                            $DNS_client_stream.Read($DNS_client_receive,0,$DNS_client_receive.Length) > $null
                            $DNS_update_response_status = Get-DNSUpdateResponseStatus $DNS_client_receive
                            Write-Output $DNS_update_response_status
                            $DNS_client.Close()
                            $DNS_client_stream.Close()
                        }

                    }

                }

            }

        }

    }

}
