namespace PingCastleAutoUpdater.ConfigurationMerge
{
    using System.Xml;

    public interface IConfigSaver
    {
        void SaveConfig(XmlDocument config, string path);
    }
}