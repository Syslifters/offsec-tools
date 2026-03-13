namespace PingCastleAutoUpdater.ConfigurationMerge
{
    using System;
    using System.Xml;

    public class ConfigLoader : IConfigLoader
    {
        public XmlDocument LoadConfig(string path)
        {
            try
            {
                var xmlDoc = new XmlDocument();
                xmlDoc.Load(path);
                return xmlDoc;
            }
            catch (Exception ex)
            {
                throw new ConfigException($"Failed to load config file: {path}", ex);
            }
        }
    }
}