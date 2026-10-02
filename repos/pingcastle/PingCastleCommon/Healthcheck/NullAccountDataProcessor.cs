namespace PingCastleCommon.Healthcheck
{
    using PingCastle.Data;
    using System;
    using System.Collections.Generic;
    using PingCastle.ADWS;
    using PingCastle.Healthcheck;

    /// <summary>
    /// No-op account data processor for non-Windows platforms.
    /// </summary>
    public class NullAccountDataProcessor : IAccountDataProcessor
    {
        public bool ProcessAccount(IAddAccountData data, ADItem item, bool isComputer,
            DateTime dcWin2008Install, List<HealthcheckAccountDetailData> honeyPot = null)
        {
            return false;
        }

        public HealthcheckAccountDetailData GetAccountDetail(ADItem x) => new();

        public void ProcessSIDHistory(ADItem x, IAddAccountData data)
        {
            // Do nothing.
        }

        public bool IsComputerActive(ADItem x) => false;
    }
}
