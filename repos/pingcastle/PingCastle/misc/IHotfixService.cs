using System.Collections.Generic;
using PingCastle.UserInterface;

namespace PingCastle.misc
{
    /// <summary>
    /// Interface for hotfix detection services.
    /// Provides methods to retrieve installed hotfixes from remote computers.
    /// </summary>
    public interface IHotfixService
    {
        /// <summary>
        /// Attempts to retrieve installed hotfixes from a remote computer using WMI instead of registry scanning.
        /// This method is less likely to be flagged by antivirus scanners compared to remote registry access.
        /// </summary>
        /// <param name="hostName">Target computer hostname or IP address</param>
        /// <param name="hotfixes">Output set of discovered KB numbers (e.g., "KB4012598")</param>
        /// <returns>True if hotfixes were successfully retrieved, false otherwise</returns>
        bool TryGetInstalledHotfixes(string hostName, out HashSet<string> hotfixes);

        /// <summary>
        /// Attempts to retrieve installed hotfixes from a remote computer using WMI instead of registry scanning.
        /// This method is less likely to be flagged by antivirus scanners compared to remote registry access.
        /// </summary>
        /// <param name="hostName">Target computer hostname or IP address</param>
        /// <param name="hotfixes">Output set of discovered KB numbers (e.g., "KB4012598")</param>
        /// <param name="ui">User interface for displaying messages</param>
        /// <returns>True if hotfixes were successfully retrieved, false otherwise</returns>
        bool TryGetInstalledHotfixes(string hostName, out HashSet<string> hotfixes, IUserInterface ui);
    }
}