using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Management;
using System.Text.RegularExpressions;
using PingCastle.UserInterface;

namespace PingCastle.misc
{
    internal class WmiHotfixHelper : IHotfixService
    {
        private static readonly Regex KbRegex = new Regex(@"KB(\d+)", RegexOptions.Compiled | RegexOptions.IgnoreCase);

        /// <summary>
        /// Initializes a new instance of the WmiHotfixHelper class.
        /// </summary>
        public WmiHotfixHelper()
        {
        }

        /// <summary>
        /// Attempts to retrieve installed hotfixes from a remote computer using WMI.
        /// </summary>
        /// <param name="hostName">Target computer hostname or IP address</param>
        /// <param name="hotfixes">Output set of discovered KB numbers (e.g., "KB4012598")</param>
        /// <returns>True if hotfixes were successfully retrieved, false otherwise</returns>
        public bool TryGetInstalledHotfixes(string hostName, out HashSet<string> hotfixes)
        {
            return TryGetInstalledHotfixes(hostName, out hotfixes, UserInterfaceFactory.GetUserInterface());
        }

        /// <summary>
        /// Attempts to retrieve installed hotfixes from a remote computer using WMI.
        /// </summary>
        /// <param name="hostName">Target computer hostname or IP address</param>
        /// <param name="hotfixes">Output set of discovered KB numbers (e.g., "KB4012598")</param>
        /// <param name="ui">User interface for displaying messages</param>
        /// <returns>True if hotfixes were successfully retrieved, false otherwise</returns>
        public bool TryGetInstalledHotfixes(string hostName, out HashSet<string> hotfixes, IUserInterface ui)
        {
            hotfixes = new HashSet<string>();

            try
            {
                // Primary method: Win32_QuickFixEngineering (fastest and most direct)
                if (TryGetHotfixesFromQuickFixEngineering(hostName, hotfixes, ui))
                {
                    if (hotfixes.Count > 0)
                    {
                        Trace.WriteLine($"Retrieved {hotfixes.Count} hotfixes from {hostName} using Win32_QuickFixEngineering");
                        return true;
                    }
                }

                // Fallback method: Win32_Product (more comprehensive but slower)
                if (TryGetHotfixesFromWin32Product(hostName, hotfixes, ui))
                {
                    if (hotfixes.Count > 0)
                    {
                        Trace.WriteLine($"Retrieved {hotfixes.Count} hotfixes from {hostName} using Win32_Product fallback");
                        return true;
                    }
                }

                // If no hotfixes found using WMI methods
                Trace.WriteLine($"No hotfixes found on {hostName} using WMI methods");
                return false;
            }
            catch (Exception ex)
            {
                var msg = $"Could not retrieve hotfixes from {hostName} using WMI methods: {ex.Message}";
                Trace.WriteLine(msg);
                ui.DisplayMessage(msg);
                return false;
            }
        }

        /// <summary>
        /// Primary method: Uses Win32_QuickFixEngineering WMI class to retrieve hotfixes.
        /// This is the fastest and most direct approach, well-suited for most Windows systems.
        /// </summary>
        private bool TryGetHotfixesFromQuickFixEngineering(string hostName, HashSet<string> hotfixes, IUserInterface ui)
        {
            try
            {
                var connectionOptions = CreateConnectionOptions();
                var scope = new ManagementScope($"\\\\{hostName}\\root\\cimv2", connectionOptions);
                scope.Connect();

                var query = new ObjectQuery("SELECT * FROM Win32_QuickFixEngineering");
                using (var searcher = new ManagementObjectSearcher(scope, query))
                using (var collection = searcher.Get())
                {
                    foreach (ManagementObject obj in collection)
                    {
                        try
                        {
                            // HotFixID property contains the KB number (e.g., "KB4012598")
                            var hotfixId = obj["HotFixID"]?.ToString();
                            if (!string.IsNullOrEmpty(hotfixId))
                            {
                                ExtractKbFromString(hotfixId, hotfixes);
                            }

                            // Also check Description field as some updates might have KB numbers there
                            var description = obj["Description"]?.ToString();
                            if (!string.IsNullOrEmpty(description))
                            {
                                ExtractKbFromString(description, hotfixes);
                            }
                        }
                        catch (Exception ex)
                        {
                            Trace.WriteLine($"Error processing QuickFix entry on {hostName}: {ex.Message}");
                        }
                        finally
                        {
                            obj?.Dispose();
                        }
                    }
                }

                return true;
            }
            catch (UnauthorizedAccessException ex)
            {
                var msg = $"Access denied when querying Win32_QuickFixEngineering on {hostName}: {ex.Message}";
                Trace.WriteLine(msg);
                ui.DisplayMessage(msg);
                return false;
            }
            catch (Exception ex)
            {
                Trace.WriteLine($"Error querying Win32_QuickFixEngineering on {hostName}: {ex.Message}");
                return false;
            }
        }

        /// <summary>
        /// Fallback method: Uses Win32_Product WMI class to find installed products with KB numbers.
        /// This method is more comprehensive but slower and can trigger consistency checks.
        /// Only used when Win32_QuickFixEngineering doesn't return sufficient results.
        /// </summary>
        private bool TryGetHotfixesFromWin32Product(string hostName, HashSet<string> hotfixes, IUserInterface ui)
        {
            try
            {
                var connectionOptions = CreateConnectionOptions();
                var scope = new ManagementScope($"\\\\{hostName}\\root\\cimv2", connectionOptions);
                scope.Connect();

                // Query for products that contain KB in their name
                var query = new ObjectQuery("SELECT Name FROM Win32_Product WHERE Name LIKE '%KB%'");
                using (var searcher = new ManagementObjectSearcher(scope, query))
                using (var collection = searcher.Get())
                {
                    foreach (ManagementObject obj in collection)
                    {
                        try
                        {
                            var productName = obj["Name"]?.ToString();
                            if (!string.IsNullOrEmpty(productName))
                            {
                                ExtractKbFromString(productName, hotfixes);
                            }
                        }
                        catch (Exception ex)
                        {
                            Trace.WriteLine($"Error processing Win32_Product entry on {hostName}: {ex.Message}");
                        }
                        finally
                        {
                            obj?.Dispose();
                        }
                    }
                }

                return true;
            }
            catch (UnauthorizedAccessException ex)
            {
                var msg = $"Access denied when querying Win32_Product on {hostName}: {ex.Message}";
                Trace.WriteLine(msg);
                ui.DisplayMessage(msg);
                return false;
            }
            catch (Exception ex)
            {
                Trace.WriteLine($"Error querying Win32_Product on {hostName}: {ex.Message}");
                return false;
            }
        }

        /// <summary>
        /// Creates connection options for WMI queries with appropriate timeout and authentication settings.
        /// </summary>
        private ConnectionOptions CreateConnectionOptions()
        {
            var options = new ConnectionOptions
            {
                // Use current user's credentials (pass-through authentication)
                EnablePrivileges = true,

                // Set reasonable timeout (90 seconds for network operations)
                Timeout = TimeSpan.FromSeconds(90),

                // Use default authentication (typically NTLM/Kerberos)
                Authentication = AuthenticationLevel.PacketPrivacy,

                // Enable impersonation for proper access
                Impersonation = ImpersonationLevel.Impersonate
            };

            return options;
        }

        /// <summary>
        /// Extracts KB numbers from a given string using regex matching.
        /// Supports formats like "KB4012598", "kb1234567", etc.
        /// </summary>
        private void ExtractKbFromString(string input, HashSet<string> hotfixes)
        {
            if (string.IsNullOrEmpty(input))
                return;

            var matches = KbRegex.Matches(input);
            foreach (Match match in matches)
            {
                if (match.Success && match.Groups.Count > 1)
                {
                    var kbNumber = "KB" + match.Groups[1].Value;
                    hotfixes.Add(kbNumber);
                }
            }
        }

        /// <summary>
        /// Alternative method that could be used for additional hotfix detection via PowerShell cmdlets.
        /// This method uses Get-HotFix cmdlet through WMI if PowerShell remoting is available.
        /// Currently not implemented but could be added as another fallback mechanism.
        /// </summary>
        /// <remarks>
        /// This would require PowerShell remoting to be enabled on target systems:
        /// Invoke-Command -ComputerName $hostName -ScriptBlock { Get-HotFix | Select-Object HotFixID }
        /// </remarks>
        private bool TryGetHotfixesFromPowerShell(string hostName, HashSet<string> hotfixes, IUserInterface ui)
        {
            // Future implementation: Use PowerShell remoting if WMI methods fail
            // This would require additional dependencies and PowerShell remoting configuration
            return false;
        }
    }
}