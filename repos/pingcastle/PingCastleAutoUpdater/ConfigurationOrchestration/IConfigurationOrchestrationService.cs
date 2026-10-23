using System.IO.Compression;

namespace PingCastleAutoUpdater.ConfigurationOrchestration
{
    public interface IConfigurationOrchestrationService
    {
        void HandleXmlConfigDuringExtraction(ZipArchiveEntry entry, string targetFilePath);
        void HandleJsonConfigDuringExtraction(ZipArchiveEntry entry, string targetFilePath);
        void PerformPostExtractionConversionAndMerge();
    }
}
