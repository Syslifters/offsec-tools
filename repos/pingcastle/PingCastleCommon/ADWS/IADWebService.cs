namespace PingCastle.ADWS;

using System;
using System.Collections.Generic;
using System.Net;

public interface IADWebService : IDisposable, IADConnection
{
    public delegate void Action();

    string Server { get; set; }
    int Port { get; set; }
    NetworkCredential Credential { get; set; }
    bool useLdap { get; }
    ADDomainInfo DomainInfo { get; }
    IFileConnection FileConnection { get; }
    ADDomainInfo GetDomainInfo();
    List<OUExploration> BuildOUExplorationList(string OU, int NumberOfDepthForSplit);
    void Enumerate(string distinguishedName, string filter, string[] properties, WorkOnReturnedObjectByADWS callback);
    void Enumerate(string distinguishedName, string filter, string[] properties, WorkOnReturnedObjectByADWS callback, string scope);
    void Enumerate(Action preambleWithReentry, string distinguishedName, string filter, string[] properties, WorkOnReturnedObjectByADWS callback, string scope);
    string ConvertSIDToName(string sidstring);
    string ConvertSIDToName(string sidstring, out string referencedDomain);
    System.Security.Principal.SecurityIdentifier ConvertNameToSID(string nameToResolve);
    void ThreadInitialization();
    void Dispose();
}