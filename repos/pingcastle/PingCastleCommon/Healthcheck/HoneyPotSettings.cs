namespace PingCastle.Healthcheck;

using System;
using System.Collections.Generic;
using Microsoft.Extensions.Options;
using PingCastleCommon.Options;

public class HoneyPotSettings
{
    private static HoneyPotSettings _cachedSettings;

    public static HoneyPotSettings GetHoneyPotSettings()
    {
        if (_cachedSettings == null)
        {
            if (ServiceProviderAccessor.IsInitialized)
            {
                var options = ServiceProviderAccessor.Current.GetService(typeof(IOptions<HoneyPotOptions>)) as IOptions<HoneyPotOptions>;
                _cachedSettings = new HoneyPotSettings();
                if (options?.Value != null)
                {
                    _cachedSettings._honeyPotsCollection = new HoneyPotsCollection();
                    foreach (var potOption in options.Value.HoneyPots)
                    {
                        var potSetting = new SingleHoneyPotSettings();
                        potSetting.SamAccountName = potOption.SamAccountName;
                        potSetting.DistinguishedName = potOption.DistinguishedName;
                        _cachedSettings._honeyPotsCollection.Add(potSetting);
                    }
                }
            }
            else
            {
                throw new ApplicationException("Could not load Honeypot settings");
            }
        }

        return _cachedSettings;
    }

    private HoneyPotsCollection _honeyPotsCollection = null;

    public HoneyPotsCollection HoneyPots => _honeyPotsCollection;
}

public class HoneyPotsCollection : List<SingleHoneyPotSettings>
{
}

public class SingleHoneyPotSettings
{
    public string SamAccountName
    {
        get; set;
    }

    public string DistinguishedName
    {
        get; set;
    }
}