namespace PingCastle;

using System;

public static class ConfigurationHelper
{
    private static IServiceProvider? _serviceProvider;

    public static void SetServiceProvider(IServiceProvider serviceProvider)
    {
        _serviceProvider = serviceProvider;
    }

    public static IServiceProvider? GetServiceProvider()
    {
        return _serviceProvider;
    }

    public static T? GetService<T>(IServiceProvider? provider = null) where T : class
    {
        var sp = provider ?? _serviceProvider;
        return sp?.GetService(typeof(T)) as T;
    }
}