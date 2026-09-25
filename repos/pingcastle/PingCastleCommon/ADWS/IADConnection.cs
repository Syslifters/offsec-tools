using System.Security.Principal;

namespace PingCastle.ADWS
{
    public interface IADConnection
    {
        public ADDomainInfo GetDomainInfo();

        public void Enumerate(string distinguishedName, string filter, string[] properties, WorkOnReturnedObjectByADWS callback, string scope);

        public string ConvertSIDToName(string sidstring, out string referencedDomain);

        public SecurityIdentifier ConvertNameToSID(string nameToResolve);

        public IFileConnection FileConnection { get; }

        public void ThreadInitialization();
    }
}
