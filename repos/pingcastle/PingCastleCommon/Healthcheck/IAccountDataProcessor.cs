namespace PingCastleCommon.Healthcheck
{
    using PingCastle.Data;
    using System;
    using System.Collections.Generic;
    using PingCastle.ADWS;
    using PingCastle.Healthcheck;

    /// <summary>
    /// Interface for processing account data in a platform-agnostic way.
    /// </summary>
    public interface IAccountDataProcessor
    {
        /// <summary>
        /// Process account or computer data and update the account data structure.
        /// </summary>
        /// <param name="data">Account data container to update</param>
        /// <param name="item">AD item (user or computer) to process</param>
        /// <param name="isComputer">Whether this is a computer account (true) or user account (false)</param>
        /// <param name="dcWin2008Install">DC Windows 2008 installation date</param>
        /// <param name="honeyPot">Optional list of honey pot accounts to exclude</param>
        /// <returns>True if processing succeeded, false otherwise</returns>
        bool ProcessAccount(IAddAccountData data, ADItem item, bool isComputer,
            DateTime dcWin2008Install, List<HealthcheckAccountDetailData> honeyPot = null);

        /// <summary>
        /// Creates account detail data from an AD item.
        /// </summary>
        /// <param name="x">The ADItem containing the account data.</param>
        /// <returns>A new <see cref="HealthcheckAccountDetailData"/> containing the account data.</returns>
        HealthcheckAccountDetailData GetAccountDetail(ADItem x);

        /// <summary>
        /// Adds the SID history detail to the account data.
        /// </summary>
        /// <param name="x">The ADItem containing the account data.</param>
        /// <param name="data">The data object to which the history should be added.</param>
        void ProcessSIDHistory(ADItem x, IAddAccountData data);

        /// <summary>
        /// Determines if the computer detailed in the ADItem is active.
        /// </summary>
        /// <param name="x">The ADItem containing the computer data.</param>
        /// <returns><c>true</c> if the computer is active, otherwise <c>false</c></returns>
        bool IsComputerActive(ADItem x);
    }
}
