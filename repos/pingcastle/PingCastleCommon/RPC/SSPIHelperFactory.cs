namespace PingCastleCommon.RPC;

/// <summary>
/// Windows-specific factory for creating SSPIHelper instances.
/// Only instantiable on Windows platforms with SSPI support.
/// </summary>
public class SSPIHelperFactory : ISSPIHelperFactory
{
    public ISSPIHelper Create(string remotePrincipal)
    {
        return new SSPIHelper(remotePrincipal);
    }
}
