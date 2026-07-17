using System.Net;
using System.Security.Principal;

namespace PingCastle.ADWS
{
    public interface IIdentityProvider
    {
        WindowsIdentity GetWindowsIdentityForUser(NetworkCredential credential, string remoteserver);

        string ConvertSIDToName(WindowsIdentity identity, string sidstring, string server, out string referencedDomain);
    }
}
