namespace PingCastleCommon.Scanners;

using System;
using System.Collections.Generic;
using System.Net;
using PingCastle.Scanners;

/// <summary>
/// No-op implementation of ISmb2Protocol for non-Windows platforms.
/// Provides graceful degradation when SMB2 protocol operations are unavailable.
/// </summary>
public class NullSmb2Protocol : ISmb2Protocol
{
    public string LogPrefix { get; set; }

    public SMB2_NegotiateResponse SendNegotiateRequest(int dialect)
    {
        return default(SMB2_NegotiateResponse);
    }

    public SMB2_SessionSetupResponse SendSessionSetupRequests(NetworkCredential optionalCredential = null)
    {
        return default(SMB2_SessionSetupResponse);
    }

    public SMB2_TreeConnectResponse SendTreeConnect(string target)
    {
        return default(SMB2_TreeConnectResponse);
    }

    public List<NetworkInfo> GetNetworkInterfaceInfo()
    {
        return null;
    }
}
