namespace PingCastle.ADWS;

using System.Net;

/// <summary>
/// Factory for creating AD connections (ADWS, LDAP, or Linux-based).
/// Delegates to ADWebService which handles connection type selection and fallback logic.
/// </summary>
public class ADConnectionFactory : IADConnectionFactory
{
    /// <summary>
    /// Creates a Disposable AD connection with the specified parameters.
    /// The returned connection implements IDisposable and should be disposed after use.
    /// </summary>
    /// <param name="server">The AD server hostname or IP address</param>
    /// <param name="port">The LDAP port (typically 389 for LDAP or 636 for LDAPS)</param>
    /// <param name="credential">Network credentials for the connection</param>
    /// <param name="identityProvider">Instance for connection to use to retrieve identity and convert SIDs</param>
    /// <returns>An IADWebService instance (ADWebService) ready for use</returns>
    public IADWebService CreateConnection(string server, int port, NetworkCredential credential, IIdentityProvider identityProvider, IWindowsNativeMethods nativeMethods)
    {
        // ADWebService implements IADConnection and handles:
        // - Connection type selection (ADWS, LDAP, or Linux)
        // - Fallback logic (ADWS -> LDAP if needed)
        // - Proper disposal via IDisposable
        return new ADWebService(server, port, credential, identityProvider, nativeMethods);
    }
}