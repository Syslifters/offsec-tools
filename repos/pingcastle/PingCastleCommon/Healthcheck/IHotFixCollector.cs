// Copyright (c) Ping Castle. All rights reserved.
// https://www.pingcastle.com
//
// Licensed under the Non-Profit OSL. See LICENSE file in the project root for full license information.
//
using System.Collections.Generic;
using PingCastle.UserInterface;

namespace PingCastle.Healthcheck
{
    /// <summary>
    /// Interface for collecting installed hotfixes from remote computers.
    /// </summary>
    public interface IHotFixCollector
    {
        /// <summary>
        /// Attempts to retrieve installed hotfixes from a remote computer using WMI.
        /// If WMI fails, logs a warning and returns an empty HashSet.
        /// Hotfix retrieval requires privileged mode access to remote systems.
        /// </summary>
        /// <param name="hostName">Target computer hostname or IP address</param>
        /// <param name="hotfixes">Output set of discovered KB numbers (e.g., "KB4012598")</param>
        /// <param name="isPrivilegedMode">Whether privileged mode is active. Hotfix collection requires high privileges.</param>
        /// <returns>True if hotfixes were successfully retrieved, false if WMI failed or privilege check failed</returns>
        bool TryGetInstalledHotfixes(string hostName, out HashSet<string> hotfixes, bool isPrivilegedMode = true);

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
        bool TryGetInstalledHotfixes(string hostName, out HashSet<string> hotfixes, IUserInterface ui, bool isPrivilegedMode = true);
    }
}
