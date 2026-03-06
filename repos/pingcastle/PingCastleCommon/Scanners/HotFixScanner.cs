namespace PingCastle.Scanners
{
    using PingCastle.ADWS;
    using PingCastle.Healthcheck;
    using PingCastle.UserInterface;
    using PingCastleCommon.Utility;
    using System;
    using System.Collections.Generic;
    using System.Diagnostics;

    /// <summary>
    /// Base class for scanners that need to retrieve and analyze installed hotfixes.
    /// Provides common functionality for hotfix retrieval.
    /// </summary>
    abstract public class HotFixScanner : ScannerBase
    {
        protected readonly IHotFixCollector HotfixCollector;
        protected readonly IOperatingSystemInfoProvider OsInfoProvider;

        protected HotFixScanner(IHotFixCollector hotfixCollector, IOperatingSystemInfoProvider osInfoProvider, IIdentityProvider identityProvider)
            : base(identityProvider)
        {
            HotfixCollector = hotfixCollector ?? throw new ArgumentNullException(nameof(hotfixCollector));
            OsInfoProvider = osInfoProvider ?? throw new ArgumentNullException(nameof(osInfoProvider));
        }

        /// <summary>
        /// Retrieves installed hotfixes for a domain controller using WMI.
        /// </summary>
        /// <param name="computerName">The name of the computer to scan.</param>
        /// <returns>A HashSet of installed KB numbers, or empty set if retrieval fails</returns>
        protected HashSet<string> RetrieveInstalledHotfixes(string computerName)
        {
            try
            {
                var hotfixCollector = HotfixCollector;
                IUserInterface ui = UserInterfaceFactory.GetUserInterface();

                if (hotfixCollector.TryGetInstalledHotfixes(computerName, out HashSet<string> hotfixes, ui))
                {
                    Trace.WriteLine($"Retrieved {hotfixes.Count} hotfixes for {computerName}");
                    return hotfixes;
                }

                Trace.WriteLine($"Unable to retrieve hotfixes for {computerName}");
                return new HashSet<string>();
            }
            catch (Exception ex)
            {
                Trace.WriteLine($"Error retrieving hotfixes for {computerName}: {ex.Message}");
                return new HashSet<string>();
            }
        }

        /// <summary>
        /// Retrieves the operating system information for a domain controller.
        /// </summary>
        /// <param name="computerName">The computer/domain controller name</param>
        /// <returns>The operating system name (e.g. "Windows Server 2025"), or "Unknown" if retrieval fails</returns>
        protected string RetrieveOperatingSystem(string computerName)
        {
            try
            {
                string osName = OsInfoProvider.GetOperatingSystemName(computerName);
                if (!string.IsNullOrEmpty(osName))
                {
                    Trace.WriteLine($"Retrieved OS name for {computerName}: {osName}");
                    return osName;
                }

                Trace.WriteLine($"Unable to retrieve OS information for {computerName}");
                return "Unknown";
            }
            catch (Exception ex)
            {
                Trace.WriteLine($"Error retrieving OS for {computerName}: {ex.Message}");
                return "Unknown";
            }
        }
    }
}
