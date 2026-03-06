using System.Collections.Generic;
using System.Diagnostics;
using PingCastle.UserInterface;
using PingCastle.misc;

namespace PingCastle.Healthcheck
{
    using PingCastleCommon.Utility;

    public class HotFixCollector : IHotFixCollector
    {
        private readonly IHotfixService _hotfixService;

        /// <summary>
        /// Initializes a new instance of the HotFixCollector class.
        /// </summary>
        /// <param name="hotfixService">The hotfix service to use for retrieving hotfixes</param>
        public HotFixCollector(IHotfixService hotfixService)
        {
            _hotfixService = hotfixService ?? throw new System.ArgumentNullException(nameof(hotfixService));
        }

        /// <summary>
        /// Attempts to retrieve installed hotfixes from a remote computer using WMI.
        /// If WMI fails, logs a warning and returns an empty HashSet.
        /// Hotfix retrieval requires privileged mode access to remote systems.
        /// </summary>
        /// <param name="hostName">Target computer hostname or IP address</param>
        /// <param name="hotfixes">Output set of discovered KB numbers (e.g., "KB4012598")</param>
        /// <param name="isPrivilegedMode">Whether privileged mode is active. Hotfix collection requires high privileges.</param>
        /// <returns>True if hotfixes were successfully retrieved, false if WMI failed or privilege check failed</returns>
        public bool TryGetInstalledHotfixes(string hostName, out HashSet<string> hotfixes, bool isPrivilegedMode = true)
        {
            return TryGetInstalledHotfixes(hostName, out hotfixes, UserInterfaceFactory.GetUserInterface(), isPrivilegedMode);
        }


        /// <summary>
        /// Attempts to retrieve installed hotfixes from a remote computer using WMI.
        /// If WMI fails, logs a warning and returns an empty HashSet.
        /// Hotfix retrieval requires privileged mode access to remote systems.
        /// </summary>
        /// <param name="hostName">Target computer hostname or IP address</param>
        /// <param name="hotfixes">Output set of discovered KB numbers (e.g., "KB4012598")</param>
        /// <param name="ui">User interface for displaying messages</param>
        /// <param name="isPrivilegedMode">Whether privileged mode is active. Hotfix collection requires high privileges.</param>
        /// <returns>True if hotfixes were successfully retrieved, false if WMI failed or privilege check failed</returns>
        public bool TryGetInstalledHotfixes(string hostName, out HashSet<string> hotfixes, IUserInterface ui, bool isPrivilegedMode = true)
        {
            hotfixes = new HashSet<string>();

            if (!isPrivilegedMode)
            {
                Trace.WriteLine($"Hotfix collection skipped for {hostName} - not in privileged mode");
                return false;
            }

            // Validate hostname to prevent command/script injection
            if (!NetworkHelper.IsValidHostName(hostName))
            {
                var errorMsg = $"Invalid hostname '{hostName}' - contains potentially malicious characters or exceeds length limits";
                Trace.WriteLine(errorMsg);
                ui.DisplayMessage(errorMsg);
                return false;
            }

            // Try WMI-based hotfix detection
            if (_hotfixService.TryGetInstalledHotfixes(hostName, out hotfixes, ui))
            {
                if (hotfixes.Count > 0)
                {
                    Trace.WriteLine($"Successfully retrieved {hotfixes.Count} hotfixes from {hostName} using WMI");
                    return true;
                }
            }

            // WMI failed - log warning and return empty HashSet
            var warningMsg = $"WMI hotfix detection failed for {hostName}. Hotfix information will be unavailable for security analysis.";
            Trace.WriteLine(warningMsg);
            ui.DisplayMessage(warningMsg);

            hotfixes = new HashSet<string>();
            return false;
        }
    }
}