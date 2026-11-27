#nullable enable

namespace PingCastleAutoUpdater.ConfigurationOrchestration;

using System.IO.Compression;

/// <summary>
/// Strategy interface for handling configuration extraction and merging operations.
/// Allows different implementations for real vs. dry-run execution.
/// </summary>
public interface IConfigurationStrategy
{
    /// <summary>
    /// Handle XML configuration files during extraction from the update package
    /// </summary>
    void HandleXmlConfigDuringExtraction(ZipArchiveEntry entry, string targetFilePath);

    /// <summary>
    /// Handle JSON configuration files during extraction from the update package
    /// </summary>
    void HandleJsonConfigDuringExtraction(ZipArchiveEntry entry, string targetFilePath);

    /// <summary>
    /// Perform configuration conversion and merging after extraction is complete
    /// </summary>
    void PerformPostExtractionConversionAndMerge();

    /// <summary>
    /// Perform initial state migration before update extraction.
    /// Handles case where both XML and JSON config exist initially.
    /// Migrates XML settings into JSON, with XML values taking precedence.
    /// Returns true if migration was performed, false otherwise.
    /// </summary>
    bool PerformInitialStateMigration();
}
