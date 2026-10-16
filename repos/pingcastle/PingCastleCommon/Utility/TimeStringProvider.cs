namespace PingCastleCommon.Utility;

using System;

public static class TimeStampProvider
{
    public static string LongFormatTimestamp() => $"{DateTime.Now.ToLongTimeString()}";
}