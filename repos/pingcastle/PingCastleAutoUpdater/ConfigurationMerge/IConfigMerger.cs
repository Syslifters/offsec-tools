namespace PingCastleAutoUpdater.ConfigurationMerge
{
    using System.Collections.Generic;
    using System.Xml;

    public interface IConfigMerger
    {
        IReadOnlyList<string> MergedElements { get; }

        IReadOnlyList<string> NewElements { get; }

        XmlDocument MergeConfigs(XmlDocument target, XmlDocument source);
    }
}