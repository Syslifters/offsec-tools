namespace PingCastleAutoUpdater.ConfigurationMerge
{
    using System;
    using System.Collections.Generic;
    using PingCastleAutoUpdater.ConfigurationConversion;

    public class ConfigMergeService
    {
        private readonly IConfigLoader _configLoader;
        private readonly IConfigMerger _configMerger;
        private readonly IConfigSaver _configSaver;
        private ConversionReport _lastReport;

        public ConfigMergeService(
            IConfigLoader configLoader,
            IConfigMerger configMerger,
            IConfigSaver configSaver)
        {
            _configLoader = configLoader;
            _configMerger = configMerger;
            _configSaver = configSaver;
        }

        /// <summary>
        /// Gets the last merge report generated
        /// </summary>
        public ConversionReport LastReport => _lastReport;

        public void MergeConfigFiles(string targetPath, string sourcePath)
        {
            if (targetPath == null)
            {
                throw new ArgumentNullException(nameof(targetPath));
            }

            if (sourcePath == null)
            {
                throw new ArgumentNullException(nameof(sourcePath));
            }

            try
            {
                var targetConfig = _configLoader.LoadConfig(targetPath);
                var sourceConfig = _configLoader.LoadConfig(sourcePath);

                var mergedConfig = _configMerger.MergeConfigs(targetConfig, sourceConfig);

                _configSaver.SaveConfig(mergedConfig, targetPath);

                // Generate report from merger results
                _lastReport = new ConversionReport
                {
                    SourcePath = sourcePath,
                    TargetPath = targetPath,
                    Success = true,
                    Timestamp = DateTime.Now,
                    SectionsConverted = new List<string>(_configMerger.MergedElements),
                    TotalSettingsMapped = _configMerger.NewElements.Count + _configMerger.MergedElements.Count,
                    IsDryRun = false
                };

                // Track new elements added
                foreach (var element in _configMerger.NewElements)
                {
                    _lastReport.MappedSettings[$"Added: {element}"] = "New element from source";
                }
            }
            catch (Exception ex)
            {
                // Create error report
                _lastReport = new ConversionReport
                {
                    SourcePath = sourcePath,
                    TargetPath = targetPath,
                    Success = false,
                    ErrorMessage = ex.Message,
                    Exception = ex,
                    Timestamp = DateTime.Now,
                    IsDryRun = false
                };
                throw;
            }
        }
    }
}