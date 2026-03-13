namespace PingCastle.Healthcheck;

using misc;

public static class HotFixCollectorFactory
{
    /// <summary>
    /// Creates a new instance of HotFixCollector with the default WMI hotfix service.
    /// </summary>
    /// <returns>A new HotFixCollector instance</returns>
    public static HotFixCollector Create<T>() where T : IHotfixService, new()
    {
        return new HotFixCollector(new T());
    }
}