namespace PingCastleCommon.RPC;

/// <summary>
/// Factory interface for creating ISSPIHelper instances.
/// Enables dependency injection of SSPI provider implementations.
/// </summary>
public interface ISSPIHelperFactory
{
    ISSPIHelper Create(string remotePrincipal);
}
