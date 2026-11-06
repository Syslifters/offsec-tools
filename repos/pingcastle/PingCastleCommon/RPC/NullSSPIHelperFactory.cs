namespace PingCastleCommon.RPC;

/// <summary>
/// No-op factory for creating NullSSPIHelper instances on non-Windows platforms.
/// </summary>
public class NullSSPIHelperFactory : ISSPIHelperFactory
{
    public ISSPIHelper Create(string remotePrincipal)
    {
        return new NullSSPIHelper();
    }
}
