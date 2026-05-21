using System.Net;
using System.Security.Principal;

namespace PingCastle.ADWS
{
    public class NullIdentityProvider : IIdentityProvider
    {
        public WindowsIdentity GetWindowsIdentityForUser(NetworkCredential credential, string remoteserver)
        {
            return null;
        }

        public string ConvertSIDToName(WindowsIdentity identity, string sidstring, string server, out string referencedDomain)
        {
            referencedDomain = null;
            return null;
        }
    }
}
