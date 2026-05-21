using System.Collections.Generic;

namespace PingCastleCommon.Options
{
    public class HoneyPotOptions
    {
        public const string SectionName = "HoneyPot";

        public List<SingleHoneyPotOption> HoneyPots { get; set; } = new();
    }

    public class SingleHoneyPotOption
    {
        public string? SamAccountName { get; set; }

        public string? DistinguishedName { get; set; }
    }
}
