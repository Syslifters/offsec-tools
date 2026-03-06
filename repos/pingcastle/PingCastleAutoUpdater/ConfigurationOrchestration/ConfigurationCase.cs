#nullable enable

namespace PingCastleAutoUpdater.ConfigurationOrchestration;

public enum ConfigurationCase
{
    NoAction,           // Case 1: No initial config, no update (also when no update in other scenarios)
    NoneToJson,         // Case 2: No initial config, JSON in update
    JsonMerge,          // Case 3: JSON config initially, JSON in update
    XmlToJsonUpdate     // Case 4: XML config initially, JSON in update
}