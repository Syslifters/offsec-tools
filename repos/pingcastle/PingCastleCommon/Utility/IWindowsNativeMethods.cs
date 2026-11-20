namespace PingCastle;

using System;
using System.Security.Principal;

/// <summary>
/// Public interface for Windows native method operations.
/// </summary>
public interface IWindowsNativeMethods
{
    /// <summary>
    /// Converts a SID string to a domain\account name using the Windows API.
    /// </summary>
    string ConvertSIDToName(string sidstring, string server, out string referencedDomain);

    /// <summary>
    /// Gets a SecurityIdentifier for a domain name using the Windows API.
    /// </summary>
    SecurityIdentifier GetSidFromDomainName(string server, string domainToResolve);

    /// <summary>
    /// Gets the startup time of a remote computer.
    /// </summary>
    DateTime GetStartupTime(string server);

    /// <summary>
    /// Gets the version information of a remote computer.
    /// </summary>
    string GetComputerVersion(string server);

    /// <summary>
    /// Splits a command line string into individual arguments, handling quoted strings.
    /// </summary>
    string[] SplitArguments(string commandLine);
}