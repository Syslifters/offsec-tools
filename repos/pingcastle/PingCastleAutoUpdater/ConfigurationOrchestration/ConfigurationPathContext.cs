namespace PingCastleAutoUpdater.ConfigurationOrchestration;

using System;
using System.IO;

public class ConfigurationPathContext
{
    public string ExeDirectory { get; }
    public string JsonConfigPath { get; }
    public string JsonBackupPath { get; }
    public string XmlConfigPath { get; }
    public string XmlBackupPath { get; }
    public string TempJsonPath { get; }
    public string TempXmlPath { get; }

    public ConfigurationPathContext(string exeDirectory)
    {
        ExeDirectory = exeDirectory ?? throw new ArgumentNullException(nameof(exeDirectory));
        JsonConfigPath = Path.Combine(exeDirectory, "appsettings.console.json");
        JsonBackupPath = JsonConfigPath + ".bak";
        XmlConfigPath = Path.Combine(exeDirectory, "PingCastle.exe.config");
        XmlBackupPath = XmlConfigPath + ".bak";
        TempJsonPath = Path.Combine(exeDirectory, "tempNew_appsettings.console.json");
        TempXmlPath = Path.Combine(exeDirectory, "tempNew_PingCastle.exe.config");
    }

    public ConfigurationState DetectCurrentState() =>
        new()
        {
            HadJsonInitially = File.Exists(JsonBackupPath),
            HadXmlInitially = File.Exists(XmlBackupPath),
            HasJsonNow = File.Exists(JsonConfigPath),
            HasXmlNow = File.Exists(XmlConfigPath),
            HasNewJsonFromUpdate = File.Exists(TempJsonPath),
            HasNewXmlFromUpdate = File.Exists(TempXmlPath)
        };
}