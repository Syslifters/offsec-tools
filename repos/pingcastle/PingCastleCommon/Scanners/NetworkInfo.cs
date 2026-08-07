namespace PingCastle.Scanners;

using System.Net;
using PingCastleCommon.Scanners;

public class NetworkInfo
{
    public SMB2_NETWORK_INTERFACE_INFO_Capability Capability { get; set; }
    public ulong LinkSpeed { get; set; }
    public IPAddress IP { get; set; }

    public uint Index { get; set; }
}