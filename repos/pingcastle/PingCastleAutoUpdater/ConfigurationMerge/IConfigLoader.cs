namespace PingCastleAutoUpdater.ConfigurationMerge
{
    using System.Xml;

    public interface IConfigLoader
    {
        XmlDocument LoadConfig(string path);
    }
}