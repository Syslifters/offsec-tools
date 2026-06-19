#nullable enable

namespace PingCastleAutoUpdater.ConfigurationOrchestration;

/// <summary>
/// Represents the configuration state before and after update extraction, used to determine
/// the appropriate configuration merge/conversion strategy.
///
/// The system supports 4 configuration cases based on:
/// 1. What configuration existed initially (JSON and/or XML)
/// 2. What configuration was in the update (JSON only)
///
/// Supported Configuration Cases:
///
/// Initial State: NONE (no config files exist)
/// ┌─────────────────────────────────────────────────────────────────────┐
/// │ Case 1: No initial config, no update        → NoAction              │
/// │ Case 2: No initial config, JSON in update   → Extract JSON          │
/// └─────────────────────────────────────────────────────────────────────┘
///
/// Initial State: JSON ONLY
/// ┌─────────────────────────────────────────────────────────────────────┐
/// │ Case 3: JSON only, JSON in update           → Merge JSON            │
/// └─────────────────────────────────────────────────────────────────────┘
///
/// Initial State: XML ONLY
/// ┌─────────────────────────────────────────────────────────────────────┐
/// │ Case 4: XML only, JSON in update            → Convert XML→JSON+Merge│
/// └─────────────────────────────────────────────────────────────────────┘
///
/// Note: Only JSON updates are supported. XML updates are no longer processed.
/// </summary>
public class ConfigurationState
{
    public bool HadJsonInitially { get; set; }
    public bool HadXmlInitially { get; set; }
    public bool HasJsonNow { get; set; }
    public bool HasXmlNow { get; set; }
    public bool HasNewJsonFromUpdate { get; set; }
    public bool HasNewXmlFromUpdate { get; set; }

    public ConfigurationCase DetermineCase()
    {
        // Short-circuit if no update
        if (!HasNewJsonFromUpdate)
        {
            return ConfigurationCase.NoAction;
        }

        // Initial State: JSON only → Merge JSON
        if (HadJsonInitially && !HadXmlInitially)
        {
            return ConfigurationCase.JsonMerge;
        }

        // Initial State: XML only → Convert XML to JSON and merge
        if (HadXmlInitially && !HadJsonInitially)
        {
            return ConfigurationCase.XmlToJsonUpdate;
        }

        // Initial State: Nothing → Extract JSON
        return ConfigurationCase.NoneToJson;
    }
}