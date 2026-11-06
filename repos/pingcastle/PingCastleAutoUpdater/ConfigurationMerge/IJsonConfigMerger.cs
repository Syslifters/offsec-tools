namespace PingCastleAutoUpdater.ConfigurationMerge
{
    /// <summary>
    /// Interface for merging JSON configuration files while preserving user settings.
    /// </summary>
    public interface IJsonConfigMerger
    {
        /// <summary>
        /// Merges source JSON configuration into target JSON configuration.
        /// Preserves existing values in target while adding new properties from source.
        /// </summary>
        /// <param name="targetPath">Path to the target JSON configuration file</param>
        /// <param name="sourcePath">Path to the source JSON configuration file</param>
        void MergeJsonConfigFiles(string targetPath, string sourcePath);
    }
}
