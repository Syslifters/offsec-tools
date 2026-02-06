namespace PingCastleAutoUpdater.ConfigurationConversion
{
    /// <summary>
    /// Interface for converting XML configuration files to JSON format.
    /// </summary>
    public interface IConfigConverter
    {
        /// <summary>
        /// Converts XML configuration file to JSON format.
        /// </summary>
        /// <param name="xmlConfigPath">Path to the XML configuration file</param>
        /// <param name="jsonConfigPath">Path where the JSON configuration file should be saved</param>
        /// <param name="deleteSourceOnCompletion"></param>
        /// <param name="createBackup">Whether to create a backup of the XML file. Only create backups for original user configs, not temp files</param>
        void ConvertXmlConfigToJson(
            string xmlConfigPath,
            string jsonConfigPath,
            bool deleteSourceOnCompletion = true,
            bool createBackup = false);
    }
}
