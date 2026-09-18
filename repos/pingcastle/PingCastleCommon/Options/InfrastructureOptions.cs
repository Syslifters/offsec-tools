using System.Collections.Generic;

namespace PingCastleCommon.Options
{
    public class InfrastructureOptions
    {
        public const string SectionName = "Infrastructure";

        public List<SingleRiverbedOption> Riverbeds { get; set; } = new();
    }

    public class SingleRiverbedOption
    {
        public string SamAccountName { get; set; } = string.Empty;
    }
}
