// Copyright (c) Ping Castle. All rights reserved.
// https://www.pingcastle.com
//
// Licensed under the Non-Profit OSL. See LICENSE file in the project root for full license information.
//
namespace PingCastleCommon.Utility
{
    /// <summary>
    /// Interface for retrieving operating system information from remote computers.
    /// </summary>
    public interface IOperatingSystemInfoProvider
    {
        /// <summary>
        /// Retrieves the operating system name from a remote computer.
        /// </summary>
        /// <param name="computerName">The computer name to query</param>
        /// <returns>The operating system name (e.g., "Microsoft Windows Server 2025 Standard"), or empty string if retrieval fails</returns>
        string GetOperatingSystemName(string computerName);
    }
}
