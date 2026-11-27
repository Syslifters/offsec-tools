#nullable enable

namespace PingCastleAutoUpdater.ConfigurationOrchestration;

using System.IO.Compression;

/// <summary>
/// Orchestration service for configuration extraction and merging.
/// Uses the Strategy Pattern to delegate behavior to either real or dry-run strategies.
/// </summary>
public class ConfigurationOrchestrationService : IConfigurationOrchestrationService
{
    private readonly IConfigurationStrategy _strategy;

    public ConfigurationOrchestrationService(
        ConfigurationPathContext pathContext,
        bool dryRun)
    {
        // Instantiate the appropriate strategy based on dry-run flag
        _strategy = dryRun
            ? new DryRunConfigurationStrategy(pathContext)
            : new RealConfigurationStrategy(pathContext);
    }

    public void HandleXmlConfigDuringExtraction(ZipArchiveEntry entry, string targetFilePath)
        => _strategy.HandleXmlConfigDuringExtraction(entry, targetFilePath);

    public void HandleJsonConfigDuringExtraction(ZipArchiveEntry entry, string targetFilePath)
        => _strategy.HandleJsonConfigDuringExtraction(entry, targetFilePath);

    public void PerformPostExtractionConversionAndMerge()
        => _strategy.PerformPostExtractionConversionAndMerge();

    public bool PerformInitialStateMigration()
        => _strategy.PerformInitialStateMigration();
}
