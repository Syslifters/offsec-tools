namespace PingCastle.misc;

using System.Collections.Generic;
using System.Text.RegularExpressions;

public static class OperatingSystemHelper
{
    public static string GetOperatingSystem(string os)
    {
        if (string.IsNullOrEmpty(os))
        {
            return "OperatingSystem not set";
        }

        os = os.Replace('\u00A0', ' ');

        var osPatterns = new Dictionary<string, string>
        {
            { @"windows(.*) 2000", "Windows 2000" },
            { @"windows server(.*) 2003", "Windows 2003" },
            { @"windows server(.*) 2008", "Windows 2008" },
            { @"windows server(.*) 2012", "Windows 2012" },
            { @"windows server(.*) 2016", "Windows 2016" },
            { @"windows server(.*) 2019", "Windows 2019" },
            { @"windows server(.*) 2022", "Windows 2022" },
            { @"windows server(.*) 2025", "Windows 2025" },
            { @"windows(.*) Embedded", "Windows Embedded" },
            { @"windows(.*) 7", "Windows 7" },
            { @"windows(.*) 8", "Windows 8" },
            { @"windows(.*) XP", "Windows XP" },
            { @"windows(.*) 10", "Windows 10" },
            { @"windows(.*) 11", "Windows 11" },
            { @"windows(.*) Vista", "Windows Vista" },
            { @"windows(.*) NT", "Windows NT" },
        };

        foreach (var pattern in osPatterns)
        {
            if (Regex.IsMatch(os, pattern.Key, RegexOptions.IgnoreCase))
            {
                return pattern.Value;
            }
        }

        return os;
    }
}