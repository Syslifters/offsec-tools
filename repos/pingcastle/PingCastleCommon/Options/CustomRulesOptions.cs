using System.Collections.Generic;

namespace PingCastleCommon.Options
{
    public class CustomRulesOptions
    {
        public const string SectionName = "CustomRules";

        public List<CustomRuleOption> CustomRules { get; set; } = new();
    }

    public class CustomRuleOption
    {
        public string RiskId { get; set; } = string.Empty;

        public int? MaturityLevel { get; set; }

        public List<ComputationOption> Computations { get; set; } = new();
    }

    public class ComputationOption
    {
        public string Type { get; set; } = string.Empty;

        public int Score { get; set; }

        public int Order { get; set; } = 1;

        public int? Threshold { get; set; }
    }
}
