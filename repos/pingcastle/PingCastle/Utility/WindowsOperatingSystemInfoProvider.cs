namespace PingCastle.Utility;

using System;
using System.Diagnostics;
using System.Management;
using PingCastleCommon.Utility;

/// <summary>
/// Windows-specific implementation for retrieving operating system information via WMI.
/// This class uses Windows Management Instrumentation to query OS details from remote computers.
/// </summary>
public class WindowsOperatingSystemInfoProvider : IOperatingSystemInfoProvider
{
    /// <summary>
    /// Retrieves the operating system name from a remote computer using WMI.
    /// </summary>
    /// <param name="computerName">The computer name to query</param>
    /// <returns>The operating system name (e.g., "Microsoft Windows Server 2025 Standard"), or empty string if retrieval fails</returns>
    public string GetOperatingSystemName(string computerName)
    {
        try
        {
            var connectionOptions = new ConnectionOptions();
            connectionOptions.Impersonation = ImpersonationLevel.Impersonate;

            var managementScope = new ManagementScope($"\\\\{computerName}\\root\\cimv2", connectionOptions);
            managementScope.Connect();

            var query = new ObjectQuery("SELECT Caption FROM Win32_OperatingSystem");
            var searcher = new ManagementObjectSearcher(managementScope, query);
            var results = searcher.Get();

            foreach (ManagementObject mo in results)
            {
                string caption = mo["Caption"]?.ToString();
                if (!string.IsNullOrEmpty(caption))
                {
                    return caption;
                }
            }
        }
        catch (Exception ex)
        {
            Trace.WriteLine($"WMI retrieval failed for {computerName}: {ex.Message}");
        }

        return string.Empty;
    }
}