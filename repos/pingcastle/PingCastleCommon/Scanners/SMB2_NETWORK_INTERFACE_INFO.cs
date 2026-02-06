namespace PingCastleCommon.Scanners;

using System;
using System.Runtime.InteropServices;

public struct SMB2_NETWORK_INTERFACE_INFO
{
    public int Next;
    public UInt32 IfIndex;
    public SMB2_NETWORK_INTERFACE_INFO_Capability Capability;
    public UInt32 Reserved;
    public UInt64 LinkSpeed;
    public UInt16 SockAddr_Storage_Family;
    [MarshalAs(UnmanagedType.ByValArray, SizeConst = 14)]
    public byte[] SockAddr_Storage_Buffer;
}